# Create flow log with ER instance
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

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = var.er_vpc_attachment_name
  auto_create_vpc_routes = var.er_vpc_attachment_auto_create_vpc_routes
}

resource "huaweicloud_lts_group" "test" {
  group_name  = var.lts_group_name
  ttl_in_days = var.lts_group_ttl_in_days
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = var.lts_stream_name
}

resource "huaweicloud_er_flow_log" "test" {
  name           = var.er_flow_log_name
  instance_id    = huaweicloud_er_instance.test.id
  log_store_type = var.er_flow_log_store_type
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  resource_type  = var.er_flow_log_resource_type
  resource_id    = huaweicloud_er_vpc_attachment.test.id
}
