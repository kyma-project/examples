FROM golang:1.12 as builder
ARG DOCK_PKG_DIR=/go/src/github.com/kyma-project/examples/http-db-service/
ADD . $DOCK_PKG_DIR
WORKDIR $DOCK_PKG_DIR
RUN go get -t -d -v -insecure ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN go test ./...


FROM scratch
LABEL source=git@github.com:kyma-project/examples.git 
WORKDIR /app/
COPY --from=builder /go/src/github.com/kyma-project/examples/http-db-service/main /app/
COPY --from=builder /go/src/github.com/kyma-project/examples/http-db-service/docs/api/api.yaml /app/
CMD ["./main"]

EXPOSE 8017:8017