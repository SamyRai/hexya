# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN CGO_ENABLED=1 go build -o courses-service .

# Final stage
FROM golang:1.22
WORKDIR /
COPY --from=builder /build/courses-service /courses-service
EXPOSE 8081
CMD ["/courses-service"]
