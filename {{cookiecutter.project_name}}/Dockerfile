# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /build
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download
COPY src/ ./

ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 go build -o app ./cmd/server

FROM alpine:latest AS runtime

RUN mkdir -m 0755 -p /srv
COPY --from=build /build/app /srv/app

ARG APP_UID="2000"
ARG APP_GID="2000"
ARG APP_USER="app"
ARG APP_GROUP="app"

RUN \
    addgroup -g "${APP_GID}" "${APP_GROUP}" && \
    adduser \
      -u "${APP_UID}" \
      -g "" \
      -G "${APP_GROUP}" \
      -h "/home/${APP_USER}" \
      -D \
      -s "/bin/sh" \
      "${APP_USER}"

ENTRYPOINT [ "/srv/app" ]
