# Realtime Chat Project
> **"GKE 기반의 MSA 아키텍처에서 비용 최적화를 위해 단일 VM 환경으로 마이그레이션하여 운영 비용 $0를 달성한 실시간 채팅 서비스"**

## 1. 프로젝트 소개

본 프로젝트는 초기에 **Google Kubernetes Engine (GKE)** 기반의 고가용성 마이크로서비스 아키텍처로 시작하여, 실제 서비스 규모와 비용 효율성을 고려해 **단일 VM (Docker Compose)** 환경으로 성공적으로 마이그레이션한 엔지니어링 포트폴리오입니다.

Go의 강력한 동시성 처리 모델을 활용하여 고성능 실시간 통신을 구현했으며, 복잡한 클라우드 인프라(K8s, ArgoCD)를 직접 구축하고 운영해본 경험을 바탕으로, 상황에 맞는 **"적정 기술(Appropriate Technology)"** 을 선택하고 전환하는 과정을 담았습니다.

- **Period**: 2025.11 - 2025.12 (개인 프로젝트)
- **Tech Stack**: Go, React, Docker, Terraform, Nginx

---

## 2. 시스템 아키텍처 변화

비용 최적화와 유지보수 효율성을 위해 인프라 구조를 대대적으로 리팩토링했습니다.

### 2.1. Phase 1: GKE 기반 고가용성 아키텍처 (Initial)
> **"확장성과 무중단 배포에 초점을 맞춘 Cloud-Native 학습 단계"**

![!GKE Architecture](https://github.com/user-attachments/assets/placeholder-gke)

- **Kubernetes (GKE)**: 컨테이너 오케스트레이션을 통한 자동 복구 및 스케일링
- **Argo CD & GitOps**: Git을 단일 진실 공급원(SSOT)으로 하는 선언적 배포 파이프라인
- **Observability**: Prometheus & Grafana를 활용한 메트릭 수집 및 시각화

### 2.2. Phase 2: 단일 VM 비용 최적화 아키텍처 (Current)
> **"월 유지비용 $0를 목표로 한 실용주의적 아키텍처"**

![!Single VM Architecture](https://github.com/user-attachments/assets/placeholder-vm)

- **Docker Compose**: 복잡한 K8s 리소스를 간소화하여 단일 호스트 내에서 서비스 오케스트레이션
- **Nginx Reverse Proxy**: 로드 밸런싱 및 SSL/TLS 종단 처리
- **SSH Tunneling**: 외부 IP 노출 없이 로컬에서 안전하게 원격 DB에 접속 및 관리
- **Cost Optimization**: GKE 클러스터 비용 제거 및 프리티어 VM 활용으로 운영 비용 최소화

---

## 3. 핵심 기능 및 주요 성과

### 3.1. Infrastructure Cost Optimization
> **"GKE에서 Single VM으로, 성능 저하 없이 월 $0 달성"**

초기에는 학습 목적으로 GKE를 도입했으나, 개인 프로젝트 규모에 비해 높은 비용(Cluster Management Fee 등)이 발생했습니다. 이를 해결하기 위해 아키텍처 다이어트(Architecture Diet)를 감행했습니다.
- **마이그레이션**: Kubernetes Manifest를 Docker Compose로 변환하고, 인그레스 컨트롤러를 Nginx 설정으로 최적화했습니다.
- **성과**: 월 10만원 이상 발생하던 고정 클라우드 비용을 **$0(GCP/Oracle Free Tier)** 로 줄이면서도, 서비스의 핵심인 실시간성은 그대로 유지했습니다.

### 3.2. WebSocket 기반 실시간 통신
> **"Go Goroutine을 활용한 지연 없는 양방향 통신"**

- **Go Backend**: Node.js 대비 적은 메모리로 대규모 동시 접속을 처리하기 위해 Go 언어를 선택했습니다. Goroutine과 Channel을 활용해 효율적인 메시지 브로커를 직접 구현했습니다.
- **Connection Health**: Heartbeat 메커니즘을 도입하여 비정상적인 연결 종료를 즉시 감지하고 리소스를 해제합니다.
- **Stateless Auth**: JWT 기반 인증을 통해 소켓 연결 시 오버헤드를 줄이고 서버 확장성을 확보했습니다.

### 3.3. IaC & GitOps Automation
> **"인프라의 코드화(IaC)로 리전마저 자유롭게 이동"**

- **Terraform**: 단순한 VM 생성뿐만 아니라 VPC, 방화벽, 서브넷 등 모든 네트워크 환경을 코드로 정의했습니다. 덕분에 리전 변경(Asia → US)이나 계정 마이그레이션 시 `terraform apply` 한 번으로 동일한 환경을 10분 내에 복구할 수 있었습니다.
- **CI/CD Pipeline Evolution**:
    - (Early) GitHub Actions + Argo CD: 이미지 빌드 및 GitOps 배포
    - (Current) GitHub Actions + SSH Remote: Docker Multi-stage Build로 이미지를 **90% 이상 경량화**하고, 단일 VM에 최적화된 SSH 기반 배포 자동화 구축

### 3.4. Security & UX Engineering

- **SSL/TLS Automation**: Certbot을 Docker 컨테이너로 띄워 Nginx와 연동, Let's Encrypt 인증서의 발급과 갱신을 완전 자동화했습니다.
- **Seamless Integration**: `iframe` 위젯 형태로 블로그에 채팅 서비스를 통합하면서도, **반응형 디자인**과 **자동 재연결(Auto-reconnect)** 로직을 통해 사용자가 네트워크 단절을 인지하지 못하도록 UX를 고도화했습니다.

---

## 4. 기술 스택

| 분류 | 기술 | 선정 이유 및 활용 |
|:---:|:---:|---|
| **Backend** | **Go (Golang)** | Goroutine을 이용한 고성능 동시성 처리, 정적 타입 언어의 안정성 |
| | **Gin Framework** | 가볍고 빠른 Go 웹 프레임워크로 REST API 및 WebSocket 핸들링 |
| **Frontend** | **React (Vite)** | 컴포넌트 기반 UI 설계, Vite를 이용한 빠른 빌드 및 모듈 교체(HMR) |
| | **TypeScript** | 엄격한 타입 체크로 런타임 에러 방지 및 유지보수성 향상 |
| **Infrastructure** | **Docker Compose** | 다중 컨테이너(App, DB, Proxy)의 일관된 실행 환경 보장 |
| | **Terraform** | 클라우드 인프라 프로비저닝 자동화 및 형상 관리 |
| | **Nginx** | 리버스 프록시, 정적 파일 서빙, SSL 터미네이션 |
| **Database** | **PostgreSQL** | 관계형 데이터 모델링 및 트랜잭션 보장 |
| **DevOps** | **GitHub Actions** | CI/CD 파이프라인 구축 및 배포 프로세스 자동화 |

---

## 5. 트러블슈팅 및 성능 개선

### Case 1: WebSocket 연결 유지와 Sticky Session
- **문제**: 초기 K8s 환경에서 Ingress를 통해 로드 밸런싱을 할 때, 핸드셰이크 요청과 후속 소켓 연결이 서로 다른 파드(Pod)로 연결되어 끊기는 현상 발생.
- **해결**: Nginx Ingress Controller에 `ip_hash` 및 Sticky Session 어노테이션을 적용하여 클라이언트의 세션이 특정 파드에 고정되도록 네트워크 라우팅 정책을 수정함. VM 환경으로 이전 후에는 Nginx 설정을 통해 동일한 메커니즘 구현.

### Case 2: Docker 이미지 크기 최적화
- **문제**: 초기 Go 애플리케이션 빌드 이미지가 300MB를 초과하여 배포 시간이 오래 걸리고 대역폭 낭비 발생.
- **해결**: Docker **Multi-stage Build**를 적용. 빌드 단계(Golang 이미지)와 실행 단계(Alpine/Scratch 이미지)를 분리하여, 최종 이미지 크기를 **10MB 대**로 90% 이상 경량화 성공.

---
