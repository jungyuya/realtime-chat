# =========================================
# 1단계: 빌더(Builder) 스테이지
# =========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 루트 컨텍스트 기준 경로
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download

# 백엔드 소스 코드만 복사
COPY apps/backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/main.go


# =========================================
# 2단계: 최종(Final) 스테이지
# =========================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

# ca-certificates는 외부 HTTPS 통신 시 필요할 수 있으므로,
# 안정성을 위해 일단 포함하는 것이 좋습니다.
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["/app/server"]