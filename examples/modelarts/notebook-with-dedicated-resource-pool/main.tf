data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

# Make sure open the full ingress access for 111, 2048, 2049, 2051, 2052 and 20048 ports and about TCP and UDP protocols.
resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "111,2048,2049,2051,2052,20048"
}

# ST.001 Disable
resource "huaweicloud_networking_secgroup_rule" "udp_ingress_access" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  ports             = "111,2048,2049,2051,2052,20048"
}
# ST.001 Enable

resource "huaweicloud_sfs_turbo" "test" {
  name                  = var.turbo_name
  size                  = var.turbo_size
  share_proto           = var.turbo_share_proto
  share_type            = var.turbo_share_type
  hpc_bandwidth         = var.turbo_hpc_bandwidth
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = try(data.huaweicloud_availability_zones.test.names[0], null)
  enterprise_project_id = var.enterprise_project_id

  depends_on = [
    huaweicloud_networking_secgroup_rule.test,
    huaweicloud_networking_secgroup_rule.udp_ingress_access,
  ]

  lifecycle {
    ignore_changes = [
      availability_zone,
    ]
  }
}

resource "huaweicloud_modelarts_network" "test" {
  name = var.network_name
  cidr = var.network_cidr

  sfs_turbos {
    name = huaweicloud_sfs_turbo.test.name
    id   = huaweicloud_sfs_turbo.test.id
  }
}

resource "huaweicloud_modelarts_workspace" "test" {
  count = var.workspace_name != "" ? 1 : 0

  name = var.workspace_name
}

data "huaweicloud_modelarts_resource_flavors" "test" {
  count = var.resource_pool_flavor_id != "" ? 0 : 1

  type = "Dedicate"
}

locals {
  available_resource_flavors = [
    for o in try(data.huaweicloud_modelarts_resource_flavors.test[0].flavors, []) : o if lookup(o.az_status, try(data.huaweicloud_availability_zones.test.names[0], null), "soldout") == "normal"
  ]
}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name         = var.resource_pool_name
  scope        = var.resource_pool_scope
  network_id   = huaweicloud_modelarts_network.test.id
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(huaweicloud_modelarts_workspace.test[0].id, null)

  resources {
    flavor_id = var.resource_pool_flavor_id != "" ? var.resource_pool_flavor_id : try(local.available_resource_flavors[0].id, null)
    count     = 1
  }

  # If you want to change the `flavor` or other fields, you need to remove the corresponding fields from "lifecycle.ignore_changes".
  lifecycle {
    ignore_changes = [
      resources,
    ]
  }
}

data "huaweicloud_modelarts_notebook_flavors" "test" {
  count = var.notebook_flavor_id != "" ? 0 : 1

  type     = "DEDICATED"
  category = var.notebook_flavor_category
}

locals {
  available_notebook_flavors = [for o in try(data.huaweicloud_modelarts_notebook_flavors.test[0].flavors, []) : o if !o.sold_out]
}

data "huaweicloud_modelarts_notebook_images" "test" {
  count = var.notebook_image_id != "" ? 0 : 1

  type     = var.notebook_image_type
  cpu_arch = try(data.huaweicloud_modelarts_notebook_flavors.test[0].flavors[0].arch, "x86_64")
}

locals {
  available_notebook_images = [
    for o in try(data.huaweicloud_modelarts_notebook_images.test[0].images, []) : o if contains(o.resource_categories, var.notebook_flavor_category) && o.status == "ACTIVE" && contains(o.dev_services, "NOTEBOOK")
  ]
}

resource "huaweicloud_kps_keypair" "test" {
  count = var.notebook_key_pair_name == "" && length(var.allowed_access_ips) > 0 ? 1 : 0

  name = var.keypair_name
}

data "huaweicloud_modelartsv2_resource_pools" "test" {}

locals {
  resource_pool = try([for pool in data.huaweicloud_modelartsv2_resource_pools.test.resource_pools : pool if pool.metadata[0].name ==
  huaweicloud_modelarts_resource_pool.test.id][0], {})
}

resource "huaweicloud_modelarts_notebook" "test" {
  name               = var.notebook_name
  flavor_id          = var.notebook_flavor_id != "" ? var.notebook_flavor_id : try(local.available_notebook_flavors[0].id, null)
  image_id           = var.notebook_image_id != "" ? var.notebook_image_id : try(local.available_notebook_images[0].id, null)
  description        = var.notebook_description
  pool_id            = huaweicloud_modelarts_resource_pool.test.id
  workspace_id       = try(local.resource_pool.metadata[0].labels["os.modelarts/workspace.id"], null)
  key_pair           = var.notebook_key_pair_name != "" ? var.notebook_key_pair_name : try(huaweicloud_kps_keypair.test[0].name, null)
  allowed_access_ips = length(var.allowed_access_ips) > 0 ? var.allowed_access_ips : null

  volume {
    type      = "EFS"
    ownership = "DEDICATED"
    uri       = format("%s:/", try(huaweicloud_modelarts_network.test.sfs_turbos[0].uri, null))
    id        = huaweicloud_sfs_turbo.test.id
  }

  tags = var.notebook_tags
}

resource "huaweicloud_modelarts_notebook_mount_storage" "test" {
  count = var.notebook_mount_storage_path != "" && var.notebook_mount_storage_local_directory != "" ? 1 : 0

  notebook_id           = huaweicloud_modelarts_notebook.test.id
  storage_path          = var.notebook_mount_storage_path
  local_mount_directory = var.notebook_mount_storage_local_directory
}
