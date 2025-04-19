# Create a minimal container to run a Golang static binary
FROM golang:1.24-alpine AS build
WORKDIR /go/src/app

# Create appuser.
# See https://stackoverflow.com/a/55757473/12429735
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN apk update && apk add git ca-certificates

COPY go.mod .
COPY go.sum .
RUN go mod download

ARG APPLICATION="go-executable"
ARG BUILD_RFC3339="1970-01-01T00:00:00Z"
ARG PACKAGE="user/repo"
ARG REVISION="local"
ARG VERSION="dirty"
ARG GO_LDFLAGS="-w -s \
    -X github.com/jnovack/release.Application=${APPLICATION} \
    -X github.com/jnovack/release.BuildRFC3339=${BUILD_RFC3339} \
    -X github.com/jnovack/release.Package=${PACKAGE} \
    -X github.com/jnovack/release.Revision=${REVISION} \
    -X github.com/jnovack/release.Version=${VERSION} \
    "

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "${GO_LDFLAGS}" -o /go/bin/${APPLICATION} cmd/${APPLICATION}/*

# Copy binary to single-serve container
FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER appuser:appuser

ARG APPLICATION="go-executable"
ARG BUILD_RFC3339="1970-01-01T00:00:00Z"
ARG DESCRIPTION="no description"
ARG PACKAGE="user/repo"
ARG REVISION="local"
ARG VERSION="dirty"

LABEL org.opencontainers.image.ref.name="${PACKAGE}" \
    org.opencontainers.image.created=$BUILD_RFC3339 \
    org.opencontainers.image.authors="Justin J. Novack <jnovack@gmail.com>" \
    org.opencontainers.image.documentation="https://github.com/${PACKAGE}/README.md" \
    org.opencontainers.image.description="${DESCRIPTION}" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.source="https://github.com/${PACKAGE}" \
    org.opencontainers.image.revision=$REVISION \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.url="https://hub.docker.com/r/${PACKAGE}/"

COPY --from=build /go/bin/${APPLICATION} /app

EXPOSE 8000

ENTRYPOINT ["/app"]
