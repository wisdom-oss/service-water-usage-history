"""Package containing the code which will be the API later on"""
import datetime
import email.utils
import hashlib
import logging
import typing
import uuid

import amqp_rpc_client
import fastapi
import py_eureka_client.eureka_client
import pytz as pytz
import sqlalchemy.exc
import orjson
import database.tables
import api.handler
import configuration
import database
import exceptions
import models.internal
import tools
from api import security


# %% Global Clients
_amqp_client: typing.Optional[amqp_rpc_client.Client] = None
_service_registry_client: typing.Optional[py_eureka_client.eureka_client.EurekaClient] = None

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
    consumer: uuid.UUID = fastapi.Query(default=..., alias="consumer"),
):
    query = sqlalchemy.select(
        [
            database.tables.usages.c.year,
            database.tables.usages.c.value.label("usage"),
            database.tables.usages.c.recorded.label("recorded_at"),
        ],
        database.tables.usages.c.consumer == consumer,
    ).order_by(database.tables.usages.c.year)
    query_result = database.engine.execute(query).fetchall()
    if len(query_result) == 0:
        return fastapi.Response(status_code=204)
    return query_result
