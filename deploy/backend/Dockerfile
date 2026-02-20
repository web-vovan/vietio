# =========================
# Stage 1 — Builder
# =========================
FROM golang:1.25-bookworm AS builder

# Устанавливаем зависимости для CGO (goheif)
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libheif-dev \
    libde265-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Сначала копируем только go.mod/sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Затем копируем остальной код
COPY . .

# Включаем CGO
ENV CGO_ENABLED=1

# Собираем бинарник
RUN go build -ldflags="-s -w" -o app ./cmd


# =========================
# Stage 2 — Runtime
# =========================
FROM debian:bookworm-slim

# Только runtime-библиотеки
RUN apt-get update && apt-get install -y \
    libheif1 \
    libde265-0 \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/app .

EXPOSE 8888

CMD ["./app"]
