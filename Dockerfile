# STAGE 1: Build binary
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Build binary yang statis (tanpa ketergantungan library luar)
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# STAGE 2: Runner (Gunakan Alpine agar ukuran < 50MB)
FROM alpine:latest
WORKDIR /root/
# Salin binary dari stage builder
COPY --from=builder /app/main .
# Salin .env (opsional, di produksi biasanya pakai Env Var Cloud)
COPY --from=builder /app/.env . 

EXPOSE 8080
CMD ["./main"]