FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
RUN go build -o main .

# Expose port
EXPOSE 3001

# Command untuk run
CMD ["./main"]