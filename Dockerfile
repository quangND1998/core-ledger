# Builder stage - dùng Go >= 1.25.3 theo go.mod
FROM golang:1.25.3-alpine AS builder
RUN apk add --no-cache git build-base

WORKDIR /src

# copy tất cả (hỗ trợ replace local)
COPY . .

# optional: use module proxy
RUN go env -w GOPROXY=https://proxy.golang.org,direct

# download deps và build 2 binary: core app và dashboard
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/core-ledger ./ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/dashboard ./dashboard

# Final runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/core-ledger .
COPY --from=builder /app/dashboard .
ENV GIN_MODE=release
# optional TZ
ENV TZ=Asia/Ho_Chi_Minh
# default ports (mappings set in docker-compose)
EXPOSE 8080 8081
CMD ["./core-ledger"]