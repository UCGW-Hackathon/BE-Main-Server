FROM golang:1.24.5-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/situkang-api .

FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app && apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /out/situkang-api ./situkang-api

RUN mkdir -p /app/logs /app/uploads && chown -R app:app /app

USER app

ENV APP_ENV=production
ENV HOST_ADDRESS=0.0.0.0
ENV HOST_PORT=8080
ENV LOG_PATH=logs
ENV DB_PORT=5432
ENV DB_SSLMODE=require
ENV JWT_ACCESS_TOKEN_TTL_SECONDS=3600
ENV JWT_REFRESH_TOKEN_TTL_HOURS=720
ENV AUTO_MIGRATE=true
ENV SEED_DATA=true
ENV UPLOAD_BASE_URL=/uploads

EXPOSE 8080

CMD ["./situkang-api"]
