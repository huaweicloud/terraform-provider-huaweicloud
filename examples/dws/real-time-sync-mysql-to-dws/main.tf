data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
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

resource "huaweicloud_networking_secgroup" "test" {
  name                  = var.security_group_name
  delete_default_rules  = var.security_group_delete_default_rules
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

# DWS and RDS ports are opened to DLI elastic resource pool
resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "${var.dws_port},${var.rds_db_port}"
  remote_ip_prefix  = var.elastic_resource_pool_cidr
}

data "huaweicloud_rds_flavors" "test" {
  count = var.rds_flavor_id == "" ? 1 : 0

  db_type           = "MySQL"
  db_version        = var.rds_db_version
  instance_mode     = var.rds_instance_mode
  vcpus             = var.rds_flavor_vcpus
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_rds_instance" "test" {
  name              = var.rds_instance_name
  flavor            = var.rds_flavor_id != "" ? var.rds_flavor_id : try(data.huaweicloud_rds_flavors.test[0].flavors[0].name, null)
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = var.availability_zone != "" ? [var.availability_zone] : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1), null)

  db {
    type     = "MySQL"
    version  = var.rds_db_version
    port     = var.rds_db_port
    password = var.rds_db_password
  }

  volume {
    type = var.rds_volume_type
    size = var.rds_volume_size
  }

  lifecycle {
    ignore_changes = [flavor]
  }
}

data "huaweicloud_dws_flavors" "test" {
  count = var.dws_node_type == "" || var.dws_version == "" ? 1 : 0

  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  vcpus             = var.dws_flavor_vcpus
  memory            = var.dws_flavor_memory
  datastore_type    = var.dws_datastore_type
}

resource "huaweicloud_dws_cluster" "test" {
  name                  = var.dws_cluster_name
  node_type             = var.dws_node_type != "" ? var.dws_node_type : try(data.huaweicloud_dws_flavors.test[0].flavors[0].flavor_id, null)
  number_of_node        = var.dws_number_of_node
  number_of_cn          = var.dws_number_of_cn
  version               = var.dws_version != "" ? var.dws_version : try(data.huaweicloud_dws_flavors.test[0].flavors[0].datastore_version, null)
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  user_name             = var.dws_admin_user_name
  user_pwd              = var.dws_admin_user_pwd
  port                  = var.dws_port
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  volume {
    type     = var.dws_volume_type
    capacity = var.dws_volume_capacity
  }
}

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = var.elastic_resource_pool_name
  description           = var.elastic_resource_pool_description
  min_cu                = var.elastic_resource_pool_min_cu
  max_cu                = var.elastic_resource_pool_max_cu
  cidr                  = var.elastic_resource_pool_cidr
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  label                 = var.elastic_resource_pool_label
}

resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1
  name                       = var.queue_name
  queue_type                 = "general"
  cu_count                   = var.queue_cu_count
  description                = var.queue_description
  enterprise_project_id      = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_dli_datasource_connection" "test" {
  name      = var.datasource_connection_name
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = huaweicloud_dli_datasource_connection.test.id
  elastic_resource_pools = [huaweicloud_dli_elastic_resource_pool.test.name]

  depends_on = [huaweicloud_dli_queue.test]
}
