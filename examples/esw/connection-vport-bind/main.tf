data "huaweicloud_esw_flavors" "test" {}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = 2

  vpc_id     = huaweicloud_vpc.test.id
  name       = format("%s-%d", var.subnet_name, count.index)
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
    virsubnet_id = huaweicloud_vpc_subnet.test[0].id
    tunnel_ip    = var.esw_instance_tunnel_ip
  }

  charge_infos {
    charge_mode = "postPaid"
  }
}

resource "huaweicloud_esw_connection" "test" {
  instance_id  = huaweicloud_esw_instance.test.id
  name         = var.esw_connection_name
  vpc_id       = huaweicloud_vpc.test.id
  virsubnet_id = huaweicloud_vpc_subnet.test[1].id
  fixed_ips    = var.esw_connection_fixed_ips

  remote_infos {
    segmentation_id = var.esw_connection_segmentation_id
    tunnel_ip       = huaweicloud_esw_instance.test.tunnel_info[0].tunnel_ip
    tunnel_port     = huaweicloud_esw_instance.test.tunnel_info[0].tunnel_port
  }
}

resource "huaweicloud_vpc_subnet_private_ip" "test" {
  subnet_id    = huaweicloud_vpc_subnet.test[1].id
  device_owner = "neutron:VIP_PORT"
}

resource "huaweicloud_esw_connection_vport_bind" "test" {
  connection_id = huaweicloud_esw_connection.test.id
  port_id       = huaweicloud_vpc_subnet_private_ip.test.id
}
