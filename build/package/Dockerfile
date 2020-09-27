# Create a minimal container to run a Golang static binary
FROM golang:1.13-alpine as build

WORKDIR /go/src/app

# Create appuser.
# See https://stackoverflow.com/a/55757473/12429735RUN
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

WORKDIR /go/src
ADD . /go/src
RUN go get -d
ARG VERSION=${VERSION}
ARG COMMIT=${COMMIT}
ARG BUILD_RFC3339=${BUILD_RFC3339}
RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildRFC3339=${BUILD_RFC3339}" -o whoami

# Copy binary to single-serve container
FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /go/src/whoami /
USER appuser:appuser
ENTRYPOINT ["/whoami"]
EXPOSE 8000
