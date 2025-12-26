data "huaweicloud_esw_flavors" "test" {}

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

resource "huaweicloud_esw_instance" "test" {
  name        = var.esw_instance_name
  flavor_ref  = try(data.huaweicloud_esw_flavors.test.flavors[0].name, "")
  ha_mode     = var.esw_instance_ha_mode
  description = var.esw_instance_description

  availability_zones {
    primary = try(data.huaweicloud_esw_flavors.test.flavors.0.available_zones[0], "")
    standby = try(data.huaweicloud_esw_flavors.test.flavors.0.available_zones[1], "")
  }

  tunnel_info {
    vpc_id       = huaweicloud_vpc.test.id
    virsubnet_id = huaweicloud_vpc_subnet.test.id
    tunnel_ip    = var.esw_instance_tunnel_ip
  }

  charge_infos {
    charge_mode = "postPaid"
  }
}
