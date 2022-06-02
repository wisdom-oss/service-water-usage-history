# Water Usage History
This service enables the frontend to display the usage history of a consumer

## Deployment
Use the following snippet to add this service to your custom deployment
```yaml
services:
  water-usage-history:
    build: https://github.com/wisdom-oss/service-water-usage-history#dev
    image: wisdom-oss/services/water-usage-history
    restart: always
    stop_grace_period: 1m
    deploy:
      mode: replicated
      replicas: 3
    expose:
      - 5000
    depends_on:
      - service-registry
      - message-broker
      - postgres
    environment:
      - CONFIG_DB_DSN=postgresql://postgres:<<gen-postgres-pass>>@postgres:5432/wisdom
      - CONFIG_AMQP_DSN=amqp://wisdom:<<gen-pass-rabbitmq>>@message-broker/%2F?heartbeat=600
      - CONFIG_SERVICE_REGISTRY_HOST=service-registry
      - CONFIG_SERVICE_NAME=water-usage-history
    logging:
      driver: "json-file"
      options:
        max-size: 5m
        max-file: 5
```
