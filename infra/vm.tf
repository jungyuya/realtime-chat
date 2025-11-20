# ------------------------------------------------------------------------------
# Compute Engine VM 인스턴스 (DB 서버용)
# ------------------------------------------------------------------------------
resource "google_compute_instance" "db_server" {
  name         = "db-server"
  machine_type = "e2-micro"
  zone         = var.gcp_zone

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network    = google_compute_network.vpc_network.id
    subnetwork = google_compute_subnetwork.subnet.id
    access_config {
      // Ephemeral public IP
    }
  }

  tags = ["http-server", "https-server"]

  # VM이 시작될 때 실행할 자동화 스크립트
  metadata_startup_script = <<-EOF
    #! /bin/bash
    
    # 1. 로그 파일 설정 (디버깅용)
    exec > /var/log/startup-script.log 2>&1
    echo "Start setup..."

    # 2. Docker 설치 (공식 문서 기준)
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl gnupg
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg

    echo \
      "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
      "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # 3. PostgreSQL 컨테이너 실행
    # 이미 실행 중인지 확인하고, 없으면 실행합니다 (재시작 시 중복 방지)
    if [ ! "$(sudo docker ps -q -f name=postgres-db)" ]; then
        if [ "$(sudo docker ps -aq -f name=postgres-db)" ]; then
            # 컨테이너가 중지된 상태라면 재시작
            sudo docker start postgres-db
        else
            # 컨테이너가 아예 없으면 새로 생성 및 실행
            sudo docker run -d \
              --name postgres-db \
              -p 5432:5432 \
              -e POSTGRES_USER=postgres \
              -e POSTGRES_PASSWORD=mysecretpassword \
              -e POSTGRES_DB=chatdb \
              -v postgres-data:/var/lib/postgresql/data \
              postgres:15-alpine
        fi
    fi

    echo "Setup finished successfully!"
  EOF
}