data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

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
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

# ST.001 Disable
resource "huaweicloud_vpc_subnet" "eni" {
  # ST.001 Enable
  count = var.eni_ipv4_subnet_id == "" ? 1 : 0

  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  name              = var.eni_subnet_name
  cidr              = var.eni_subnet_cidr != "" ? var.eni_subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 1)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : var.eni_subnet_cidr != "" ? cidrhost(var.eni_subnet_cidr, 1) : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 1), 1)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_cce_cluster" "test" {
  name                         = var.cluster_name
  flavor_id                    = var.cluster_flavor_id
  cluster_version              = var.cluster_version
  cluster_type                 = var.cluster_type
  container_network_type       = var.container_network_type
  vpc_id                       = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  enable_distribute_management = true
  subnet_id                    = var.subnet_id != "" ? var.subnet_id : huaweicloud_vpc_subnet.test[0].id
  eni_subnet_id                = var.eni_ipv4_subnet_id != "" ? var.eni_ipv4_subnet_id : huaweicloud_vpc_subnet.eni[0].ipv4_subnet_id
  description                  = var.cluster_description
  tags                         = var.cluster_tags
}

data "huaweicloud_compute_flavors" "test" {
  count = var.node_flavor_id == "" ? 1 : 0

  performance_type  = var.node_flavor_performance_type
  cpu_core_count    = var.node_flavor_cpu_core_count
  memory_size       = var.node_flavor_memory_size
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_cce_partition" "test" {
  count = var.node_partition == "" ? 1 : 0

  cluster_id           = huaweicloud_cce_cluster.test.id
  name                 = var.partition_name
  category             = var.partition_category
  public_border_group  = var.partition_public_border_group
  partition_subnet_id  = huaweicloud_vpc_subnet.eni[0].id
  container_subnet_ids = [huaweicloud_vpc_subnet.eni[0].ipv4_subnet_id]
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = var.node_name
  flavor_id         = var.node_flavor_id != "" ? var.node_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  password          = var.node_password
  partition         = var.node_partition != "" ? var.node_partition : huaweicloud_cce_partition.test[0].id

  root_volume {
    volumetype = var.root_volume_type
    size       = var.root_volume_size
  }

  dynamic "data_volumes" {
    for_each = var.data_volumes_configuration

    content {
      volumetype = data_volumes.value.volumetype
      size       = data_volumes.value.size
    }
  }

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
    ]
  }
}

resource "huaweicloud_cce_node_pool" "test" {
  count = var.node_pool_name != "" ? 1 : 0

  cluster_id         = huaweicloud_cce_cluster.test.id
  name               = var.node_pool_name
  os                 = var.node_pool_os_type
  flavor_id          = var.node_flavor_id != "" ? var.node_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  initial_node_count = var.node_pool_initial_node_count
  availability_zone  = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  password           = var.node_pool_password
  type               = "vm"
  partition          = var.node_partition != "" ? var.node_partition : huaweicloud_cce_partition.test[0].id

  root_volume {
    volumetype = var.root_volume_type
    size       = var.root_volume_size
  }

  dynamic "data_volumes" {
    for_each = var.data_volumes_configuration

    content {
      volumetype = data_volumes.value.volumetype
      size       = data_volumes.value.size
    }
  }

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
    ]
  }
}
