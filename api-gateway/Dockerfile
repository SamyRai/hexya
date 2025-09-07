# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /build
COPY go.mod .
COPY main.go .
RUN go build -o api-gateway .

# Final stage
FROM scratch
COPY --from=builder /build/api-gateway /api-gateway
EXPOSE 8080
CMD ["/api-gateway"]
