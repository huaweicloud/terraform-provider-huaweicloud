data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
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
  name                 = var.security_group_name
  delete_default_rules = true
}

data "huaweicloud_dbss_flavors" "test" {
  count = var.instance_flavor == "" ? 1 : 0
}

locals {
  vpc_name      = var.vpc_name
  subnet_name   = var.subnet_name
  instance_name = var.instance_name

  product_spec_desc = jsonencode(
  {
    "specDesc" : {
      "zh-cn" : {
        "主机名称" : local.instance_name,
        "虚拟私有云" : local.vpc_name,
        "子网" : local.subnet_name
      },
      "en-us" : {
        "Instance Name" : local.instance_name,
        "VPC" : local.vpc_name,
        "Subnet" : local.subnet_name
      }
    }
  }
  )
}

resource "huaweicloud_dbss_instance" "test" {
  name                  = var.instance_name
  availability_zone     = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  flavor                = var.instance_flavor == "" ? try(data.huaweicloud_dbss_flavors.test[0].flavors[0].id, null) : var.instance_flavor
  product_spec_desc     = local.product_spec_desc
  resource_spec_code    = var.instance_spec_code
  description           = var.instance_description
  tags                  = var.instance_tags
  enterprise_project_id = var.enterprise_project_id
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew
}

resource "huaweicloud_dbss_ecs_database" "test" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
  name        = var.database_name
  type        = var.database_type
  version     = var.database_version
  ip          = var.database_ip_address
  port        = var.database_port
  os          = var.database_os
  charset     = var.database_charset
}
