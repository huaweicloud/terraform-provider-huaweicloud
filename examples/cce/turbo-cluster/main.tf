data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  count = var.vpc_id == "" && var.subnet_id == "" ? 1 : 0

  name = var.vpc_name
  cidr = var.vpc_cidr
}

# One for the turbo cluster and one for the ENI configuration
resource "huaweicloud_vpc_subnet" "test" {
  count = var.subnet_id == "" ? 1 : 0

  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  name              = var.subnet_name
  cidr              = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : var.subnet_cidr != "" ? cidrhost(var.subnet_cidr, 1) : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0), 1)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
}

# ST.001 Disable
resource "huaweicloud_vpc_subnet" "eni" {
# ST.001 Enable
  count = var.eni_ipv4_subnet_id == "" ? 1 : 0

  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  name              = var.eni_subnet_name
  cidr              = var.eni_subnet_cidr != "" ? var.eni_subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 1)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : var.eni_subnet_cidr != "" ? cidrhost(var.eni_subnet_cidr, 1) : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 1), 1)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
}

resource "huaweicloud_vpc_eip" "test" {
  count = var.eip_address == "" ? 1 : 0

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = var.cluster_name
  flavor_id              = var.cluster_flavor_id
  cluster_version        = var.cluster_version
  cluster_type           = var.cluster_type
  container_network_type = var.container_network_type
  vpc_id                 = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  subnet_id              = var.subnet_id != "" ? var.subnet_id : huaweicloud_vpc_subnet.test[0].id
  eni_subnet_id          = var.eni_ipv4_subnet_id != "" ? var.eni_ipv4_subnet_id : huaweicloud_vpc_subnet.eni[0].ipv4_subnet_id
  eip                    = var.eip_address != "" ? var.eip_address : huaweicloud_vpc_eip.test[0].address
  description            = var.cluster_description
  tags                   = var.cluster_tags
}
