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

  dynamic "configs" {
    for_each = var.instance_configs

    content {
      name  = configs.value.name
      value = configs.value.value
    }
  }
}

data "huaweicloud_dms_rocketmq_broker" "test" {
  count = length(var.topic_brokers) == 0 && huaweicloud_dms_rocketmq_instance.test.engine_version == "4.8.0" ? 1 : 0

  instance_id = huaweicloud_dms_rocketmq_instance.test.id
}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id  = huaweicloud_dms_rocketmq_instance.test.id
  name         = var.topic_name
  message_type = var.topic_message_type
  queue_num    = var.topic_queue_num
  permission   = var.topic_permission

  dynamic "brokers" {
    for_each = length(var.topic_brokers) > 0 ? var.topic_brokers : try(data.huaweicloud_dms_rocketmq_broker.test[0].brokers, [])

    content {
      name = brokers.value
    }
  }
}

resource "huaweicloud_dms_rocketmq_message_send" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  topic       = huaweicloud_dms_rocketmq_topic.test.name
  body        = var.message_body

  dynamic "property_list" {
    for_each = var.message_properties

    content {
      name  = property_list.value.name
      value = property_list.value.value
    }
  }
}
