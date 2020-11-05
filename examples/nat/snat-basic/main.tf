resource "huaweicloud_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id      = huaweicloud_vpc.vpc_1.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "huaweicloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_gateway" "nat_1" {
  name                = var.nat_gatway_name
  description         = "test for terraform examples"
  spec                = "1"
  internal_network_id = huaweicloud_vpc_subnet.subnet_1.id
  router_id           = huaweicloud_vpc.vpc_1.id
}

resource "huaweicloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = huaweicloud_nat_gateway.nat_1.id
  network_id     = huaweicloud_vpc_subnet.subnet_1.id
  floating_ip_id = huaweicloud_vpc_eip.eip_1.id
}
