# VPC 네트워크
resource "google_compute_network" "vpc_network" {
  name                    = "realtime-chat-vpc"
  auto_create_subnetworks = false
}

# 서브넷
resource "google_compute_subnetwork" "subnet" {
  name          = "realtime-chat-subnet"
  network       = google_compute_network.vpc_network.id
  ip_cidr_range = "10.10.10.0/24"
  region        = var.gcp_region
}

# 방화벽 규칙 (SSH + HTTP + HTTPS)
resource "google_compute_firewall" "allow_web_ssh" {
  name    = "allow-web-ssh"
  network = google_compute_network.vpc_network.name
  
  allow {
    protocol = "tcp"
    ports    = ["22", "80", "443"]
  }

  # 전 세계 접속 허용 (SSH Key로 보안 유지)
  source_ranges = ["0.0.0.0/0"]
}