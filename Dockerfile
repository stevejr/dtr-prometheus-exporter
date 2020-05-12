FROM golang:1.13-alpine as build

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o app

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
LABEL maintainer="srichards@mirantis.com"

# Ballerina runtime distribution filename.
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

# Labels.
LABEL com.mirantis.schema-version="1.0"
LABEL com.mirantis.build-date=$BUILD_DATE
LABEL com.mirantis.name="dockerps/dtr-prometheus-exporter"
LABEL com.mirantis.description="Docker Trusted Registry custom Prometheus metrics exporter"
LABEL com.mirantis.url="https://mirantis.com/"
LABEL com.mirantis.vcs-url="https://github.com/stevejr/dtr-prometheus-exporter"
LABEL com.mirantis.vcs-ref=$VCS_REF
LABEL com.mirantis.version=$BUILD_VERSION
LABEL com.mirantis.docker.cmd="docker run \
-d \
-p 9580:9580 \
--mount type=bind,source=[YOUR DTR CERTS DIR],target=/dtrcerts,readonly \
-e CONNECTION_STRING=[YOUR CONNECTION STRING] \
-e DTR_CA=/dtrcerts/[YOUR CA.PEM FILENAME] \
-e DTR_CERT=/dtrcerts/[YOUR CERT.PEM FILENAME] \
-e DTR_KEY=/dtrcerts/[YOUR KEY.PEM FILENAME] \
-e DTR_USERNAME=[YOUR DTR USERNAME] \
-e DTR_PASSWORD=[YOUR DTR PASSWORD] \
dockerps/dtr-prometheus-exporter:alpine"

COPY --from=build /build /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 9580
ENTRYPOINT [ "/app" ]