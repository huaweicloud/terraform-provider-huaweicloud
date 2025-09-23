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

resource "huaweicloud_dms_rocketmq_migration_task" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = var.migration_task_name
  overwrite   = var.migration_task_overwrite
  type        = var.migration_task_type

  dynamic "topic_configs" {
    for_each = var.migration_task_topic_configs

    content {
      order             = topic_configs.value["order"]
      perm              = topic_configs.value["perm"]
      read_queue_num    = topic_configs.value["read_queue_nums"]
      topic_filter_type = topic_configs.value["topic_filter_type"]
      topic_name        = topic_configs.value["topic_name"]
      topic_sys_flag    = topic_configs.value["topic_sys_flag"]
      write_queue_num   = topic_configs.value["write_queue_nums"]
    }
  }

  dynamic "subscription_groups" {
    for_each = var.migration_task_subscription_groups

    content {
      consume_broadcast_enable          = subscription_groups.value["consume_broadcast_enable"]
      consume_enable                    = subscription_groups.value["consume_enable"]
      consume_from_min_enable           = subscription_groups.value["consume_from_min_enable"]
      group_name                        = subscription_groups.value["group_name"]
      notify_consumerids_changed_enable = subscription_groups.value["notify_consumerids_changed_enable"]
      retry_max_times                   = subscription_groups.value["retry_max_times"]
      retry_queue_num                   = subscription_groups.value["retry_queue_num"]
      which_broker_when_consume_slow    = subscription_groups.value["which_broker_when_consume_slow"]
    }
  }

  dynamic "vhosts" {
    for_each = var.migration_task_vhosts

    content {
      name = vhosts.value["name"]
    }
  }

  dynamic "queues" {
    for_each = var.migration_task_queues

    content {
      name    = queues.value["name"]
      vhost   = queues.value["vhost"]
      durable = queues.value["durable"]
    }
  }

  dynamic "exchanges" {
    for_each = var.migration_task_exchanges

    content {
      name    = exchanges.value["name"]
      vhost   = exchanges.value["vhost"]
      type    = exchanges.value["type"]
      durable = exchanges.value["durable"]
    }
  }

  dynamic "bindings" {
    for_each = var.migration_task_bindings

    content {
      vhost            = bindings.value["vhost"]
      source           = bindings.value["source"]
      destination      = bindings.value["destination"]
      destination_type = bindings.value["destination_type"]
      routing_key      = bindings.value["routing_key"]
    }
  }
}
