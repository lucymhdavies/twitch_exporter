# From:
# https://github.com/prometheus/client_golang/blob/master/examples/simple/Dockerfile
#
# This Dockerfile builds an image for a client_golang example.
#
# Use as (from the root for the client_golang repository):
#    docker build -f examples/$name/Dockerfile -t prometheus/golang-example-$name .

# Builder image, where we build the example.
FROM golang:1.10.0 AS builder
WORKDIR /go/src/github.com/lucymhdavies/go-experiments/prometheus/twitch
ADD . .
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o app

# Final image.
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/lucymhdavies/go-experiments/prometheus/twitch .
EXPOSE 8080
ENTRYPOINT ["/app"]
