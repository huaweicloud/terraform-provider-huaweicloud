output "vip_address" {
  value = huaweicloud_networking_vip.vip_1.ip_address
}

output "instance_0" {
  value = huaweicloud_compute_instance.mycompute[0].network.0.fixed_ip_v4
}

output "instance_1" {
  value = huaweicloud_compute_instance.mycompute[1].network.0.fixed_ip_v4
}
