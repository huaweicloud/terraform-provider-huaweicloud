data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
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

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip  = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
  ipv6_enable = true
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name                  = var.loadbalancer_name
  vpc_id                = huaweicloud_vpc.test.id
  ipv4_subnet_id        = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id       = huaweicloud_vpc_subnet.test.id
  availability_zone     = var.availability_zone != "" ? [var.availability_zone] : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1), null)
  cross_vpc_backend     = var.loadbalancer_cross_vpc_backend
  description           = var.loadbalancer_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.loadbalancer_tags
}

resource "huaweicloud_elb_listener" "test" {
  loadbalancer_id    = huaweicloud_elb_loadbalancer.test.id
  name               = var.listener_name
  protocol           = var.listener_protocol
  protocol_port      = var.listener_port
  server_certificate = var.listener_server_certificate
  ca_certificate     = var.listener_ca_certificate
  sni_certificate    = var.listener_sni_certificates
  sni_match_algo     = var.listener_sni_match_algo
  security_policy_id = var.listener_security_policy_id
  http2_enable       = var.listener_http2_enable

  dynamic "port_ranges" {
    for_each = var.listener_port_ranges

    content {
      start_port = port_ranges.value["start_port"]
      end_port   = port_ranges.value["end_port"]
    }
  }

  idle_timeout                = var.listener_idle_timeout
  request_timeout             = var.listener_request_timeout
  response_timeout            = var.listener_response_timeout
  description                 = var.listener_description
  tags                        = var.listener_tags
  advanced_forwarding_enabled = var.listener_advanced_forwarding_enabled
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

resource "huaweicloud_elb_monitor" "test" {
  pool_id     = huaweicloud_elb_pool.test.id
  protocol    = var.health_check_protocol
  interval    = var.health_check_interval
  timeout     = var.health_check_timeout
  max_retries = var.health_check_max_retries
  port        = var.health_check_port
  url_path    = var.health_check_url_path
  status_code = var.health_check_status_code
  http_method = var.health_check_http_method
  domain_name = var.health_check_domain_name
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = var.pool_protocol == "UDP" ? "udp" : "tcp"
  # The backend server port and check health port.
  ports             = var.health_check_port != null ? join(",", distinct([var.health_check_port, var.member_protocol_port])) : var.member_protocol_port
  # The CIDR to which the ELB backend subnet belongs.
  remote_ip_prefix  = huaweicloud_vpc_subnet.test.cidr
}

# ST.001 Disable
resource "huaweicloud_networking_secgroup_rule" "in_v4_icmp" {
  # ST.001 Enable
  count = var.pool_protocol == "UDP" ? 1 : 0

  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = huaweicloud_vpc_subnet.test.cidr
}

# ST.001 Disable
resource "huaweicloud_networking_secgroup_rule" "in_v6_icmp" {
  # ST.001 Enable
  count = var.pool_protocol == "UDP" ? 1 : 0

  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv6"
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = huaweicloud_vpc_subnet.test.ipv6_cidr
}

# ST.001 Disable
resource "huaweicloud_networking_secgroup_rule" "in_v6" {
  # ST.001 Enable
  count = var.pool_protocol == "UDP" ? 1 : 0

  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv6"
  direction         = "ingress"
  protocol          = "udp"
  # The backend server port and check health port.
  ports             = var.health_check_port != null ? join(",", distinct([var.health_check_port, var.member_protocol_port])) : var.member_protocol_port
  # The CIDR to which the ELB backend subnet belongs.
  remote_ip_prefix  = huaweicloud_vpc_subnet.test.ipv6_cidr
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

  lifecycle {
    ignore_changes = [
      flavor_id,
      image_id,
      availability_zone
    ]
  }
}

resource "huaweicloud_elb_member" "test" {
  pool_id       = huaweicloud_elb_pool.test.id
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = var.member_protocol_port
  weight        = var.member_weight
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
