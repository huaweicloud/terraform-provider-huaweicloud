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

  name = var.vpc_name
  cidr = var.vpc_cidr
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

  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_compute_instance" "test" {
  provider = huaweicloud.sharer

  name               = var.instance_name
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], null)
  flavor_id          = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  image_id           = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  admin_pass         = var.administrator_password

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  delete_disks_on_termination = true

  dynamic "data_disks" {
    for_each = var.instance_data_disks

    content {
      type = data_disks.value.type
      size = data_disks.value.size
    }
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

resource "huaweicloud_cbr_vault" "test" {
  provider = huaweicloud.sharer

  name             = var.vault_name
  type             = var.vault_type
  consistent_level = var.vault_consistent_level
  protection_type  = var.vault_protection_type
  size             = var.vault_size
}

resource "huaweicloud_ims_ecs_whole_image" "test" {
  provider = huaweicloud.sharer

  name        = var.whole_image_name
  instance_id = huaweicloud_compute_instance.test.id
  vault_id    = huaweicloud_cbr_vault.test.id
  description = var.whole_image_description
}

resource "huaweicloud_images_image_share" "test" {
  provider = huaweicloud.sharer

  source_image_id    = huaweicloud_ims_ecs_whole_image.test.id
  target_project_ids = var.accepter_project_ids
}

# ST.001 Disable
resource "huaweicloud_cbr_vault" "accepter" {
  provider = huaweicloud.accepter

  name             = var.accepter_vault_name
  type             = var.accepter_vault_type
  consistent_level = var.accepter_vault_consistent_level
  protection_type  = var.accepter_vault_protection_type
  size             = var.accepter_vault_size
}

resource "huaweicloud_images_image_share_accepter" "accepter" {
  provider = huaweicloud.accepter

  image_id = huaweicloud_ims_ecs_whole_image.test.id
  vault_id = huaweicloud_cbr_vault.accepter.id

  depends_on = [huaweicloud_images_image_share.test]
}

data "huaweicloud_availability_zones" "accepter" {
  provider = huaweicloud.accepter
}

data "huaweicloud_compute_flavors" "accepter" {
  provider = huaweicloud.accepter

  availability_zone = try(data.huaweicloud_availability_zones.accepter.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

resource "huaweicloud_vpc" "accepter" {
  provider = huaweicloud.accepter

  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "accepter" {
  provider = huaweicloud.accepter

  vpc_id     = huaweicloud_vpc.accepter.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.accepter.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.accepter.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "accepter" {
  provider = huaweicloud.accepter

  name                 = var.security_group_name
  delete_default_rules = true
}

# Create new ECS instance in accepter account using the accepted shared whole image
resource "huaweicloud_compute_instance" "accepter" {
  provider = huaweicloud.accepter

  depends_on = [huaweicloud_images_image_share_accepter.accepter]

  name               = var.accepter_instance_name
  availability_zone  = try(data.huaweicloud_availability_zones.accepter.names[0], null)
  flavor_id          = try(data.huaweicloud_compute_flavors.accepter.flavors[0].id, null)
  image_id           = huaweicloud_images_image_share_accepter.accepter.image_id
  security_group_ids = [huaweicloud_networking_secgroup.accepter.id]
  admin_pass         = var.administrator_password

  network {
    uuid = huaweicloud_vpc_subnet.accepter.id
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
# ST.001 Enable
