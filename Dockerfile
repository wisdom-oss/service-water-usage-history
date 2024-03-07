FROM golang:alpine AS build-service
COPY . /tmp/src
WORKDIR /tmp/src
RUN mkdir -p /tmp/build && \
    go mod download && \
    go build -o /tmp/build/app

FROM alpine:latest
COPY --from=build-service /tmp/build/app /service
COPY resources /
ENTRYPOINT ["/service"]
EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=15s CMD /service -healthcheck