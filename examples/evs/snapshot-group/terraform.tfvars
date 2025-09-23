vpc_name          = "evs-test-vpc"
subnet_name       = "evs-test-subnet"
secgroup_name     = "evs-test-sg"
ecs_instance_name = "evs-test-ecs"

volume_configuration = [
  {
    name        = "evs-test-volume1"
    size        = 50
    volume_type = "SSD"
    device_type = "VBD"
  },
  {
    name        = "evs-test-volume2"
    size        = 100
    volume_type = "SAS"
    device_type = "SCSI"
  }
]
