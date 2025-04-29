# List the availability zones in the current region.
data "huaweicloud_availability_zones" "test" {}

# Create a postpaid EVS volume.
resource "huaweicloud_evs_volume" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "test-name"
  size              = 100
  description       = "test description"
  volume_type       = "GPSSD2"
  iops              = 3000
  throughput        = 125
  device_type       = "SCSI"
  multiattach       = false

  tags = {
    foo = "bar"
    key = "value"
  }
}