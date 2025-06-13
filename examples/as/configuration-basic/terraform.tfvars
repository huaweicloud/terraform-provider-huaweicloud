# VPC
security_group_name = "tf_test_secgroup_demo"

# KPS
kps_key_pair_name = "tf_test_keypair_demo"
kps_public_key    = "your_public_key" # Please replace with your actual public key

# AS
as_configuration_name = "tf_test_as_configuration"
as_metadata = {
  some_key = "some_value"
}
as_user_data = <<EOT
#!/bin/sh
echo "Hello World! The time is now $(date -R)!" | tee /root/output.txt
EOT

as_disks = [
  {
    size        = 40
    volume_type = "SSD"
    disk_type   = "SYS"
  }
]

as_public_ip = [
  {
    eip = {
      ip_type = "5_bgp"
      bandwidth = {
        size          = 10
        share_type    = "PER"
        charging_mode = "traffic"
      }
    }
  }
]
