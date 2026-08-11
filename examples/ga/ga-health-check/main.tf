# GA accelerator with IPV4 and IPV6
resource "huaweicloud_ga_accelerator" "test" {
  name        = var.accelerator_name
  description = var.accelerator_description

  ip_sets {
    ip_type = "IPV4"
    area    = var.ip_area
  }

  ip_sets {
    ip_type = "IPV6"
    area    = var.ip_area
  }

  tags = var.tags
}

# GA listener
resource "huaweicloud_ga_listener" "test" {
  accelerator_id = huaweicloud_ga_accelerator.test.id
  name           = var.listener_name
  protocol       = var.listener_protocol
  description    = var.listener_description

  port_ranges {
    from_port = var.port_from
    to_port   = var.port_to
  }

  tags = var.tags
}

# GA endpoint group
resource "huaweicloud_ga_endpoint_group" "test" {
  name        = var.endpoint_group_name
  description = var.endpoint_group_description
  region_id   = var.backend_region

  listeners {
    id = huaweicloud_ga_listener.test.id
  }
}

# GA health check
resource "huaweicloud_ga_health_check" "test" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  enabled           = var.health_check_enabled
  protocol          = "TCP"
  interval          = var.health_check_interval
  max_retries       = var.health_check_max_retries
  port              = var.health_check_port
  timeout           = var.health_check_timeout
}
