# Create a minimal container to run a Golang static binary
FROM golang:1.13
WORKDIR /go/src
ADD . /go/src
RUN go get -d
ARG VERSION=${VERSION}
ARG COMMIT=${COMMIT}
ARG BUILD_RFC3339=${BUILD_RFC3339}
RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildRFC3339=${BUILD_RFC3339}" -o whoami

# Copy binary to single-serve container
FROM scratch
COPY --from=0 /go/src/whoami /
ENTRYPOINT ["/whoami"]
EXPOSE 80
