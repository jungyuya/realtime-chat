resource "google_container_cluster" "primary" {
  name     = "realtime-chat-cluster"
  location = var.gcp_region

  network    = google_compute_network.vpc_network.name
  subnetwork = google_compute_subnetwork.subnet.name

  remove_default_node_pool = true
  initial_node_count       = 1

  # [추가] 기본 노드 풀의 설정을 직접 제어합니다.
  # 이 블록을 추가하면, Terraform이 임시로 생성하는 기본 노드 풀조차도
  # 아래의 설정을 따르게 됩니다.
  node_config {
    disk_size_gb = 30
    machine_type = "e2-medium"
  }

  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"
}

# ------------------------------------------------------------------------------
# GKE 커스텀 노드 풀
# ------------------------------------------------------------------------------
# 이 리소스는 그대로 둡니다.
# 클러스터가 생성된 후, 이 커스텀 노드 풀이 최종적으로 사용됩니다.
resource "google_container_node_pool" "primary_nodes" {
  name     = "default-pool"
  cluster  = google_container_cluster.primary.name
  location = var.gcp_region

  initial_node_count = 1 # initial_node_count로 변경하는 것이 더 명확합니다.

  node_config {
    machine_type = "e2-medium"
    disk_size_gb = 30

    service_account = "default"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }

  # [추가] 클러스터 리소스가 먼저 생성된 후에 이 노드 풀이 생성되도록
  # 명시적인 의존성을 추가하여 안정성을 높입니다.
  depends_on = [
    google_container_cluster.primary
  ]
}