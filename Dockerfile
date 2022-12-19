FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL maintainer="anushasankaranarayanan@github.com"

WORKDIR /go

COPY build/microservice api/openapi.json /go/

RUN chown 1001:root /go

USER 1001

CMD ["./microservice"]
