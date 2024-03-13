FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goflet -trimpath \
    -ldflags="-s -w -X 'github.com/vvbbnn00/goflet/base.Version=$(git tag --sort=-creatordate | head -n 1) (commit:$(git rev-parse --short HEAD))'"

FROM alpine:3.19

WORKDIR /data/

COPY --from=builder /app/goflet /bin/goflet
VOLUME /data/
EXPOSE 8080

CMD ["goflet"]
