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
  cidr = var.network_cidr # The recommended connecting CIDR about SFS Turbo.

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

resource "huaweicloud_smn_topic" "test" {
  count = var.topic_name != "" ? 1 : 0

  name = var.topic_name
}

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name         = var.training_job_name
    workspace_id = var.workspace_id != "" ? var.workspace_id : try(huaweicloud_modelarts_workspace.test[0].id, null)
    annotations  = var.training_job_annotations
    description  = var.training_job_description
  }

  algorithm {
    code_dir = var.training_job_code_dir
    command  = var.training_job_command

    engine {
      image_url      = var.training_job_engine.image_url
      engine_id      = var.training_job_engine.id
      engine_version = var.training_job_engine.version
      engine_name    = var.training_job_engine.name
    }


    dynamic "inputs" {
      for_each = var.training_job_inputs

      content {
        local_dir = inputs.value.local_dir

        remote {
          dataset {
            id           = inputs.value.dataset.id
            name         = inputs.value.dataset.name
            version_id   = inputs.value.dataset.version_id
            service_type = inputs.value.dataset.service_type
          }
        }
      }
    }

    environments = var.training_job_environments
  }

  spec {
    resource {
      pool_id    = huaweicloud_modelarts_resource_pool.test.id
      node_count = var.resource_node_count
    }

    dynamic "volumes" {
      for_each = var.training_job_volumes

      content {
        dynamic "nfs" {
          for_each = volumes.value.nfs != null ? [volumes.value.nfs] : []

          content {
            nfs_server_path = nfs.value.nfs_server_path
            local_path      = nfs.value.local_path
            read_only       = nfs.value.read_only
          }
        }

        dynamic "pfs" {
          for_each = volumes.value.pfs != null ? [volumes.value.pfs] : []

          content {
            pfs_path   = pfs.value.pfs_path
            local_path = pfs.value.local_path
          }
        }

        dynamic "obs" {
          for_each = volumes.value.obs != null ? [volumes.value.obs] : []

          content {
            obs_path   = obs.value.obs_path
            local_path = obs.value.local_path
          }
        }
      }
    }

    dynamic "log_export_path" {
      for_each = var.training_job_log_export_path_obs_url != "" ? [1] : []

      content {
        obs_url = var.training_job_log_export_path_obs_url
      }
    }

    dynamic "log_export_config" {
      for_each = var.training_job_log_export_config_version != "" ? [1] : []

      content {
        version = var.training_job_log_export_config_version
      }
    }

    dynamic "auto_stop" {
      for_each = var.training_job_auto_stop_duration != 0 ? [1] : []

      content {
        time_unit = "HOURS"
        duration  = var.training_job_auto_stop_duration
      }
    }

    dynamic "notification" {
      for_each = var.training_job_notification_topic_urn != "" || var.topic_name != "" ? [1] : []

      content {
        topic_urn = var.training_job_notification_topic_urn != "" ? var.training_job_notification_topic_urn : var.topic_name != "" ? huaweicloud_smn_topic.test[0].id : ""
        events    = var.training_job_notification_events
      }
    }

    dynamic "custom_metrics" {
      for_each = var.custom_metrics != null ? [var.custom_metrics] : []

      content {
        dynamic "exec" {
          for_each = custom_metrics.value.exec != null ? [custom_metrics.value.exec] : []

          content {
            command = exec.value.command
          }
        }

        dynamic "http_get" {
          for_each = custom_metrics.value.http_get != null ? [custom_metrics.value.http_get] : []
          content {
            path = http_get.value.path
            port = http_get.value.port
          }
        }
      }
    }

    dynamic "asset_model" {
      for_each = var.training_job_asset_model != null ? [var.training_job_asset_model] : []

      content {
        name    = asset_model.value.name
        version = asset_model.value.version
        type    = asset_model.value.type
        code    = asset_model.value.code
        desc    = asset_model.value.desc
        series  = asset_model.value.series
      }
    }

    dynamic "output_model" {
      for_each = var.training_job_output_model != null ? [var.training_job_output_model] : []

      content {
        obs {
          obs_path   = output_model.value.obs_path
          local_path = output_model.value.local_path
        }
      }
    }
  }

  tags = var.training_job_tags
}
