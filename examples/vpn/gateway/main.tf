# Create a VPN gateway
data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = var.vpn_gateway_flavor
  attachment_type = var.vpn_gateway_attachment_type
}

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

resource "huaweicloud_vpc_eip" "test" {
  count = 2

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = "${var.bandwidth_name}-${count.index}"
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = var.vpn_gateway_name
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    try(data.huaweicloud_vpn_gateway_availability_zones.test.names[0], "default_value"),
    try(data.huaweicloud_vpn_gateway_availability_zones.test.names[1], "default_value")
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test[0].id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test[1].id
  }

  delete_eip_on_termination = var.vpn_gateway_delete_eip_on_termination
}
