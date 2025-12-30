data "huaweicloud_availability_zones" "test" {
  provider = huaweicloud.sharer
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  provider = huaweicloud.sharer

  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  provider = huaweicloud.sharer

  flavor_id  = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_vpc" "test" {
  provider = huaweicloud.sharer

  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_vpc_subnet" "test" {
  provider = huaweicloud.sharer

  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  provider = huaweicloud.sharer

  name                  = var.security_group_name
  delete_default_rules  = true
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_compute_instance" "test" {
  provider = huaweicloud.sharer

  name                        = var.instance_name
  availability_zone           = try(data.huaweicloud_availability_zones.test.names[0], null)
  flavor_id                   = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  image_id                    = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  security_group_ids          = [huaweicloud_networking_secgroup.test.id]
  admin_pass                  = var.administrator_password
  delete_disks_on_termination = true
  enterprise_project_id       = var.enterprise_project_id

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      availability_zone,
      flavor_id,
      image_id,
      admin_pass,
    ]
  }
}

resource "huaweicloud_evs_volume" "test" {
  provider = huaweicloud.sharer

  name              = var.data_volume_name
  server_id         = huaweicloud_compute_instance.test.id
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  volume_type       = var.data_volume_type
  size              = var.data_volume_size
}

resource "huaweicloud_ims_evs_data_image" "test" {
  provider = huaweicloud.sharer

  name                  = var.data_image_name
  volume_id             = huaweicloud_evs_volume.test.id
  description           = var.data_image_description
  enterprise_project_id = var.enterprise_project_id
}

data "huaweicloud_identity_auth_projects" "test" {
  provider = huaweicloud.accepter
}

locals {
  accepter_project_ids = [for project in data.huaweicloud_identity_auth_projects.test.projects : project.id if project.name == var.region_name]
}

resource "huaweicloud_images_image_share" "test" {
  provider = huaweicloud.sharer

  source_image_id    = huaweicloud_ims_evs_data_image.test.id
  target_project_ids = local.accepter_project_ids
}

resource "huaweicloud_images_image_share_accepter" "test" {
  provider = huaweicloud.accepter

  image_id = huaweicloud_ims_evs_data_image.test.id

  depends_on = [huaweicloud_images_image_share.test]
}

# ST.001 Disable
data "huaweicloud_availability_zones" "accepter" {
  provider = huaweicloud.accepter
}

resource "huaweicloud_vpc" "accepter" {
  provider = huaweicloud.accepter

  name = var.accepter_vpc_name
  cidr = var.accepter_vpc_cidr
}

resource "huaweicloud_vpc_subnet" "accepter" {
  provider = huaweicloud.accepter

  vpc_id     = huaweicloud_vpc.accepter.id
  name       = var.accepter_subnet_name
  cidr       = cidrsubnet(huaweicloud_vpc.accepter.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.accepter.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "accepter" {
  provider = huaweicloud.accepter

  name                 = var.accepter_security_group_name
  delete_default_rules = true
}

data "huaweicloud_compute_flavors" "accepter" {
  count = var.accepter_instance_flavor_id == "" ? 1 : 0

  provider = huaweicloud.accepter

  availability_zone = try(data.huaweicloud_availability_zones.accepter.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "accepter" {
  count = var.accepter_instance_image_id == "" ? 1 : 0

  provider = huaweicloud.accepter

  flavor_id  = var.accepter_instance_flavor_id != "" ? var.accepter_instance_flavor_id : try(data.huaweicloud_compute_flavors.accepter[0].flavors[0].id, null)
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_compute_instance" "accepter" {
  provider = huaweicloud.accepter

  name                        = var.accepter_instance_name
  availability_zone           = try(data.huaweicloud_availability_zones.accepter.names[0], null)
  flavor_id                   = var.accepter_instance_flavor_id != "" ? var.accepter_instance_flavor_id : try(data.huaweicloud_compute_flavors.accepter[0].flavors[0].id, null)
  image_id                    = var.accepter_instance_image_id != "" ? var.accepter_instance_image_id : try(data.huaweicloud_images_images.accepter[0].images[0].id, null)
  security_group_ids          = [huaweicloud_networking_secgroup.accepter.id]
  admin_pass                  = var.administrator_password
  delete_disks_on_termination = true

  network {
    uuid = huaweicloud_vpc_subnet.accepter.id
  }

  depends_on = [huaweicloud_images_image_share_accepter.test]

  lifecycle {
    ignore_changes = [
      availability_zone,
      flavor_id,
      image_id,
      admin_pass,
    ]
  }
}

resource "huaweicloud_evsv3_volume" "accepter" {
  provider = huaweicloud.accepter

  name              = var.accepter_data_volume_name
  availability_zone = try(data.huaweicloud_availability_zones.accepter.names[0], null)
  volume_type       = var.accepter_data_volume_type
  size              = var.accepter_data_volume_size
  image_id          = huaweicloud_images_image_share_accepter.test.image_id
}

resource "huaweicloud_compute_volume_attach" "accepter" {
  provider = huaweicloud.accepter

  instance_id = huaweicloud_compute_instance.accepter.id
  volume_id   = huaweicloud_evsv3_volume.accepter.id
}
# ST.001 Enable
