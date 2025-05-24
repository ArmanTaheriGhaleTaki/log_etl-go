FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o log_etl .
FROM alpine:latest
WORKDIR /root/

# Copy the built binary from the builder stage
# COPY --from=builder /app/log_etl .

# Expose the port the application runs on
# EXPOSE 8080

CMD ["./log_etl"]