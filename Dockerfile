# ---- Build Stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o globalizer ./cmd/globalizer

# ---- Runtime Stage ----
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/globalizer .

ENV GLOBALIZER_DB_PATH=/data/globalizer.db
VOLUME ["/data"]

EXPOSE 8080

ENTRYPOINT ["./globalizer"]
CMD ["serve"]
