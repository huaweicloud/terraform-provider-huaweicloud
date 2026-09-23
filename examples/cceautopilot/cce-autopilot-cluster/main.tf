resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
}

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = var.cluster_name
  flavor      = "cce.autopilot.cluster"
  alias       = var.cluster_alias
  description = var.cluster_description
  category    = var.cluster_category
  type        = var.cluster_type
  annotations = var.cluster_annotations
  custom_san  = var.cluster_custom_san
  eip_id      = var.cluster_eip_id

  enable_snat             = var.cluster_enable_snat
  enable_swr_image_access = var.cluster_enable_swr_image_access

  delete_efs   = var.cluster_delete_efs
  delete_eni   = var.cluster_delete_eni
  delete_net   = var.cluster_delete_net
  delete_obs   = var.cluster_delete_obs
  delete_sfs30 = var.cluster_delete_sfs_turbo

  lts_reclaim_policy = var.cluster_lts_reclaim_policy

  host_network {
    vpc    = huaweicloud_vpc.test.id
    subnet = huaweicloud_vpc_subnet.test.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
    }
  }

  service_network {
    ipv4_cidr = var.cluster_service_network_cidr
  }

  extend_param {
    enterprise_project_id = var.cluster_enterprise_project_id
  }

  dynamic "configurations_override" {
    for_each = length(var.cluster_configurations_override) > 0 ? [1] : []

    content {
      name = var.cluster_configurations_override_name

      dynamic "configurations" {
        for_each = var.cluster_configurations_override

        content {
          name  = configurations.value.name
          value = configurations.value.value
        }
      }
    }
  }

  tags = var.cluster_tags

  lifecycle {
    ignore_changes = [
      annotations,
      tags["CCE-Cluster-ID"],
      service_network,
    ]
  }
}
