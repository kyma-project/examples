FROM golang:1.12-alpine

ENV SRC_DIR=/go/src/github.com/kyma-project/examples/tests/http-db-service

ADD . $SRC_DIR

WORKDIR $SRC_DIR

LABEL source=git@github.com:kyma-project/examples.git

ENTRYPOINT go test -v ./...