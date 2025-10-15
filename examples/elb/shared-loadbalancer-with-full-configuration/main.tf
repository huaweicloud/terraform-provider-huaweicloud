data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
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
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
  dns_list   = var.subnet_dns_list
}

resource "huaweicloud_lb_loadbalancer" "test" {
  name          = var.loadbalancer_name
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_vpc_eip" "test" {
  count = var.is_associate_eip && var.eip_address == "" ? 1 : 0

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_vpc_eipv3_associate" "test" {
  count = var.is_associate_eip ? 1 : 0

  publicip_id             = var.eip_address != "" ? var.eip_address : huaweicloud_vpc_eip.test[0].id
  associate_instance_type = "ELB"
  associate_instance_id   = huaweicloud_lb_loadbalancer.test.id
}

resource "huaweicloud_lb_certificate" "test" {
  count = var.listener_protocol == "TERMINATED_HTTPS" && var.listener_default_tls_container_ref == "" ? 1 : 0

  name        = var.listener_server_certificate_name
  private_key = var.listener_server_certificate_private_key
  certificate = var.listener_server_certificate_certificate
}

resource "huaweicloud_lb_listener" "test" {
  loadbalancer_id             = huaweicloud_lb_loadbalancer.test.id
  name                        = var.listener_name
  protocol                    = var.listener_protocol
  protocol_port               = var.listener_port
  description                 = var.listener_description
  tags                        = var.listener_tags
  http2_enable                = var.listener_http2_enable
  default_tls_container_ref   = var.listener_default_tls_container_ref != "" ? var.listener_default_tls_container_ref : try(huaweicloud_lb_certificate.test[0].id, null)
  client_ca_tls_container_ref = var.listener_client_ca_tls_container_ref
  sni_container_refs          = var.listener_sni_container_refs
  tls_ciphers_policy          = var.listener_tls_ciphers_policy

  dynamic "insert_headers" {
    for_each = length(var.listener_insert_headers) > 0 ? [var.listener_insert_headers] : []

    content {
      x_forwarded_elb_ip = insert_headers.value["x_forwarded_elb_ip"]
      x_forwarded_host   = insert_headers.value["x_forwarded_host"]
    }
  }
}

resource "huaweicloud_lb_pool" "test" {
  listener_id = huaweicloud_lb_listener.test.id
  name        = var.pool_name
  protocol    = var.pool_protocol
  lb_method   = var.pool_method
  description = var.pool_description

  dynamic "persistence" {
    for_each = var.pool_persistence != null ? [var.pool_persistence] : []

    content {
      type        = persistence.value["type"]
      cookie_name = persistence.value["cookie_name"]
      timeout     = persistence.value["timeout"]
    }
  }
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

# The subnet ID must belong to the VPC of the load balancer.
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

resource "huaweicloud_lb_member" "test" {
  pool_id       = huaweicloud_lb_pool.test.id
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = var.member_protocol_port
  weight        = var.member_weight
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = var.health_check_port != null ? join(",", distinct([var.health_check_port, var.member_protocol_port])) : var.member_protocol_port
  remote_ip_prefix  = var.security_group_rule_remote_ip_prefix
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_lb_monitor" "test" {
  pool_id        = huaweicloud_lb_pool.test.id
  name           = var.health_check_name
  type           = var.health_check_type
  delay          = var.health_check_delay
  timeout        = var.health_check_timeout
  max_retries    = var.health_check_max_retries
  port           = var.health_check_port
  url_path       = var.health_check_url_path
  http_method    = var.health_check_http_method
  expected_codes = var.health_check_expected_codes
  domain_name    = var.health_check_domain_name

  depends_on = [huaweicloud_networking_secgroup_rule.test]
}
