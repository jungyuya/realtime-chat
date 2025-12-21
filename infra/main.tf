provider "google" {
  credentials = file("gcp-credentials.json")
  project     = var.gcp_project_id
  region      = var.gcp_region
  zone        = var.gcp_zone
}

# VM용 고정 IP 예약 (Regional)
resource "google_compute_address" "vm_static_ip" {
  name   = "realtime-chat-vm-ip"
  region = var.gcp_region
}

# Artifact Registry (이미지 저장소)
# 리전이 us-central1으로 변경되어 새로 생성됩니다.
resource "google_artifact_registry_repository" "docker_repo" {
  repository_id = "realtime-chat-repo"
  location      = var.gcp_region
  format        = "DOCKER"
  description   = "Docker repository for realtime-chat project"

  cleanup_policy_dry_run = false

  # 최근 10개 버전 유지 정책
  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"
    condition {
      newer_than = "2592000s" # 30일
    }
  }
}