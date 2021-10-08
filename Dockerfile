FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /wedding-api


#FROM gcr.io/distroless/base-debian10
FROM golang:latest

WORKDIR /

COPY --from=builder /wedding-api /wedding-api
COPY --from=builder /app/config.* /

ENV PORT=80

ENTRYPOINT ["/wedding-api"]
