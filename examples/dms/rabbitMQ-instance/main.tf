# Create a VPC.
resource "huaweicloud_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = "192.168.0.0/24"
}

# Create a subnet under the VPC that created above.
resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = var.subnet_name
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

# Create a security group.
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = var.security_group_name
  description = "terraform security group"
}

# List the availability zones in the current region.
data "huaweicloud_availability_zones" "zones" {}

# Query flavor information based on flavorID and storage I/O specification.
# Make sure the flavors are available in the availability zone.
data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type              = "cluster"
  flavor_id         = "c6.2u4g.cluster"
  storage_spec_code = "dms.physical.storage.ultra.v2"

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
  ]
}

# Create the DMS RabbitMQ instance.
resource "huaweicloud_dms_rabbitmq_instance" "instance_1" {
  name        = "instance_1"
  description = "rabbitmq test"

  access_user = "user"
  password    = "Rabbitmqtest@123"

  vpc_id            = huaweicloud_vpc.vpc_1.id
  network_id        = huaweicloud_vpc_subnet.vpc_subnet_1.id
  security_group_id = huaweicloud_networking_secgroup.secgroup.id

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
  ]

  flavor_id         = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0].id
  engine_version    = "3.8.35"
  broker_num        = 3
  storage_space     = 600
  storage_spec_code = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0].ios[0].storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
