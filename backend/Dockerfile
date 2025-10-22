FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
# Compilar estático (CGO disabled)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/backend

FROM gcr.io/distroless/static
COPY --from=builder /usr/local/bin/backend /backend
EXPOSE 8000
ENTRYPOINT ["/backend"]
