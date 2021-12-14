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

# Find the product ID according to the Kafka instance information to be created.
data "huaweicloud_dms_product" "product_1" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  bandwidth         = "100MB"
  storage_spec_code = "dms.physical.storage.high"
}

# Create the DMS Kafka instance.
resource "huaweicloud_dms_kafka_instance" "kafka_instance_1" {
  name        = "instance_1"
  description = "kafka instance demo"

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1],
    data.huaweicloud_availability_zones.zones.names[2],
  ]

  product_id        = data.huaweicloud_dms_product.product_1.id
  engine_version    = "2.3.0"
  storage_spec_code = data.huaweicloud_dms_product.product_1.storage_spec_code

  vpc_id            = huaweicloud_vpc.vpc_1.id
  network_id        = huaweicloud_vpc_subnet.vpc_subnet_1.id
  security_group_id = huaweicloud_networking_secgroup.secgroup.id

  access_user      = var.access_user_name
  password         = var.access_user_password
  manager_user     = var.manager_user
  manager_password = var.manager_password
}
