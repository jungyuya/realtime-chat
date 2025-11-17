# infra/main.tf
# 하드코딩 x / 변수화 필수!!

provider "google" {
  credentials = file("gcp-credentials.json")
  project = var.gcp_project_id 
  region  = var.gcp_region
  zone    = var.gcp_zone
}

