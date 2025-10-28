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

data "huaweicloud_dms_kafka_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  type               = var.instance_flavor_type
  availability_zones = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1)) : var.availability_zones
  storage_spec_code  = var.instance_storage_spec_code
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = var.instance_name
  availability_zones = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1)) : var.availability_zones
  engine_version     = var.instance_engine_version
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_dms_kafka_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  storage_spec_code  = var.instance_storage_spec_code
  storage_space      = var.instance_storage_space
  broker_num         = var.instance_broker_num
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  access_user        = var.instance_access_user_name
  password           = var.instance_access_user_password
  enabled_mechanisms = var.instance_enabled_mechanisms

  dynamic "port_protocol" {
    for_each = [var.port_protocol]

    content {
      private_plain_enable          = port_protocol.value.private_plain_enable
      private_sasl_ssl_enable       = port_protocol.value.private_sasl_ssl_enable
      private_sasl_plaintext_enable = port_protocol.value.private_sasl_plaintext_enable
      public_plain_enable           = port_protocol.value.public_plain_enable
      public_sasl_ssl_enable        = port_protocol.value.public_sasl_ssl_enable
      public_sasl_plaintext_enable  = port_protocol.value.public_sasl_plaintext_enable
    }
  }

  # If you want to change some of the following parameters, you need to remove the corresponding fields from "lifecycle.ignore_changes".
  lifecycle {
    ignore_changes = [
      availability_zones,
      flavor_id,
    ]
  }
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id      = huaweicloud_dms_kafka_instance.test.id
  name             = var.topic_name
  partitions       = var.topic_partitions
  replicas         = var.topic_replicas
  aging_time       = var.topic_aging_time
  sync_replication = var.topic_sync_replication
  sync_flushing    = var.topic_sync_flushing
  description      = var.topic_description

  dynamic "configs" {
    for_each = var.topic_configs

    content {
      name  = configs.value.name
      value = configs.value.value
    }
  }
}

resource "huaweicloud_dms_kafka_message_produce" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  topic       = huaweicloud_dms_kafka_topic.test.name
  body        = var.message_body

  dynamic "property_list" {
    for_each = var.message_properties

    content {
      name  = property_list.value.name
      value = property_list.value.value
    }
  }
}
