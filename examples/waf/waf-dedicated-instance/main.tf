data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = var.subnet_name
  cidr              = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip        = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_waf_dedicated_instance" "test" {
  name                  = var.waf_dedicated_instance_name
  available_zone        = data.huaweicloud_availability_zones.test.names[0]
  specification_code    = var.waf_dedicated_instance_specification_code
  ecs_flavor            = data.huaweicloud_compute_flavors.test.flavors[0].id
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = var.enterprise_project_id

  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

resource "huaweicloud_waf_policy" "test" {
  name  = var.waf_policy_name
  level = 1

  # Make sure that a dedicated instance has been created.
  depends_on = [huaweicloud_waf_dedicated_instance.test]
}
