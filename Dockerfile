FROM golang:1.19.4 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /wedding-api


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /wedding-api /wedding-api
COPY --from=builder /app/config.* /

ENV PORT=80
ENV GIN_MODE=release

ENTRYPOINT ["/wedding-api"]
