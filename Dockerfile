FROM golang:latest as builder

WORKDIR /app

COPY ./ ./

RUN go build -o main .

FROM alpine:latest

COPY --from=builder /app /app

CMD ["/app/main"]
