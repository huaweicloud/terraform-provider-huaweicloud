data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  count = var.vpc_id == "" && var.subnet_id == "" ? 1 : 0

  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = var.subnet_id == "" ? 1 : 0

  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  name              = var.subnet_name
  cidr              = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : var.subnet_cidr != "" ? cidrhost(var.subnet_cidr, 1) : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0), 1)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
}

resource "huaweicloud_vpc_eip" "test" {
  count = var.eip_address == "" ? 1 : 0

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = var.cluster_name
  flavor_id              = var.cluster_flavor_id
  cluster_version        = var.cluster_version
  cluster_type           = var.cluster_type
  container_network_type = var.container_network_type
  vpc_id                 = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  subnet_id              = var.subnet_id != "" ? var.subnet_id : huaweicloud_vpc_subnet.test[0].id
  eip                    = var.eip_address != "" ? var.eip_address : huaweicloud_vpc_eip.test[0].address
}

data "huaweicloud_compute_flavors" "test" {
  performance_type  = var.node_performance_type
  cpu_core_count    = var.node_cpu_core_count
  memory_size       = var.node_memory_size
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
}

resource "huaweicloud_kps_keypair" "test" {
  name = var.keypair_name
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  type                     = var.node_pool_type
  name                     = var.node_pool_name
  flavor_id                = try(data.huaweicloud_compute_flavors.test.flavors[0].id, null)
  availability_zone        = try(data.huaweicloud_availability_zones.test.names[0], null)
  os                       = var.node_pool_os_type
  initial_node_count       = var.node_pool_initial_node_count
  min_node_count           = var.node_pool_min_node_count
  max_node_count           = var.node_pool_max_node_count
  scale_down_cooldown_time = var.node_pool_scale_down_cooldown_time
  priority                 = var.node_pool_priority
  key_pair                 = huaweicloud_kps_keypair.test.name
  tags                     = var.node_pool_tags

  root_volume {
    volumetype = var.root_volume_type
    size       = var.root_volume_size
  }

  dynamic "data_volumes" {
    for_each = local.flattened_data_volumes

    content {
      volumetype    = data_volumes.value.volumetype
      size          = data_volumes.value.size
      kms_key_id    = data_volumes.value.kms_key_id
      extend_params = data_volumes.value.extend_params
    }
  }

  dynamic "storage" {
    for_each = length(local.default_data_volumes_configuration_with_virtual_spaces) + length(local.user_data_volumes_configuration_with_virtual_spaces) > 0 ? [1] : []

    content {
      # Default selector which is used to select the data volumes that used by CCE
      dynamic "selectors" {
        for_each = local.default_data_volumes_configuration_with_virtual_spaces

        content {
          name                           = "cceUse"
          type                           = "evs"
          match_label_volume_type        = selectors.value.volumetype
          match_label_size               = selectors.value.size
          match_label_count              = selectors.value.count
          match_label_metadata_encrypted = selectors.value.kms_key_id != "" && selectors.value.kms_key_id != null ? "1" : "0"
          match_label_metadata_cmkid     = selectors.value.kms_key_id != "" ? selectors.value.kms_key_id : null
        }
      }

      # User data selector which is used to select the data volumes that used by user
      dynamic "selectors" {
        for_each = local.user_data_volumes_configuration_with_virtual_spaces

        content {
          name                           = selectors.value.select_name
          type                           = "evs"
          match_label_volume_type        = selectors.value.volumetype
          match_label_size               = selectors.value.size
          match_label_count              = selectors.value.count
          match_label_metadata_encrypted = selectors.value.kms_key_id != "" && selectors.value.kms_key_id != null ? "1" : "0"
          match_label_metadata_cmkid     = selectors.value.kms_key_id != "" ? selectors.value.kms_key_id : null
        }
      }

      # Default group which is used to group the data volumes that used by CCE
      dynamic "groups" {
        for_each = local.default_data_volumes_configuration_with_virtual_spaces

        content {
          name           = "vgpaas"
          cce_managed    = true
          selector_names = ["cceUse"]

          dynamic "virtual_spaces" {
            for_each = groups.value.virtual_spaces

            content {
              name            = virtual_spaces.value.name
              size            = virtual_spaces.value.size
              lvm_lv_type     = virtual_spaces.value.lvm_lv_type
              lvm_path        = virtual_spaces.value.lvm_path
              runtime_lv_type = virtual_spaces.value.runtime_lv_type
            }
          }
        }
      }

      # User data group which is used to group the data volumes that used by user
      dynamic "groups" {
        for_each = local.user_data_volumes_configuration_with_virtual_spaces

        content {
          name           = "vg${groups.value.select_name}"
          selector_names = [groups.value.select_name]

          dynamic "virtual_spaces" {
            for_each = groups.value.virtual_spaces

            content {
              name            = virtual_spaces.value.name
              size            = virtual_spaces.value.size
              lvm_lv_type     = virtual_spaces.value.lvm_lv_type
              lvm_path        = virtual_spaces.value.lvm_path
              runtime_lv_type = virtual_spaces.value.runtime_lv_type
            }
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
      initial_node_count
    ]
  }

  depends_on = [
    huaweicloud_kps_keypair.test
  ]
}

locals {
  flattened_data_volumes                                 = flatten([
    for v in var.data_volumes_configuration : [
      for i in range(v.count) : {
        volumetype    = v.volumetype
        size          = v.size
        kms_key_id    = v.kms_key_id
        extend_params = v.extend_params
      }
    ]
  ])
  default_data_volumes_configuration_with_virtual_spaces = [
    for v in slice(var.data_volumes_configuration, 0, 1) : v if length(v.virtual_spaces) > 0
  ]
  user_data_volumes_configuration_with_virtual_spaces    = [
    for i, v in  [
      for v in slice(var.data_volumes_configuration, 1, length(var.data_volumes_configuration)) : v if length(v.virtual_spaces) > 0
    ] : {
      select_name    = "user${i+1}"
      volumetype     = v.volumetype
      size           = v.size
      count          = v.count
      kms_key_id     = v.kms_key_id
      extend_params  = v.extend_params
      virtual_spaces = v.virtual_spaces
    }
  ]
}
