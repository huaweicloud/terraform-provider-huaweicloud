resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
}

data "huaweicloud_compute_flavors" "test" {
  count = var.dedicated_instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  performance_type  = var.dedicated_instance_performance_type
  cpu_core_count    = var.dedicated_instance_cpu_core_count
  memory_size       = var.dedicated_instance_memory_size
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_waf_dedicated_instance" "test" {
  name               = var.dedicated_instance_name
  specification_code = var.dedicated_instance_specification_code
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  available_zone     = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  ecs_flavor         = var.dedicated_instance_flavor_id == "" ? data.huaweicloud_compute_flavors.test[0].ids[0] : var.dedicated_instance_flavor_id

  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

resource "huaweicloud_waf_policy" "test" {
  name  = var.policy_name
  level = var.policy_level

  depends_on = [
    huaweicloud_waf_dedicated_instance.test
  ]
}

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain      = var.dedicated_mode_domain_name
  policy_id   = huaweicloud_waf_policy.test.id
  keep_policy = true

  server {
    client_protocol = var.dedicated_domain_client_protocol
    server_protocol = var.dedicated_domain_server_protocol
    address         = var.dedicated_domain_address == "" ? cidrhost(huaweicloud_vpc_subnet.test.cidr, 4) : var.dedicated_domain_address
    port            = var.dedicated_domain_port
    type            = var.dedicated_domain_type
    vpc_id          = huaweicloud_vpc.test.id
  }
}
