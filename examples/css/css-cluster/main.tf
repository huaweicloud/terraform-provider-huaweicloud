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
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

data "huaweicloud_css_flavors" "test" {
  count = var.cluster_flavor == "" ? 1 : 0
}

resource "huaweicloud_css_cluster" "test" {
  name              = var.cluster_name
  engine_version    = var.cluster_engine_version
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  ess_node_config {
    flavor          = var.cluster_flavor == "" ? try(data.huaweicloud_css_flavors.test[0].flavors[0].name, null) : var.cluster_flavor
    instance_number = var.cluster_instance_number

    volume {
      volume_type = var.cluster_volume_type
      size        = var.cluster_volume_size
    }
  }

  engine_type   = var.cluster_engine_type
  security_mode = var.cluster_security_mode
  password      = var.cluster_access_password
  https_enabled = var.cluster_https_enabled
  charging_mode = var.charging_mode
  period_unit   = var.period_unit
  period        = var.period
  auto_renew    = var.auto_renew
}
