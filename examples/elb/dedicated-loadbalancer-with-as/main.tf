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
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

data "huaweicloud_elb_flavors" "test" {
  type = "L7"
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name                  = var.loadbalancer_name
  vpc_id                = huaweicloud_vpc.test.id
  ipv4_subnet_id        = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  l7_flavor_id          = try(data.huaweicloud_elb_flavors.test.ids[0], null)
  availability_zone     = var.availability_zone != "" ? [var.availability_zone] : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1), null)
  cross_vpc_backend     = var.loadbalancer_cross_vpc_backend
  description           = var.loadbalancer_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.loadbalancer_tags
  force_delete          = var.loadbalancer_force_delete
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_vpc_eipv3_associate" "test" {
  publicip_id             = huaweicloud_vpc_eip.test.id
  associate_instance_type = "ELB"
  associate_instance_id   = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_listener" "test" {
  loadbalancer_id  = huaweicloud_elb_loadbalancer.test.id
  name             = var.listener_name
  protocol         = var.listener_protocol
  protocol_port    = var.listener_port
  idle_timeout     = var.listener_idle_timeout
  request_timeout  = var.listener_request_timeout
  response_timeout = var.listener_response_timeout
  description      = var.listener_description
  tags             = var.listener_tags
}

resource "huaweicloud_elb_pool" "test" {
  listener_id     = huaweicloud_elb_listener.test.id
  name            = var.pool_name
  protocol        = var.pool_protocol
  lb_method       = var.pool_method
  any_port_enable = var.pool_any_port_enable
  description     = var.pool_description

  dynamic "persistence" {
    for_each = var.pool_persistences

    content {
      type        = persistence.value["type"]
      cookie_name = persistence.value["cookie_name"]
      timeout     = persistence.value["timeout"]
    }
  }
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].ids[0], null)
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = var.configuration_name

  instance_config {
    image              = var.configuration_image_id != "" ? var.configuration_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
    flavor             = var.configuration_flavor_id != "" ? var.configuration_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
    security_group_ids = [huaweicloud_networking_secgroup.test.id]
    user_data          = var.configuration_user_data

    dynamic "disk" {
      for_each = var.configuration_disks

      content {
        size        = disk.value["size"]
        volume_type = disk.value["volume_type"]
        disk_type   = disk.value["disk_type"]
      }
    }
  }
}

resource "huaweicloud_as_group" "test" {
  scaling_group_name       = var.group_name
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  desire_instance_number   = var.group_desire_instance_number
  min_instance_number      = var.group_min_instance_number
  max_instance_number      = var.group_max_instance_number
  vpc_id                   = huaweicloud_vpc.test.id
  delete_publicip          = var.group_delete_publicip
  delete_instances         = var.group_delete_instances
  force_delete             = var.group_force_delete

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }

  lbaas_listeners {
    pool_id       = huaweicloud_elb_pool.test.id
    protocol_port = huaweicloud_elb_listener.test.protocol_port
  }

  lifecycle {
    ignore_changes = [
      # When instances are auto-scaled, the desire instance number will be changed.
      desire_instance_number,
    ]
  }
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_compute_instance" "test" {
  name              = var.instance_name
  image_id          = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  flavor_id         = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  security_groups   = [huaweicloud_networking_secgroup.test.name]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  # Prevents resource changes when data source query results change.
  lifecycle {
    ignore_changes = [
      flavor_id,
      image_id,
      availability_zone,
    ]
  }
}

resource "huaweicloud_as_instance_attach" "test" {
  scaling_group_id = huaweicloud_as_group.test.id
  instance_id      = huaweicloud_compute_instance.test.id
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name = var.alarm_rule_name

  metric {
    namespace = "SYS.AS"

    dimensions {
      name  = "AutoScalingGroup"
      value = huaweicloud_as_group.test.id
    }
  }

  dynamic "condition" {
    for_each = var.alarm_rule_conditions

    content {
      period              = condition.value["period"]
      filter              = condition.value["filter"]
      comparison_operator = condition.value["comparison_operator"]
      value               = condition.value["value"]
      unit                = condition.value["unit"]
      count               = condition.value["count"]
      alarm_level         = condition.value["alarm_level"]
      metric_name         = condition.value["metric_name"]
    }
  }

  alarm_actions {
    type              = "autoscaling"
    notification_list = []
  }
}

resource "huaweicloud_as_policy" "test" {
  scaling_policy_name = var.policy_name
  scaling_policy_type = "ALARM"
  scaling_group_id    = huaweicloud_as_group.test.id
  alarm_id            = huaweicloud_ces_alarmrule.test.id
  cool_down_time      = var.policy_cool_down_time

  scaling_policy_action {
    operation       = var.policy_operation
    instance_number = var.policy_instance_number
  }
}
