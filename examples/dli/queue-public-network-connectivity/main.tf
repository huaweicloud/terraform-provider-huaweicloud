resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = var.elastic_resource_pool_name
  description           = var.elastic_resource_pool_description
  min_cu                = var.elastic_resource_pool_min_cu
  max_cu                = var.elastic_resource_pool_max_cu
  cidr                  = var.elastic_resource_pool_cidr
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  label = {
    spec = "basic"
  }
}

resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1

  name                  = var.queue_name
  queue_type            = var.queue_type
  cu_count              = var.queue_cu_count
  description           = var.queue_description
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc" "test" {
  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_dli_datasource_connection" "test" {
  name      = var.datasource_connection_name
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id

  dynamic "hosts" {
    for_each = var.datasource_connection_hosts

    content {
      name = hosts.value.name
      ip   = hosts.value.ip
    }
  }

  dynamic "routes" {
    for_each = var.datasource_connection_routes

    content {
      name = routes.value.name
      cidr = routes.value.cidr
    }
  }
}

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = huaweicloud_dli_datasource_connection.test.id
  elastic_resource_pools = [huaweicloud_dli_elastic_resource_pool.test.name]

  depends_on = [huaweicloud_dli_queue.test]
}

resource "huaweicloud_vpc_eip" "test" {
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.eip_bandwidth_name
    size        = var.eip_bandwidth_size
    share_type  = var.eip_bandwidth_share_type
    charge_mode = var.eip_bandwidth_charge_mode
  }
}

resource "huaweicloud_nat_gateway" "test" {
  name                  = var.nat_gateway_name
  spec                  = var.nat_gateway_spec
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  description           = var.nat_gateway_description
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  floating_ip_id = huaweicloud_vpc_eip.test.id
  source_type    = 1
  cidr           = huaweicloud_dli_elastic_resource_pool.test.cidr
  description    = var.snat_description
}
