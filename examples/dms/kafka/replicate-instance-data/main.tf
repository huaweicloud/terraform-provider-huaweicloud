data "huaweicloud_availability_zones" "test" {
  count = anytrue([for v in var.instance_configurations : length(v.availability_zones) == 0]) ? 1 : 0
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

locals {
  instance_configurations_without_flavor_id = [for v in var.instance_configurations : v if v.flavor_id == ""]
}

data "huaweicloud_dms_kafka_flavors" "test" {
  count = length(local.instance_configurations_without_flavor_id)

  type               = local.instance_configurations_without_flavor_id[count.index].flavor_type
  availability_zones = length(local.instance_configurations_without_flavor_id[count.index].availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1)) : null
  storage_spec_code  = local.instance_configurations_without_flavor_id[count.index].storage_spec_code
}

resource "huaweicloud_dms_kafka_instance" "test" {
  count = length(var.instance_configurations)

  name               = var.instance_configurations[count.index].name
  availability_zones = length(var.instance_configurations[count.index].availability_zones) > 0 ? var.instance_configurations[count.index].availability_zones : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1))
  engine_version     = var.instance_configurations[count.index].engine_version
  flavor_id          = var.instance_configurations[count.index].flavor_id != "" ? var.instance_configurations[count.index].flavor_id : try(data.huaweicloud_dms_kafka_flavors.test[count.index].flavors[0].id, null)
  storage_spec_code  = var.instance_configurations[count.index].storage_spec_code
  storage_space      = var.instance_configurations[count.index].storage_space
  broker_num         = var.instance_configurations[count.index].broker_num
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  access_user        = var.instance_configurations[count.index].access_user
  password           = var.instance_configurations[count.index].password
  enabled_mechanisms = var.instance_configurations[count.index].enabled_mechanisms

  dynamic "port_protocol" {
    for_each = length(var.instance_configurations[count.index].port_protocol) > 0 ? [var.instance_configurations[count.index].port_protocol] : []

    content {
      private_plain_enable          = port_protocol.value.private_plain_enable
      private_sasl_ssl_enable       = port_protocol.value.private_sasl_ssl_enable
      private_sasl_plaintext_enable = port_protocol.value.private_sasl_plaintext_enable
    }
  }

  lifecycle {
    ignore_changes = [
      availability_zones,
      flavor_id,
    ]
  }
}

resource "huaweicloud_dms_kafka_topic" "test" {
  count = length(var.task_topics) == 0 ? 1 : 0

  instance_id      = huaweicloud_dms_kafka_instance.test[0].id
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

resource "huaweicloud_dms_kafka_smart_connect" "test" {
  instance_id       = huaweicloud_dms_kafka_instance.test[0].id
  storage_spec_code = var.smart_connect_storage_spec_code
  bandwidth         = var.smart_connect_bandwidth
  node_count        = var.smart_connect_node_count
}

resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test[0].id
  task_name   = var.task_name
  source_type = "KAFKA_REPLICATOR_SOURCE"
  start_later = var.task_start_later
  topics      = length(var.task_topics) > 0 ? var.task_topics : huaweicloud_dms_kafka_topic.test[*].name

  source_task {
    peer_instance_id              = huaweicloud_dms_kafka_instance.test[1].id
    direction                     = var.task_direction
    replication_factor            = var.task_replication_factor
    task_num                      = var.task_task_num
    provenance_header_enabled     = var.task_provenance_header_enabled
    sync_consumer_offsets_enabled = var.task_sync_consumer_offsets_enabled
    rename_topic_enabled          = var.task_rename_topic_enabled
    consumer_strategy             = var.task_consumer_strategy
    compression_type              = var.task_compression_type
    topics_mapping                = var.task_topics_mapping
    security_protocol             = try(huaweicloud_dms_kafka_instance.test[1].port_protocol[0].private_sasl_ssl_enable, false) ? "SASL_SSL" : try(huaweicloud_dms_kafka_instance.test[1].port_protocol[0].private_sasl_plaintext_enable, false) ? "PLAINTEXT" : null
    sasl_mechanism                = try(tolist(huaweicloud_dms_kafka_instance.test[1].enabled_mechanisms)[0], null)
    user_name                     = try(huaweicloud_dms_kafka_instance.test[1].access_user, null)
    password                      = try(huaweicloud_dms_kafka_instance.test[1].password, null)
  }

  depends_on = [huaweicloud_dms_kafka_smart_connect.test]
}
