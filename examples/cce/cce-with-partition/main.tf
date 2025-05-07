resource "huaweicloud_vpc" "test" {
  name = var.random_resource_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = var.random_resource_name
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)

  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
}

resource "huaweicloud_vpc_subnet" "eni_network" {
  vpc_id = huaweicloud_vpc.test.id

  name       = format("%s_eni_usage", var.random_resource_name)
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 2)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 2), 1)

  availability_zone = var.iec_availability_zone
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = var.random_resource_name
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "eni"

  enable_distribute_management = true

  eni_subnet_id = join(",", [
    huaweicloud_vpc_subnet.test.ipv4_subnet_id,
  ])

  lifecycle {
    ignore_changes = [
      eni_subnet_id,
    ]
  }
}

resource "huaweicloud_cce_partition" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id

  name                 = var.iec_availability_zone
  category             = "HomeZone"
  public_border_group  = var.iec_partition_border_group
  partition_subnet_id  = huaweicloud_vpc_subnet.eni_network.id
  container_subnet_ids = [huaweicloud_vpc_subnet.eni_network.ipv4_subnet_id]
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = var.iec_availability_zone
  performance_type  = "computingv3"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = var.random_resource_name
  flavor_id         = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  availability_zone = var.iec_availability_zone
  partition         = huaweicloud_cce_partition.test.id
  password          = "Overlord!!52259"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = var.random_resource_name
  flavor_id                = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  initial_node_count       = 1
  availability_zone        = var.iec_availability_zone
  password                 = "Overlord!!52259"
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  partition                = huaweicloud_cce_partition.test.id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
