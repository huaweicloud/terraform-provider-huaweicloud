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

data "huaweicloud_dms_rocketmq_availability_zones" "test" {
  count = length(var.availability_zones) == 0 ? 1 : 0
}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  type               = var.instance_flavor_type
  availability_zones = length(var.availability_zones) != 0 ? var.availability_zones : slice(data.huaweicloud_dms_rocketmq_availability_zones.test[0].availability_zones[*].code, 0, var.availability_zones_count)
  charging_mode      = var.charging_mode
}

resource "huaweicloud_vpc_eip" "test" {
  count = var.instance_enable_publicip ? var.instance_publicip_id != "" ? 0 : var.instance_eips_count : 0

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = format("%s-%d", var.bandwidth_name, count.index)
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name                  = var.instance_name
  flavor_id             = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_dms_rocketmq_flavors.test[0].flavors[0].id, null)
  engine_version        = var.instance_engine_version != "" ? var.instance_engine_version : try(data.huaweicloud_dms_rocketmq_flavors.test[0].versions[0], null)
  storage_spec_code     = var.instance_storage_spec_code != "" ? var.instance_storage_spec_code : try(data.huaweicloud_dms_rocketmq_flavors.test[0].flavors[0].ios[0].storage_spec_code, null)
  storage_space         = var.instance_storage_space
  availability_zones    = length(var.availability_zones) != 0 ? var.availability_zones : slice(data.huaweicloud_dms_rocketmq_availability_zones.test[0].availability_zones[*].code, 0, var.availability_zones_count)
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  broker_num            = var.instance_broker_num
  description           = var.instance_description
  tags                  = var.instance_tags
  enterprise_project_id = var.enterprise_project_id
  enable_acl            = var.instance_enable_acl
  tls_mode              = var.instance_tls_mode
  enable_publicip       = var.instance_enable_publicip
  publicip_id           = var.instance_enable_publicip ? var.instance_publicip_id != "" ? var.instance_publicip_id : join(",", huaweicloud_vpc_eip.test[*].id) : null

  dynamic "configs" {
    for_each = var.instance_configs

    content {
      name  = configs.value.name
      value = configs.value.value
    }
  }

  charging_mode = var.charging_mode
  period_unit   = var.period_unit
  period        = var.period
  auto_renew    = var.auto_renew
}
