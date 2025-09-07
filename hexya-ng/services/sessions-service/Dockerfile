# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod .
COPY main.go .
RUN go build -o sessions-service .

# Final stage
FROM scratch
COPY --from=builder /build/sessions-service /sessions-service
EXPOSE 8082
CMD ["/sessions-service"]
