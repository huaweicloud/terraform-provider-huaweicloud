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

# EIP in the backend region
resource "huaweicloud_vpc_eip" "test" {
  region = var.backend_region

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.eip_name
    size        = var.bandwidth_size
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

# GA endpoint
resource "huaweicloud_ga_endpoint" "test" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  resource_id       = huaweicloud_vpc_eip.test.id
  ip_address        = huaweicloud_vpc_eip.test.address
  resource_type     = "EIP"
  weight            = var.endpoint_weight
}
