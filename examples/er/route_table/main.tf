# Create a route table with ER instance
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name               = var.er_instance_name
  asn                = var.er_instance_asn
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = var.er_route_table_name
}
