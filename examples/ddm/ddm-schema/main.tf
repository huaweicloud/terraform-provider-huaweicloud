data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) == 0 ? 1 : 0
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "random_password" "test" {
  count = var.rds_instance_password == "" ? 1 : 0

  length           = 12
  special          = true
  override_special = "~!@#%^*-_+?"
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
}

data "huaweicloud_rds_flavors" "test" {
  count = var.instance_flavor == "" ? 1 : 0

  db_type       = var.database_type
  db_version    = var.database_version
  instance_mode = var.instance_mode
  group_type    = var.instance_group_type
  vcpus         = var.instance_flavor_vcpus
}

resource "huaweicloud_rds_instance" "test" {
  name              = var.rds_instance_name
  availability_zone = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1)) : var.availability_zones
  flavor            = var.instance_flavor == "" ? try(data.huaweicloud_rds_flavors.test[0].flavors[0].name, null) : var.instance_flavor
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  db {
    type     = var.database_type
    version  = var.database_version
    port     = var.database_port
    password = var.rds_instance_password == "" ? try(random_password.test[0].result, null) : var.rds_instance_password
  }

  volume {
    type = var.volume_type
    size = var.volume_size
  }
}

data "huaweicloud_ddm_engines" "test" {
  count = var.instance_engine_id == "" ? 1 : 0
}

data "huaweicloud_ddm_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  engine_id = var.instance_engine_id == "" ? try(data.huaweicloud_ddm_engines.test[0].engines[0].id, null)  : var.instance_engine_id
}

resource "huaweicloud_ddm_instance" "test" {
  name               = var.ddm_instance_name
  availability_zones = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1)) : var.availability_zones
  engine_id          = var.instance_engine_id == "" ? try(data.huaweicloud_ddm_engines.test[0].engines[0].id, null)  : var.instance_engine_id
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_ddm_flavors.test[0].flavors[0].id, null)  : var.instance_flavor_id
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  node_num           = var.instance_node_num

  dynamic "parameters" {
    for_each = var.instance_parameters

    content {
      name  = parameters.value.name
      value = parameters.value.value
    }
  }
}

resource "huaweicloud_ddm_schema" "test" {
  instance_id  = huaweicloud_ddm_instance.test.id
  name         = var.schema_name
  shard_mode   = var.schema_shard_mode
  shard_number = var.schema_shard_number

  data_nodes {
    id             = huaweicloud_rds_instance.test.id
    admin_user     = "root"
    admin_password = var.rds_instance_password == "" ? try(random_password.test[0].result, null) : var.rds_instance_password
  }

  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
