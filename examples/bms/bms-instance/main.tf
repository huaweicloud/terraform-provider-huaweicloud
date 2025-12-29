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

data "huaweicloud_bms_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  cpu_arch          = "aarch64"
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  os         = "Huawei Cloud EulerOS"
  image_type = "Ironic"
}

resource "huaweicloud_kps_keypair" "test" {
  name = var.keypair_name
}

resource "huaweicloud_bms_instance" "test" {
  name                  = var.instance_name
  user_id               = var.instance_user_id
  availability_zone     = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  vpc_id                = huaweicloud_vpc.test.id
  flavor_id             = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_bms_flavors.test[0].flavors[0].id, null)
  image_id              = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  security_groups       = [huaweicloud_networking_secgroup.test.id]
  key_pair              = huaweicloud_kps_keypair.test.name
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.instance_tags
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  # If you want to change some of the following parameters, you need to remove the corresponding fields from "lifecycle.ignore_changes".
  lifecycle {
    ignore_changes = [
      availability_zone,
      flavor_id,
      image_id,
    ]
  }
}
