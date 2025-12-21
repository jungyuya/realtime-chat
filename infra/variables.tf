variable "gcp_project_id" {
  description = "The GCP project ID to deploy resources into."
  type        = string
}

variable "gcp_region" {
  description = "The GCP region for all resources."
  type        = string
  default     = "us-central1" # [중요] 평생 무료 등급 리전
}

variable "gcp_zone" {
  description = "The GCP zone for all resources."
  type        = string
  default     = "us-central1-f" # 해당 리전의 영역
}