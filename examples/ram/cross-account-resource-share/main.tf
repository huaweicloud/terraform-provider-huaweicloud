# Create resources for sharing
data "huaweicloud_account" "test" {}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.secgroup_name
  description          = var.secgroup_description
  delete_default_rules = var.secgroup_delete_default_rules
}

locals {
  urn_prefix = "urn:hws:ecs:${var.region_name}:${data.huaweicloud_account.test.id}"

  # Use provided resource_urns if specified, otherwise use created resources
  resource_urns = length(var.resource_urns) > 0 ? var.resource_urns : [
    "${local.urn_prefix}:vpc:${huaweicloud_vpc.test.id}",
    "${local.urn_prefix}:subnet:${huaweicloud_vpc_subnet.test.id}",
    "${local.urn_prefix}:security-group:${huaweicloud_networking_secgroup.test.id}"
  ]
}

# Create a resource share instance
resource "huaweicloud_ram_resource_share" "test" {
  name           = var.resource_share_name
  description    = var.description
  principals     = var.principals
  resource_urns  = local.resource_urns
  permission_ids = var.permission_ids
}
