data "huaweicloud_availability_zones" "test" {}

# Create a VPC
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

# Create a subnet
resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
}

# Create a security group
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.secgroup_name
  delete_default_rules = true
}

# Create a SFS Turbo file system
resource "huaweicloud_sfs_turbo" "test" {
  name              = var.turbo_name
  size              = var.turbo_size
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
}

# Optional: Create a backup policy
resource "huaweicloud_cbr_policy" "test" {
  count = var.enable_policy ? 1 : 0

  name        = "${var.vault_name}-policy"
  type        = "backup"
  time_period = 20
  time_zone   = "UTC+08:00"
  enabled     = true

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00"]
  }
}

# Create a CBR vault for SFS Turbo backup
resource "huaweicloud_cbr_vault" "test" {
  name                  = var.vault_name
  type                  = "turbo"
  protection_type       = var.protection_type
  size                  = var.vault_size
  auto_expand           = var.auto_expand
  enterprise_project_id = var.enterprise_project_id
  backup_name_prefix    = var.backup_name_prefix
  is_multi_az           = var.is_multi_az

  resources {
    includes = [
      huaweicloud_sfs_turbo.test.id
    ]
  }

  dynamic "policy" {
    for_each = var.enable_policy ? [1] : []
    content {
      id = huaweicloud_cbr_policy.test[0].id
    }
  }

  tags = var.tags
}
