output "vip_address" {
  description = "The virtual IP address"
  value       = huaweicloud_networking_vip.vip_1.ip_address
}

output "instance_0" {
  description = "The IP address of instance 0"
  value       = huaweicloud_compute_instance.mycompute[0].network[0].fixed_ip_v4
}

output "instance_1" {
  description = "The IP address of instance 1"
  value       = huaweicloud_compute_instance.mycompute[1].network[0].fixed_ip_v4
}
