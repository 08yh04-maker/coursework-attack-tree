FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY *.go ./
RUN go build -o analyzer main.go models.go analyzer.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/analyzer .

EXPOSE 8081

CMD ["./analyzer"]
