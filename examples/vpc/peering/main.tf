data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  count = length(var.vpc_configurations)

  name                  = lookup(var.vpc_configurations[count.index], "vpc_name", null)
  cidr                  = lookup(var.vpc_configurations[count.index], "vpc_cidr", null)
  enterprise_project_id = lookup(var.vpc_configurations[count.index], "enterprise_project_id", null)
}

resource "huaweicloud_vpc_subnet" "test" {
  count = length(var.vpc_configurations)

  vpc_id     = huaweicloud_vpc.test[count.index].id
  name       = lookup(var.vpc_configurations[count.index], "subnet_name", null)
  cidr       = try(cidrsubnet(lookup(var.vpc_configurations[count.index], "vpc_cidr", null), 6, 32), null)
  gateway_ip = try(cidrhost(cidrsubnet(lookup(var.vpc_configurations[count.index], "vpc_cidr", null), 6, 32), 1), null)
}

resource "huaweicloud_vpc_peering_connection" "test" {
  count = length(var.vpc_configurations) == 2 ? 1 : 0

  name        = var.peering_connection_name
  vpc_id      = try(huaweicloud_vpc.test[0].id, null) # source VPC
  peer_vpc_id = try(huaweicloud_vpc.test[1].id, null) # target VPC
}

resource "huaweicloud_vpc_route" "test" {
  count = length(var.vpc_configurations)

  vpc_id      = huaweicloud_vpc.test[count.index].id
  destination = try(cidrsubnet(lookup(var.vpc_configurations[count.index], "vpc_cidr", null), 6, 33), null)
  type        = "peering"
  nexthop     = try(huaweicloud_vpc_peering_connection.test[0].id, null)
}
