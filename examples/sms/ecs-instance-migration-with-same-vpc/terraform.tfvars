vpc_name                    = "tf_test_vpc"
subnet_name                 = "tf_test_subnet"
security_group_name         = "tf_test_security_group"
instance_name               = "tf_test_source_server"
instance_admin_password     = "YourPassword123!"
destination_instance_name   = "tf_test_destination_server"
source_server_name          = "tf_test_source_server"
source_server_os_version    = "UBUNTU_24_4_64BIT"
source_server_agent_version = "25.2.0"

source_server_disks = {
  name            = "Disk 0"
  partition_style = "MBR"
  device_use      = "BOOT"
  size            = 42949672960
  used_size       = 42949672960

  physical_volumes = [
    {
      device_use  = "OS"
      file_system = "ext4"
      mount_point = "/"
      name        = "/dev/vda1"
      size        = 42943137792
      used_size   = 0
    }
  ]
}

migrate_task_type        = "MIGRATE_FILE"
task_target_server_disks = {
  name        = "/dev/sda"
  size        = 85899345920
  device_type = "NORMAL"

  physical_volumes = [
    {
      device_type  = "OS"
      file_system  = "ext4"
      mount_point  = "/"
      name         = "/dev/sda1"
      size         = 85899345920
      volume_index = 0
    }
  ]
}
