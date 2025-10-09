security_group_name    = "tf_test_secgroup_demo"
keypair_name           = "tf_test_keypair_demo"
configuration_name     = "tf_test_as_configuration"
configuration_metadata = {
  some_key = "some_value"
}

configuration_user_data = <<EOT
# !/bin/sh
echo "Hello World! The time is now $(date -R)!" | tee /root/output.txt
EOT

configuration_disks = [
  {
    size        = 40
    volume_type = "SSD"
    disk_type   = "SYS"
  }
]

configuration_public_eip_settings = [
  {
    ip_type   = "5_bgp"
    bandwidth = {
      size          = 10
      share_type    = "PER"
      charging_mode = "traffic"
    }
  }
]
