# infra/variables.tf

variable "gcp_project_id" {
  description = "The GCP project ID to deploy resources into."
  type        = string
}

variable "gcp_region" {
  description = "The GCP region for all resources."
  type        = string
  default     = "asia-northeast3" # 기본값 설정
}

variable "gcp_zone" {
  description = "The GCP zone for all resources."
  type        = string
  default     = "asia-northeast3-a" # 기본값 설정
}