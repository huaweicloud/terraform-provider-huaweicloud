# create vpc and subnet
resource "huaweicloud_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = var.subnet_name
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

# create security group
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = var.security_group_name
  description = "terraform security group"
}

data "huaweicloud_availability_zones" "zones" {}

data "huaweicloud_compute_flavors" "flavors" {
  availability_zone = data.huaweicloud_availability_zones.zones.names[1]
  performance_type  = "normal"
  cpu_core_count    = 2
}

# create a waf dedicated instance
resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = var.waf_dedicated_instance_name
  available_zone     = data.huaweicloud_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id

  security_group = [
    huaweicloud_networking_secgroup.secgroup.id
  ]
}

resource "huaweicloud_waf_policy" "policy_1" {
  name  = var.waf_policy_name
  level = 1

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1]
  # Make sure that a dedicated instance has been created.
}

# create a dedicated mode domain name
resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain      = var.dedicatd_mode_domain_name
  proxy_id    = huaweicloud_waf_policy.policy_1.id
  keep_policy = true

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.1.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.vpc_1.id
  }
}
