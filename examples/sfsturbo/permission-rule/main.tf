data "huaweicloud_availability_zones" "test" {}

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

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = var.turbo_availability_zone != "" ? var.turbo_availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
  name              = var.turbo_name
  size              = var.turbo_size
  share_proto       = var.turbo_share_proto
  share_type        = var.turbo_share_type
  hpc_bandwidth     = var.turbo_hpc_bandwidth
}

resource "huaweicloud_sfs_turbo_perm_rule" "test" {
  share_id  = huaweicloud_sfs_turbo.test.id
  ip_cidr   = var.rule_ip_cidr
  rw_type   = var.rule_rw_type
  user_type = var.rule_user_type
}
