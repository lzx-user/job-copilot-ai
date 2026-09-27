FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/job-copilot ./cmd/server

FROM alpine:3.22
RUN adduser -D -u 10001 appuser
COPY --from=builder /out/job-copilot /app/job-copilot
USER appuser
ENV APP_ENV=production APP_PORT=8080
EXPOSE 8080
CMD ["/app/job-copilot"]
