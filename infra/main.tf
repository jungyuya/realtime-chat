# infra/main.tf
provider "google" {
  credentials = file("gcp-credentials.json")
  project = var.gcp_project_id 
  region  = var.gcp_region
  zone    = var.gcp_zone
}

# ------------------------------------------------------------------------------
# Artifact Registry (Docker 이미지 저장소)
# ------------------------------------------------------------------------------
resource "google_artifact_registry_repository" "docker_repo" {
  repository_id = "realtime-chat-repo"
  location      = var.gcp_region
  format        = "DOCKER"
  description   = "Docker repository for realtime-chat project"

  # [추가] 수명주기 정책 설정
  # dry_run = false: 실제로 삭제를 수행함 (true면 로그만 남김)
  cleanup_policy_dry_run = false

  # 정책 1: 'latest' 태그가 붙은 이미지는 절대 삭제하지 않음 (보호)
  cleanup_policies {
    id     = "keep-latest"
    action = "KEEP"
    condition {
      tag_state = "TAGGED"
      tag_prefixes = ["latest"]
    }
  }

  # 정책 2: 최근 업로드된 10개 버전은 유지함 (보호)
  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"
    condition {
      # 최근 10개 버전 유지 (가장 최근에 업로드된 순서대로)
      newer_than = "2592000s" # 30일 
    }
  }
  
  # 정책 3: 위 KEEP 조건에 해당하지 않는 모든 이미지는 삭제
  cleanup_policies {
    id     = "delete-old-images"
    action = "DELETE"
    condition {
      older_than = "2592000s" # 30일 이상 된 이미지는 삭제 대상
    }
  }
}

# ------------------------------------------------------------------------------
# Global Static IP (Ingress용 고정 IP)
# ------------------------------------------------------------------------------
# GKE Ingress가 사용할 전역 고정 IP 주소를 예약합니다.
resource "google_compute_global_address" "ingress_ip" {
  name = "realtime-chat-ingress-ip"
  # GKE Ingress용 IP는 반드시 'EXTERNAL' 타입 (기본값)
}