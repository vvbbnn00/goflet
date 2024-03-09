FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goflet .

FROM alpine:latest

WORKDIR /app/
COPY --from=builder /app/goflet .
VOLUME /app/data

EXPOSE 8080

CMD ["./goflet"]
