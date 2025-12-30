data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_deh_types" "test" {
  count = var.deh_instance_host_type == "" ? 1 : 0

  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_deh_instance" "test" {
  name                  = var.deh_instance_name
  availability_zone     = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  host_type             = var.deh_instance_host_type != "" ? var.deh_instance_host_type : try(data.huaweicloud_deh_types.test[0].dedicated_host_types[0].host_type, null)
  auto_placement        = var.deh_instance_auto_placement
  enterprise_project_id = var.enterprise_project_id
  charging_mode         = var.deh_instance_charging_mode
  period_unit           = var.deh_instance_period_unit
  period                = var.deh_instance_period
  auto_renew            = var.deh_instance_auto_renew
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

data "huaweicloud_images_images" "test" {
  count = var.ecs_instance_image_id == "" ? 1 : 0

  flavor_id  = var.ecs_instance_flavor_id != "" ? var.ecs_instance_flavor_id : try(huaweicloud_deh_instance.test.host_properties[0].available_instance_capacities[0].flavor, null)
  visibility = var.ecs_instance_image_visibility
  os         = var.ecs_instance_image_os
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.ecs_instance_name
  availability_zone  = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  flavor_id          = var.ecs_instance_flavor_id != "" ? var.ecs_instance_flavor_id : try(huaweicloud_deh_instance.test.host_properties[0].available_instance_capacities[0].flavor, null)
  image_id           = var.ecs_instance_image_id != "" ? var.ecs_instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, "")
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  admin_pass         = var.ecs_instance_admin_pass

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  scheduler_hints {
    tenancy = "dedicated"
    deh_id  = huaweicloud_deh_instance.test.id
  }
}
