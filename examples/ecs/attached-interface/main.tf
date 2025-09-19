data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  performance_type  = var.instance_performance_type
  cpu_core_count    = var.instance_cpu_core_count
  memory_size       = var.instance_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = length(var.subnet_configurations)

  vpc_id     = huaweicloud_vpc.test.id
  name       = lookup(var.subnet_configurations[count.index], "subnet_name", null)
  cidr       = try(coalesce(lookup(var.subnet_configurations[count.index], "subnet_cidr", null), cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index)), null)
  gateway_ip = try(coalesce(lookup(var.subnet_configurations[count.index], "subnet_gateway_ip", null), cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index), 1)), null)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  image_id           = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  flavor_id          = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  admin_pass         = var.instance_admin_password

  network {
    uuid = try(huaweicloud_vpc_subnet.test[0].id, null)
  }

  # When using `huaweicloud_compute_interface_attach`, if the security group is not specified, the default security group will be automatically added to the ECS instance.
  lifecycle {
    ignore_changes = [
      security_group_ids
    ]
  }
}

resource "huaweicloud_compute_interface_attach" "test" {
  instance_id        = huaweicloud_compute_instance.test.id
  network_id         = var.attached_network_id != "" ? var.attached_network_id : try(huaweicloud_vpc_subnet.test[1].id, null)
  fixed_ip           = var.attached_interface_fixed_ip
  security_group_ids = var.attached_security_group_ids
}
