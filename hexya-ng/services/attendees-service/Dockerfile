# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod .
COPY main.go .
RUN go build -o attendees-service .

# Final stage
FROM scratch
COPY --from=builder /build/attendees-service /attendees-service
EXPOSE 8083
CMD ["/attendees-service"]
