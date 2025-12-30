data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = var.vpn_gateway_az_flavor
  attachment_type = var.vpn_gateway_az_attachment_type
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

  dynamic "publicip" {
    for_each = var.vpc_eip_public_ip

    content {
      type = publicip.value.type
    }
  }

  dynamic "bandwidth" {
    for_each = var.vpc_eip_bandwidth

    content {
      name        = bandwidth.value.name
      size        = bandwidth.value.size
      share_type  = bandwidth.value.share_type
      charge_mode = bandwidth.value.charge_mode
    }
  }
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = var.vpn_gateway_name
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    try(data.huaweicloud_vpn_gateway_availability_zones.test.names[0], null),
    try(data.huaweicloud_vpn_gateway_availability_zones.test.names[1], null)
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test[0].id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test[1].id
  }
}

resource "huaweicloud_vpn_customer_gateway" "test" {
  name     = var.vpn_customer_gateway_name
  id_value = var.vpn_customer_gateway_id_value
}

resource "huaweicloud_vpn_connection" "test" {
  name                = var.vpn_connection_name
  gateway_id          = huaweicloud_vpn_gateway.test.id
  gateway_ip          = huaweicloud_vpn_gateway.test.master_eip[0].id
  customer_gateway_id = huaweicloud_vpn_customer_gateway.test.id
  peer_subnets        = var.vpn_connection_peer_subnets
  vpn_type            = var.vpn_connection_vpn_type
  psk                 = var.vpn_connection_psk
  enable_nqa          = var.vpn_connection_enable_nqa
}
