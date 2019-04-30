FROM golang:1.12 as builder

WORKDIR /go/src/tracing
COPY src/order-front.go .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o order-front .

FROM scratch

WORKDIR /root/

COPY --from=builder /go/src/tracing .

EXPOSE 8080

ENTRYPOINT ["/root/order-front"]
