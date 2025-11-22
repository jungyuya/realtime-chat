# 생성된 Ingress용 고정 IP 주소를 출력합니다.
output "ingress_ip_address" {
  value       = google_compute_global_address.ingress_ip.address
  description = "The global static IP address for the GKE Ingress"
}