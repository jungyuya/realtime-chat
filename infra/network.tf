# ------------------------------------------------------------------------------
# VPC 네트워크
# ------------------------------------------------------------------------------
resource "google_compute_network" "vpc_network" {
  # 리소스의 이름 (Terraform 코드 내에서 사용할 이름)
  name                    = "realtime-chat-vpc"
  # auto_create_subnetworks를 false로 설정하여, 우리가 직접 서브넷을 제어합니다.
  auto_create_subnetworks = false
}

# ------------------------------------------------------------------------------
# 서브넷
# ------------------------------------------------------------------------------
resource "google_compute_subnetwork" "subnet" {
  name          = "realtime-chat-subnet"
  # 이 서브넷이 속할 VPC 네트워크를 지정합니다.
  # google_compute_network.vpc_network.id는 위에서 만든 VPC의 ID를 참조합니다.
  network       = google_compute_network.vpc_network.id
  # 이 서브넷이 사용할 IP 주소 범위를 지정합니다. (사설 IP 대역)
  ip_cidr_range = "10.10.10.0/24"
  # 이 서브넷이 위치할 리전을 지정합니다.
  # var.gcp_region은 변수로 관리할 수 있지만, 여기서는 main.tf의 값을 따릅니다.
  region        = "asia-northeast3"
}

# ------------------------------------------------------------------------------
# 방화벽 규칙 (Firewall Rules)
# ------------------------------------------------------------------------------
# 외부 인터넷에서 들어오는 SSH, HTTP, HTTPS 트래픽을 허용하는 규칙
resource "google_compute_firewall" "allow_http_https_ssh" {
  name    = "allow-http-https-ssh"
  # 이 방화벽 규칙이 적용될 네트워크를 지정합니다.
  network = google_compute_network.vpc_network.name

  # "ingress"는 들어오는 트래픽에 대한 규칙임을 의미합니다.
  direction = "INGRESS"

  # 어떤 프로토콜과 포트를 허용할지 정의합니다.
  allow {
    protocol = "tcp"
    ports    = ["22", "80", "443"] # SSH, HTTP, HTTPS
  }

  # 어떤 IP 주소로부터의 요청을 허용할지 지정합니다.
  # "0.0.0.0/0"은 "인터넷의 모든 IP"를 의미합니다.
  source_ranges = ["0.0.0.0/0"]
}

# GKE 클러스터의 노드들 간의 통신을 허용하는 규칙 (매우 중요)
# GKE는 내부적으로 노드들이 서로 통신해야 제대로 동작합니다.
resource "google_compute_firewall" "allow_internal_gke" {
  name    = "allow-internal-gke"
  network = google_compute_network.vpc_network.name
  
  allow {
    protocol = "tcp"
    ports    = ["0-65535"] # 모든 TCP 포트
  }
  allow {
    protocol = "udp"
    ports    = ["0-65535"] # 모든 UDP 포트
  }
  allow {
    protocol = "icmp" # Ping과 같은 통신을 위해 필요
  }

  # source_ranges 대신 source_tags를 사용하여,
  # 특정 태그가 붙은 VM들 간의 통신만 허용할 수 있습니다.
  # 여기서는 간단하게 동일 서브넷 내의 모든 통신을 허용합니다.
  source_ranges = [google_compute_subnetwork.subnet.ip_cidr_range]
}