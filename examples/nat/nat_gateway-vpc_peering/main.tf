data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_vpc" "vpcA" {
  name = "VPCA"
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnetA" {
  name              = "subnetA"
  cidr              = "192.168.1.0/24"
  gateway_ip        = var.subnetA_gateway_ip
  vpc_id            = huaweicloud_vpc.vpcA.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}

resource "huaweicloud_nat_gateway" "mygateway" {
  name        = var.gateway_name
  description = "test for terraform"
  spec        = "1"
  vpc_id      = huaweicloud_vpc.vpcA.id
  subnet_id   = huaweicloud_vpc_subnet.subnetA.id
}

resource "huaweicloud_vpc_eip" "eipA" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "bandwidth_A"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "eipC" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "bandwidth_C"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_snat_rule" "snatA" {
  nat_gateway_id = huaweicloud_nat_gateway.mygateway.id
  floating_ip_id = huaweicloud_vpc_eip.eipA.id
  subnet_id      = huaweicloud_vpc_subnet.subnetA.id
}

resource "huaweicloud_networking_secgroup" "mysecgroup" {
  name        = var.secgroup_name
  description = "This is a security group"
}

resource "huaweicloud_networking_secgroup_rule" "mysecgrouprule_1" {
  direction         = "ingress"
  description       = "This is a basic acc test"
  ethertype         = "IPv4"
  ports             = "22,80,443,3389"
  protocol          = "tcp"
  remote_ip_prefix  = var.remote_ip_prefix
  security_group_id = huaweicloud_networking_secgroup.mysecgroup.id
}

resource "huaweicloud_networking_secgroup_rule" "mysecgrouprule_2" {
  direction         = "ingress"
  description       = "This is a basic acc test"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.mysecgroup.id
}

data "huaweicloud_compute_flavors" "flavor" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "instanceA" {
  name               = "instanceA"
  image_name         = var.image_name
  flavor_id          = data.huaweicloud_compute_flavors.flavor.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.mysecgroup.id]
  admin_pass         = var.admin_pass
  availability_zone  = data.huaweicloud_availability_zones.myaz.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.subnetA.id
  }
}

resource "huaweicloud_nat_dnat_rule" "dnatA" {
  nat_gateway_id        = huaweicloud_nat_gateway.mygateway.id
  floating_ip_id        = huaweicloud_vpc_eip.eipC.id
  port_id               = huaweicloud_compute_instance.instanceA.network[0].port
  protocol              = "tcp"
  internal_service_port = var.internal_service_port
  external_service_port = var.external_service_port_A
}

resource "huaweicloud_vpc" "vpcB" {
  name = "VPCB"
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnetB" {
  name              = "subnetB"
  cidr              = "192.168.2.0/24"
  gateway_ip        = var.subnetB_gateway_ip
  vpc_id            = huaweicloud_vpc.vpcB.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}

resource "huaweicloud_vpc_eip" "eipB" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "bandwidth_B"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "eipD" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "bandwidth_D"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_snat_rule" "snatB" {
  nat_gateway_id = huaweicloud_nat_gateway.mygateway.id
  floating_ip_id = huaweicloud_vpc_eip.eipB.id
  source_type    = 1
  cidr           = "192.168.2.0/24"
}

resource "huaweicloud_compute_instance" "instanceB" {
  name               = "instanceB"
  image_name         = var.image_name
  flavor_id          = data.huaweicloud_compute_flavors.flavor.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.mysecgroup.id]
  admin_pass         = var.admin_pass
  availability_zone  = data.huaweicloud_availability_zones.myaz.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.subnetB.id
  }
}

resource "huaweicloud_nat_dnat_rule" "dnatB" {
  nat_gateway_id        = huaweicloud_nat_gateway.mygateway.id
  floating_ip_id        = huaweicloud_vpc_eip.eipD.id
  private_ip            = huaweicloud_compute_instance.instanceB.network[0].fixed_ip_v4
  protocol              = "tcp"
  internal_service_port = var.internal_service_port
  external_service_port = var.external_service_port_B
}

resource "huaweicloud_vpc_peering_connection" "peering" {
  name        = var.peer_conn_name
  vpc_id      = huaweicloud_vpc.vpcA.id
  peer_vpc_id = huaweicloud_vpc.vpcB.id
}

resource "huaweicloud_vpc_route" "vpc_route" {
  vpc_id      = huaweicloud_vpc.vpcA.id
  destination = "192.168.0.0/16"
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.peering.id
}

resource "huaweicloud_vpc_route" "vpc_route_1" {
  vpc_id      = huaweicloud_vpc.vpcB.id
  destination = "192.168.0.0/16"
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.peering.id
}

resource "huaweicloud_vpc_route" "vpc_route_2" {
  vpc_id      = huaweicloud_vpc.vpcB.id
  destination = "0.0.0.0/0"
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.peering.id
}
