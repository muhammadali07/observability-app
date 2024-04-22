# Stage 1: Build aplikasi Go
FROM golang:1.22.1 AS build-stage

WORKDIR /app

# Salin file go.mod dan go.sum terlebih dahulu untuk mengoptimalkan caching
COPY go.mod go.sum ./
RUN go mod download

# Salin kode aplikasi
COPY *.go ./

# Build aplikasi Go tanpa dependensi eksternal
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/appname

# Stage 2: Image akhir
FROM alpine:latest

# Tidak diperlukan karena tidak ada dependensi yang dibutuhkan
# Namun, jika ada dependensi, Anda harus menginstalnya di sini

# Salin binary aplikasi dari stage 1
COPY --from=build-stage /app/appname /app/appname

# Tentukan port yang digunakan oleh aplikasi
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat kontainer dimulai
CMD ["/app/appname"]
