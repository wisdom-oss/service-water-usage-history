"""Package containing the code which will be the API later on"""
import logging
import typing
import uuid

import amqp_rpc_client
import fastapi
import py_eureka_client.eureka_client
import sqlalchemy.exc

import api.handler
import configuration
import database
import database.tables
import exceptions
import models.internal
from api import security

# %% Global Clients
_amqp_client: typing.Optional[amqp_rpc_client.Client] = None
_service_registry_client: typing.Optional[py_eureka_client.eureka_client.EurekaClient] = None
_logger = logging.getLogger("API")

# %% API Setup
service = fastapi.FastAPI()
service.add_exception_handler(exceptions.APIException, api.handler.handle_api_error)
service.add_exception_handler(fastapi.exceptions.RequestValidationError, api.handler.handle_request_validation_error)
service.add_exception_handler(sqlalchemy.exc.IntegrityError, api.handler.handle_integrity_error)

# %% Configurations
_security_configuration = configuration.SecurityConfiguration()
if _security_configuration.scope_string_value is None:
    service_scope = models.internal.ServiceScope.parse_file("./configuration/scope.json")
    _security_configuration.scope_string_value = service_scope.value


# %% Routes
@service.get("/")
async def get(
    user: typing.Union[models.internal.UserAccount, bool] = fastapi.Security(
        security.is_authorized_user, scopes=[_security_configuration.scope_string_value]
    ),
    consumer=fastapi.Query(default=..., alias="consumer"),
):
    _logger.debug("Water usage history of consumer %s requested by user %s", consumer, user.username)
    query = sqlalchemy.select(
        [
            database.tables.usages.c.year,
            database.tables.usages.c.value.label("usage"),
            database.tables.usages.c.recorded.label("recorded_at"),
        ],
        database.tables.usages.c.consumer == consumer,
    ).order_by(database.tables.usages.c.year)
    _logger.debug("Built the following query: %s", query)
    _logger.debug("Executing the above printed query")
    query_result = database.engine.execute(query).all()
    _logger.debug("Finished executing the query and got all results")
    if not query_result:
        _logger.debug("There is no water usage history available for the consumer %s", consumer)
        return fastapi.Response(status_code=204)
    _logger.debug("Returning the water usage history", consumer)
    return query_result


@service.put("/file")
async def upload_file(
    usage_history: fastapi.UploadFile,
    user: typing.Union[models.internal.UserAccount, bool] = fastapi.Security(
        security.is_authorized_user, scopes=[_security_configuration.scope_string_value]
    ),
):
    pass
