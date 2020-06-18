FROM golang:1.13.7-alpine3.11 as builder

WORKDIR /app
ADD . /app
RUN go get
RUN go build -o prometheus-pusher

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/prometheus-pusher .
CMD ["./prometheus-pusher"]