FROM golang:1.27.4-alpine AS builder

WORKDIR /src

# Сначала только манифесты — чтобы кэш слоёв работал при изменении кода
COPY go.mod go.sum ./
RUN go mod download

# Затем весь остальной код
COPY . .

# CGO выключен: получаем статический бинарник, не зависящий от libc
RUN CGO_ENABLED=0 GOOS=linux go build \
        -trimpath \
        -ldflags="-s -w" \
        -o /out/threadapi \
        ./cmd/threadapi

# ─── runtime stage ───────────────────────────────────────────
FROM alpine:3.20

# ca-certificates — для SMTP TLS и HTTPS-запросов (hcaptcha, oauth);
# tzdata — если где-то используется time.LoadLocation
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 app

WORKDIR /app

COPY --from=builder /out/threadapi /app/threadapi

USER app

EXPOSE 3000

ENTRYPOINT ["/app/threadapi"]