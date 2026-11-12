# --- Stage 1: Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git jika diperlukan oleh go modules
RUN apk add --no-cache git

# Salin dependensi terlebih dahulu untuk caching
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh source code project
COPY . .

# Kompilasi aplikasi menjadi binary statis
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/schwi ./cmd/schwi

# --- Stage 2: Final Production Stage ---
FROM alpine:latest

WORKDIR /app

# Install sertifikat SSL agar bot bisa terhubung ke Discord API via HTTPS
RUN apk --no-cache add ca-certificates

# Salin hasil binary dari stage builder
COPY --from=builder /app/bin/schwi /app/schwi

# Ekspose port HTTP sesuai dengan konfigurasi PORT=3005 di .env Anda
EXPOSE 3005

# Jalankan bot
CMD ["/app/schwi"]