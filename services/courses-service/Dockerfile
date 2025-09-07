# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /build
COPY go.mod .
COPY main.go .
RUN go build -o courses-service .

# Final stage
FROM scratch
COPY --from=builder /build/courses-service /courses-service
EXPOSE 8081
CMD ["/courses-service"]
