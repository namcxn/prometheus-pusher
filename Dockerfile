FROM golang:1.13.7-alpine3.11 as builder

WORKDIR /app
ADD . /app
RUN go get
RUN go build main.go

FROM alpine:latest
WORKDIR /app
COPY queries.yaml .
COPY config-example.yaml ./config.yaml
COPY --from=builder /app/main .
CMD ["./main"]