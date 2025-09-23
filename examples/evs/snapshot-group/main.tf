data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].ids[0], "") : var.instance_flavor_id
  os         = var.instance_image_os_type
  visibility = var.instance_image_visibility
}

# Create a VPC
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

# Create a subnet
resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
}

# Create a security group
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.secgroup_name
  delete_default_rules = true
}

# Create an ECS instance
resource "huaweicloud_compute_instance" "test" {
  name              = var.ecs_instance_name
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  flavor_id         = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, "") : var.instance_flavor_id
  image_id          = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, "") : var.instance_image_id
  security_groups   = [huaweicloud_networking_secgroup.test.name]
  key_pair          = var.key_pair_name
  system_disk_type  = var.system_disk_type
  system_disk_size  = var.system_disk_size

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  volume_type       = var.volume_configuration[count.index].volume_type
  name              = var.volume_configuration[count.index].name
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
}

resource "huaweicloud_compute_volume_attach" "test" {
  count = length(var.volume_configuration)

  instance_id = huaweicloud_compute_instance.test.id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}

resource "huaweicloud_evs_snapshot_group" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  server_id             = huaweicloud_compute_instance.test.id
  volume_ids            = length(huaweicloud_evs_volume.test) > 0 ? try([for v in huaweicloud_evs_volume.test : v.id], null) : null
  instant_access        = var.instant_access
  name                  = var.snapshot_group_name
  description           = var.snapshot_group_description
  enterprise_project_id = var.enterprise_project_id
  incremental           = var.incremental
  tags                  = var.tags
}
