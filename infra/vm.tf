resource "google_compute_instance" "app_server" {
  name                      = "realtime-chat-server"
  machine_type              = "e2-micro"
  zone                      = var.gcp_zone
  allow_stopping_for_update = true

  # 부팅 디스크 설정
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts" # Ubuntu 22.04 LTS
      size  = 30
      type  = "pd-standard"
    }
  }

  # 네트워크 설정
  network_interface {
    network    = google_compute_network.vpc_network.id
    subnetwork = google_compute_subnetwork.subnet.id
    access_config {
      # 예약한 고정 IP 연결
      nat_ip = google_compute_address.vm_static_ip.address
    }
  }
  # VM에 기본 서비스 계정을 연결합니다.
  service_account {
    # email을 지정하지 않거나 "default"로 설정하면 프로젝트의 기본 Compute Engine 서비스 계정을 사용합니다.
    email  = ""
    scopes = ["cloud-platform"] # 모든 GCP API에 접근 가능한 범위를 부여 (실제 권한은 IAM으로 제어)
  }

  # 태그 설정
  tags = ["http-server", "https-server"]

  # [핵심] 자동화 스크립트
  metadata_startup_script = <<-EOF
    #! /bin/bash
    exec > /var/log/startup-script.log 2>&1
    echo "Start setup..."

    # 1. Swap 메모리 설정 (2GB)
    # 1GB RAM으로는 빌드/배포 시 멈출 수 있으므로 필수입니다.
    if [ ! -f /swapfile ]; then
        echo "Creating 2GB swapfile..."
        fallocate -l 2G /swapfile
        chmod 600 /swapfile
        mkswap /swapfile
        swapon /swapfile
        echo '/swapfile none swap sw 0 0' >> /etc/fstab
        echo "Swap created."
    else
        echo "Swap already exists."
    fi

    # 2. Docker 설치
    if ! command -v docker &> /dev/null; then
        echo "Installing Docker..."
        apt-get update
        apt-get install -y ca-certificates curl gnupg lsb-release
        mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        echo \
          "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
          $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
        echo "Docker installed."
    else
        echo "Docker already installed."
    fi

    echo "Setup finished!"
  EOF
}
