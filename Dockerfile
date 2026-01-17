# Build stage
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY main.go .

COPY ./static ./static

COPY ./templates ./templates

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage - use Alpine for minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]