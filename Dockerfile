FROM node:22-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:alpine AS backend-builder
ENV GOPROXY=https://goproxy.cn,direct
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /panel ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata su-exec
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=backend-builder /panel .
COPY --from=frontend-builder /app/web/dist ./dist
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

VOLUME ["/app/data"]
EXPOSE 5678

ENV PORT=5678
ENV GIN_MODE=release
ENV DIST_DIR=./dist
ENV DATA_DIR=./data
ENV UPLOAD_DIR=./data/uploads
ENV PUID=1000
ENV PGID=1000

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./panel"]