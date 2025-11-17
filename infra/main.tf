# Terraform 자체에 대한 설정을 정의합니다.
terraform {
  # 이 코드는 Google Cloud Provider의 최신 버전을 사용하라고 지시합니다.
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# Google Cloud Provider에 대한 설정을 정의합니다.
provider "google" {
  # 사용할 인증 정보 파일의 경로를 지정합니다.
  credentials = file("gcp-credentials.json")

  # 리소스를 생성할 프로젝트 ID와 리전(Region), 영역(Zone)을 지정합니다.
  # TODO: "your-gcp-project-id" 부분을 당신의 실제 GCP 프로젝트 ID로 바꿔야 합니다.
  project = "your-gcp-project-id"
  region  = "asia-northeast3" # 서울 리전
  zone    = "asia-northeast3-a" # 서울 리전의 a 영역
}