FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out -ldflags="-w -s" .

FROM alpine:latest
COPY --from=builder /app/out /app/out
WORKDIR /app
EXPOSE 8182
CMD ["./out"]
