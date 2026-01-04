# Create VPC and subnet
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

# Create security group
resource "huaweicloud_networking_secgroup" "test" {
  name        = var.security_group_name
  description = "Security group for DCS data migration"
}

# Get availability zones
data "huaweicloud_availability_zones" "test" {}

# Get DCS flavors
data "huaweicloud_dcs_flavors" "test" {
  cache_mode     = var.instance_cache_mode
  capacity       = var.instance_capacity
  engine_version = var.instance_engine_version
}

# Create source DCS Redis instance
resource "huaweicloud_dcs_instance" "test" {
  count = 2

  name               = "${var.instance_name}-${count.index}"
  engine             = "Redis"
  engine_version     = var.instance_engine_version
  capacity           = var.instance_capacity
  flavor             = try(data.huaweicloud_dcs_flavors.test.flavors[0].name, null)
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 2), null)
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  password           = var.instance_password

  lifecycle {
    ignore_changes = [
      security_group_id
    ]
  }
}

# ST.001 Disable
# Create full migration task
resource "huaweicloud_dcs_online_data_migration_task" "full_migration" {
  task_name          = var.full_migration_task_name
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  description        = var.full_migration_task_description
  migration_method   = "full_amount_migration"
  resume_mode        = var.full_migration_resume_mode
  bandwidth_limit_mb = var.full_migration_bandwidth_limit_mb != "" ? var.full_migration_bandwidth_limit_mb : null

  source_instance {
    id       = huaweicloud_dcs_instance.test[0].id
    password = var.instance_password
  }

  target_instance {
    id       = huaweicloud_dcs_instance.test[1].id
    password = var.instance_password
  }

  lifecycle {
    ignore_changes = [
      source_instance,
      target_instance
    ]
  }

  depends_on = [
    huaweicloud_dcs_instance.test
  ]
}

# Create incremental migration task
resource "huaweicloud_dcs_online_data_migration_task" "incremental_migration" {
  count = var.enable_incremental_migration ? 1 : 0

  task_name          = var.incremental_migration_task_name
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  description        = var.incremental_migration_task_description
  migration_method   = "incremental_migration"
  resume_mode        = var.incremental_migration_resume_mode
  bandwidth_limit_mb = var.incremental_migration_bandwidth_limit_mb != "" ? var.incremental_migration_bandwidth_limit_mb : null

  source_instance {
    id       = huaweicloud_dcs_instance.test[0].id
    password = var.instance_password
  }

  target_instance {
    id       = huaweicloud_dcs_instance.test[1].id
    password = var.instance_password
  }

  lifecycle {
    ignore_changes = [
      source_instance,
      target_instance
    ]
  }

  depends_on = [
    huaweicloud_dcs_instance.test,
    huaweicloud_dcs_online_data_migration_task.full_migration
  ]
}
# ST.001 Enable
