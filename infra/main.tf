# infra/main.tf
# 하드코딩 x / 변수화 필수!!

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
}