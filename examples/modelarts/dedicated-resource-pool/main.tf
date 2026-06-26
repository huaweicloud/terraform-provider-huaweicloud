data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  count = var.turbo_name != "" ? 1 : 0

  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_vpc_subnet" "test" {
  count = var.turbo_name != "" ? 1 : 0

  vpc_id     = huaweicloud_vpc.test[0].id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  count = var.turbo_name != "" ? 1 : 0

  name                 = var.security_group_name
  delete_default_rules = true
}

# Make sure open the full ingress access for 111, 2048, 2049, 2051, 2052 and 20048 ports and about TCP and UDP protocols.
resource "huaweicloud_networking_secgroup_rule" "test" {
  count = var.turbo_name != "" ? 1 : 0

  security_group_id = huaweicloud_networking_secgroup.test[0].id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "111,2048,2049,2051,2052,20048"
}

# ST.001 Disable
resource "huaweicloud_networking_secgroup_rule" "udp_ingress_access" {
  count = var.turbo_name != "" ? 1 : 0

  security_group_id = huaweicloud_networking_secgroup.test[0].id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  ports             = "111,2048,2049,2051,2052,20048"
}
# ST.001 Enable

resource "huaweicloud_sfs_turbo" "test" {
  count = var.turbo_name != "" ? 1 : 0

  name                  = var.turbo_name
  size                  = var.turbo_size
  share_proto           = var.turbo_share_proto
  share_type            = var.turbo_share_type
  hpc_bandwidth         = var.turbo_hpc_bandwidth
  vpc_id                = huaweicloud_vpc.test[0].id
  subnet_id             = huaweicloud_vpc_subnet.test[0].id
  security_group_id     = huaweicloud_networking_secgroup.test[0].id
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

locals {
  network_sfs_turbos = var.turbo_name != "" ? [
    {
      id   = huaweicloud_sfs_turbo.test[0].id
      name = huaweicloud_sfs_turbo.test[0].name
    }
  ] : var.network_sfs_turbos
}

resource "huaweicloud_modelarts_network" "test" {
  count = var.network_name != "" ? 1 : 0

  name = var.network_name
  cidr = var.network_cidr

  dynamic "sfs_turbos" {
    for_each = local.network_sfs_turbos

    content {
      id   = sfs_turbos.value.id
      name = sfs_turbos.value.name
    }
  }
}

resource "huaweicloud_modelarts_workspace" "test" {
  count = var.workspace_name != "" ? 1 : 0

  name = var.workspace_name
}

locals {
  is_query_resource_flavor = length([for r in var.resource_pool_resources : r if r.flavor_id == ""]) > 0
}

data "huaweicloud_modelarts_resource_flavors" "test" {
  count = local.is_query_resource_flavor ? 1 : 0

  type = "Dedicate"
}

locals {
  available_resource_flavors = [
    for o in try(data.huaweicloud_modelarts_resource_flavors.test[0].flavors, []) : o if lookup(o.az_status, try(data.huaweicloud_availability_zones.test.names[0], null), "soldout") == "normal"
  ]
}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name         = var.resource_pool_name
  description  = var.resource_pool_description
  scope        = var.resource_pool_scope
  network_id   = var.network_id != "" ? var.network_id : huaweicloud_modelarts_network.test[0].id
  workspace_id = var.workspace_id != "" ? var.workspace_id : try(huaweicloud_modelarts_workspace.test[0].id, null)

  dynamic "metadata" {
    for_each = var.resource_pool_metadata_annotations != null ? [1] : []

    content {
      annotations = var.resource_pool_metadata_annotations
    }
  }

  dynamic "resources" {
    for_each = var.resource_pool_resources

    content {
      flavor_id     = resources.value.flavor_id != "" ? resources.value.flavor_id : try(local.available_resource_flavors[0].id, null)
      count         = resources.value.count
      max_count     = resources.value.max_count
      extend_params = resources.value.extend_params

      dynamic "root_volume" {
        for_each = try(resources.value.root_volume, null) != null ? [resources.value.root_volume] : []

        content {
          volume_type = root_volume.value.volume_type
          size        = root_volume.value.size
        }
      }

      dynamic "data_volumes" {
        for_each = try(resources.value.data_volumes, [])

        content {
          volume_type   = data_volumes.value.volume_type
          size          = data_volumes.value.size
          extend_params = data_volumes.value.extend_params
          count         = try(data_volumes.value.count, null)
        }
      }

      dynamic "volume_group_configs" {
        for_each = try(resources.value.volume_group_configs, [])

        content {
          volume_group     = volume_group_configs.value.volume_group
          docker_thin_pool = try(volume_group_configs.value.docker_thin_pool, null)
          types            = try(volume_group_configs.value.types, null)

          dynamic "lvm_config" {
            for_each = volume_group_configs.value.lvm_config != null ? [volume_group_configs.value.lvm_config] : []

            content {
              lv_type = lvm_config.value.lv_type
              path    = try(lvm_config.value.path, null)
            }
          }
        }
      }

      dynamic "os" {
        for_each = resources.value.os != null ? [resources.value.os] : []

        content {
          name       = try(os.value.name, null)
          image_id   = try(os.value.image_id, null)
          image_type = try(os.value.image_type, null)
        }
      }

      dynamic "driver" {
        for_each = resources.value.driver != null ? [resources.value.driver] : []

        content {
          version = driver.value.version
        }
      }

      dynamic "creating_step" {
        for_each = resources.value.creating_step != null ? [resources.value.creating_step] : []

        content {
          step = creating_step.value.step
          type = creating_step.value.type
        }
      }
    }
  }
}
