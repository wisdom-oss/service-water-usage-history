# WISdoM OSS - Water Usage History Service
![Golang Version](https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-water-usage-history?filename=src%2Fgo.mod&style=for-the-badge)
[![OpenAPI Documentation](https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative)](./openapi.yaml)

This microservice allows authorized users to get the water usage history
of a single consumer.

## Usage
To use this service in the default deployment you need to be a member of the
`usage-history` group. If this group does not exist in your identity provider
you may need to create it before being able to access this microservice.

Further documentation is available in the 
[`docs` repository](https://github.com/wisdom-oss/docs)

## Deployment
This service is deployed by default on the WISdoM Platform. Therefore no action
is needed by you to deploy the microservice


