data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_vpc" "vpcA" {
  name = "VPCA"
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnetA" {
  name              = "subnetA"
  cidr              = var.subnetA_cidr
  gateway_ip        = var.subnetA_gateway_ip
  vpc_id            = huaweicloud_vpc.vpcA.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}

resource "huaweicloud_vpc" "vpcB" {
  name = "VPCB"
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnetB" {
  name              = "subnetB"
  cidr              = var.subnetB_cidr
  gateway_ip        = var.subnetB_gateway_ip
  vpc_id            = huaweicloud_vpc.vpcB.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}

resource "huaweicloud_vpc_route" "vpc_route" {
  vpc_id      = huaweicloud_vpc.vpcA.id
  destination = var.vpc_cidr
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.peering.id
}

resource "huaweicloud_vpc_route" "vpc_route_1" {
  vpc_id      = huaweicloud_vpc.vpcB.id
  destination = var.vpc_cidr
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.peering.id
}

resource "huaweicloud_vpc_route" "vpc_route_2" {
  vpc_id      = huaweicloud_vpc.vpcB.id
  destination = "0.0.0.0/0"
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.peering.id
}

resource "huaweicloud_vpc_peering_connection" "peering" {
  name        = var.peer_conn_name
  vpc_id      = huaweicloud_vpc.vpcA.id
  peer_vpc_id = huaweicloud_vpc.vpcB.id
}
