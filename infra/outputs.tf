# infra/outputs.tf

output "vm_ip_address" {
  value       = google_compute_address.vm_static_ip.address
  description = "The static IP address of the VM"
}