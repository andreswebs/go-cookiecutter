# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:{{ cookiecutter.go_version }}-alpine AS build

WORKDIR /build
# go.su[m] is an optional glob: the API shape has no third-party dependencies and
# therefore no go.sum, while the CLI shape does.
COPY go.mod go.su[m] ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ARG VERSION=dev
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 go build -trimpath \
    -ldflags="-s -w -buildid= -X github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/version.Override=${VERSION}" \
    -o {{ cookiecutter.project_name }} ./cmd/{% if cookiecutter.project_type == 'cli' %}{{ cookiecutter.project_name }}{% else %}server{% endif %}

FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates

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

COPY --from=build /build/{{ cookiecutter.project_name }} /usr/local/bin/{{ cookiecutter.project_name }}

USER ${APP_USER}
{% if cookiecutter.project_type == 'api' %}
EXPOSE 8080
{% endif %}
ENTRYPOINT [ "/usr/local/bin/{{ cookiecutter.project_name }}" ]
