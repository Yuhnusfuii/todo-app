# === GIAI ĐOẠN 1: BUILD BINARY ===
FROM golang:1.23-alpine AS builder

# Cài đặt git và các công cụ cần thiết (nếu có)
RUN apk add --no-cache git

WORKDIR /app

# Copy các file quản lý thư viện trước để tận dụng Docker Cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn vào container
COPY . .

# Biên dịch ứng dụng thành file nhị phân tên là "main" tĩnh (CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# === GIAI ĐOẠN 2: RUNTIME TINH GỌN ===
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Chỉ copy file binary đã build từ giai đoạn trước qua đây
COPY --from=builder /app/main .

# Mở port 3000 giống như cấu hình trong main.go
EXPOSE 3000

# Lệnh khởi chạy ứng dụng
CMD ["./main"]