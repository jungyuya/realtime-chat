# =========================================
# 1단계: 빌더(Builder) 스테이지
# =========================================
FROM node:22-alpine AS builder

RUN npm install -g pnpm

WORKDIR /app

# 의존성 설치에 필요한 모든 설정 파일 복사
COPY pnpm-lock.yaml ./
COPY package.json ./
COPY pnpm-workspace.yaml ./
COPY apps/frontend/package.json ./apps/frontend/
COPY apps/backend/package.json ./apps/backend/

# devDependencies 포함하여 모든 의존성 설치
RUN pnpm install --filter frontend --frozen-lockfile

# 빌드에 필요한 나머지 모든 소스 코드 복사
COPY . .

# 빌드 인자 설정
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# 빌드 실행
RUN pnpm --filter frontend build

# =========================================
# 2단계: 최종(Final) 스테이지
# =========================================
FROM nginx:1.27-alpine

COPY --from=builder /app/apps/frontend/dist /usr/share/nginx/html

COPY apps/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"] 