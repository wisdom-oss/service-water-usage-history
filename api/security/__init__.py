import http
import logging
import typing

import amqp_rpc_client
import fastapi.security

import configuration
import enums
import exceptions
import models.amqp
import models.internal

# %% OAuth 2.0 Scheme Setup
__wisdom_central_auth = fastapi.security.OAuth2PasswordBearer(
    tokenUrl="/api/auth/token",
    scheme_name="WISdoM Central Authorization",
    auto_error=False,
)

_logger = logging.getLogger("api.security")

# %% Required Settings for the common packages
_service_settings = configuration.ServiceConfiguration()
_amqp_settings = configuration.AMQPConfiguration()


# %% Clients needed for the security
_amqp_client = amqp_rpc_client.Client(_amqp_settings.dsn)
__logger = logging.getLogger("security")


def is_authorized_user(
    scopes: fastapi.security.SecurityScopes,
    access_token: str = fastapi.Depends(__wisdom_central_auth),
) -> typing.Union[bool, models.internal.UserAccount]:
    """
    Check if the user calling this service is authorized.

    This security dependency needs to be used as fast api dependency in the methods

    :param scopes: The scopes this used needs to have to access this service
    :type scopes: list
    :param access_token: The access token used by the user to access the service
    :type access_token: str
    :return: Status of the authorization
    :rtype: bool
    :raises exceptions.APIException: The user is not authorized to access this service
    """
    if access_token is None:
        raise exceptions.APIException(
            error_code="MISSING_BEARER_TOKEN",
            error_title="Authorization Information Missing",
            error_description="The request did not contain any authorization information",
            http_status=http.HTTPStatus.BAD_REQUEST,
        )
    # Prepare the request
    _logger.debug("Creating a new access token introspection request")
    introspection_request = models.amqp.TokenIntrospectionRequest(bearer_token=access_token, scope=scopes.scope_str)
    _logger.debug("Created the following token introspection request: %s", introspection_request.json())
    # Send the request and wait a max amount of 10 seconds until the response needs to be returned
    _logger.debug("Sending the token introspection request to the AMQP RPC authorization service")
    introspection_id = _amqp_client.send(
        introspection_request.json(by_alias=True),
        _amqp_settings.authorization_exchange,
        "authorization-service",
    )
    _logger.debug("Request successfully sent. Waiting for response")
    introspection_response_bytes = _amqp_client.await_response(introspection_id, 10)
    _logger.debug("Waiting for the response ended")
    if introspection_response_bytes is None:
        _logger.error("Token introspection timeout occurred")
        raise exceptions.APIException(
            error_code="TOKEN_VALIDATION_TIMEOUT",
            error_title="Token Validation Timeout",
            error_description="The service could not validate the used access token in a timely manner",
            http_status=http.HTTPStatus.REQUEST_TIMEOUT,
        )
    # Try to read the response
    _logger.debug("Trying to parse the response of the amqp rpc authorization service")
    token = models.internal.TokenIntrospection.parse_raw(introspection_response_bytes)
    _logger.debug("Parsed the following information from the response: %s", token.json())
    if not token.active:
        _logger.warning("An inactive token was used to access this service")
        match token.reason:
            case enums.TokenIntrospectionFailure.INVALID_TOKEN:
                raise exceptions.APIException(
                    error_code="INVALID_TOKEN",
                    error_title="Invalid Bearer Token",
                    error_description="The request did not contain the correct credentials to allow processing this "
                    "request",
                    http_status=http.HTTPStatus.UNAUTHORIZED,
                )
            case enums.TokenIntrospectionFailure.TOKEN_USED_TOO_EARLY:
                raise exceptions.APIException(
                    error_code="EXPIRED_TOKEN",
                    error_title="Expired Bearer Token",
                    error_description="The request did not contain a alive Bearer token",
                    http_status=http.HTTPStatus.UNAUTHORIZED,
                )
            case enums.TokenIntrospectionFailure.EXPIRED:
                raise exceptions.APIException(
                    error_code="EXPIRED_TOKEN",
                    error_title="Expired Bearer Token",
                    error_description="The request did not contain a alive Bearer token",
                    http_status=http.HTTPStatus.UNAUTHORIZED,
                )

            case enums.TokenIntrospectionFailure.TOKEN_USED_TOO_EARLY:
                raise exceptions.APIException(
                    error_code="TOKEN_BEFORE_CREATION",
                    error_title="Credentials used too early",
                    error_description="The credentials used for this request are currently not valid",
                    http_status=http.HTTPStatus.UNAUTHORIZED,
                )
            case enums.TokenIntrospectionFailure.NO_USER_ASSOCIATED:
                raise exceptions.APIException(
                    error_code="USER_DELETED",
                    error_title="User deleted",
                    error_description="The account used to access this resource was deleted",
                    http_status=http.HTTPStatus.UNAUTHORIZED,
                )
            case enums.TokenIntrospectionFailure.USER_DISABLED:
                raise exceptions.APIException(
                    error_code="USER_DISABLED",
                    error_title="User Disabled",
                    error_description="The account used to access this resource is currently disabled",
                    http_status=http.HTTPStatus.FORBIDDEN,
                )
            case enums.TokenIntrospectionFailure.MISSING_PRIVILEGES:
                raise exceptions.APIException(
                    error_code="MISSING_PRIVILEGES",
                    error_title="Missing Privileges",
                    error_description="The account used to access this resource does not have the privileges to access "
                    "this endpoint",
                    http_status=http.HTTPStatus.FORBIDDEN,
                )
            case _:
                raise exceptions.APIException(
                    error_code="INACTIVE_TOKEN",
                    error_title="Inactive Bearer Token",
                    error_description="The token was rejected by the authorization system, but no error code was "
                    "returned",
                    http_status=http.HTTPStatus.UNAUTHORIZED,
                )
    if token.user is None:
        _logger.debug("The user accessing this service is authorized, but no information about the used is present")
        return True
    _logger.debug("The user '%s' accessing this service is authorized to access this service", token.user.username)
    return token.user
