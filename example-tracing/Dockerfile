FROM golang:1.10 as builder

WORKDIR /go/src/example-tracing
COPY order-front.go .
RUN ls ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o order-front .

FROM scratch

WORKDIR /root/

COPY --from=builder /go/src/example-tracing .

EXPOSE 8080

ENTRYPOINT ["/root/order-front"]
