FROM golang:1.13.7-alpine3.11 as builder

WORKDIR /app
ADD . /app
RUN go get
RUN go build -o prometheus-pusher

FROM alpine:3.12

WORKDIR /app
COPY ./config-example.yaml /app/config.yaml
COPY --from=builder /app/prometheus-pusher /app/prometheus-pusher
CMD ["/app/prometheus-pusher"]