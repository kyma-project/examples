FROM golang:1.9.0 AS builder
LABEL source=git@github.com:kyma-project/examples.git
RUN mkdir -p /concept-documentation/prometheus-custom-metrics/go/
WORKDIR /concept-documentation/prometheus-custom-metrics/go/
RUN cd /concept-documentation/prometheus-custom-metrics/go/
RUN pwd

RUN go get -d -v github.com/prometheus/client_golang/prometheus

COPY go/sample_metrics.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o sample_metrics

EXPOSE 8080
EXPOSE 8081

CMD ["./sample_metrics"]