data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].ids[0], null)
  visibility = "public"
  os         = "Ubuntu"
}

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

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_kps_keypair" "test" {
  name       = var.keypair_name
  public_key = var.keypair_public_key != "" ? var.keypair_public_key : null
}

resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = var.configuration_name

  instance_config {
    image    = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
    flavor   = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
    key_name = huaweicloud_kps_keypair.test.id

    dynamic "disk" {
      for_each = var.disk_configurations

      content {
        disk_type   = disk.value.disk_type
        volume_type = disk.value.volume_type
        size        = disk.value.volume_size
      }
    }
  }
}

resource "huaweicloud_as_group" "test" {
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  vpc_id                   = huaweicloud_vpc.test.id
  scaling_group_name       = var.group_name
  desire_instance_number   = var.desire_instance_number
  min_instance_number      = var.min_instance_number
  max_instance_number      = var.max_instance_number
  delete_publicip          = var.is_delete_publicip
  delete_instances         = var.is_delete_instances ? "yes" : "no"

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }

  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}

resource "huaweicloud_smn_topic" "test" {
  name = var.topic_name
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name = var.alarm_rule_name

  metric {
    namespace = "SYS.AS"
  }

  resources {
    dimensions {
      name  = "AutoScalingGroup"
      value = huaweicloud_as_group.test.id
    }
  }

  dynamic "condition" {
    for_each = var.rule_conditions

    content {
      alarm_level         = condition.value.alarm_level
      metric_name         = condition.value.metric_name
      period              = condition.value.period
      filter              = condition.value.filter
      comparison_operator = condition.value.comparison_operator
      suppress_duration   = condition.value.suppress_duration
      value               = condition.value.value
      count               = condition.value.count
    }
  }

  alarm_actions {
    type              = "autoscaling"
    notification_list = [huaweicloud_smn_topic.test.id]
  }
}

# ST.001 Disable
resource "huaweicloud_as_policy" "scaling_up" {
  scaling_policy_name = var.scaling_up_policy_name
  scaling_policy_type = "ALARM"
  scaling_group_id    = huaweicloud_as_group.test.id
  alarm_id            = huaweicloud_ces_alarmrule.test.id
  cool_down_time      = var.scaling_up_cool_down_time

  scaling_policy_action {
    operation       = "ADD"
    instance_number = var.scaling_up_instance_number
  }
}

resource "huaweicloud_as_policy" "scaling_down" {
  scaling_policy_name = var.scaling_down_policy_name
  scaling_policy_type = "ALARM"
  scaling_group_id    = huaweicloud_as_group.test.id
  alarm_id            = huaweicloud_ces_alarmrule.test.id
  cool_down_time      = var.scaling_down_cool_down_time

  scaling_policy_action {
    operation       = "REMOVE"
    instance_number = var.scaling_down_instance_number
  }
}
# ST.001 Enable
