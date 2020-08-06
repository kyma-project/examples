FROM golang:1.14-alpine as builder

ENV BASE_APP_DIR=/go/src/github.com/kyma-project/examples/order-service \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR ${BASE_APP_DIR}

COPY . .

RUN go build -ldflags "-s -w" -a -o service main.go \
    && mkdir /app \
    && mv ./service /app/service

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
LABEL source = git@github.com:kyma-project/examples.git

COPY --from=builder /app /app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/app/service"]
