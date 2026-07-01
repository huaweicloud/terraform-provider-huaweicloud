data "huaweicloud_modelartsv2_resource_pools" "test" {}

locals {
  resource_pool = try([for pool in data.huaweicloud_modelartsv2_resource_pools.test.resource_pools : pool if pool.metadata[0].name ==
  var.service_group_pool_id][0], {})
}

resource "huaweicloud_modelartsv2_service" "test" {
  name         = var.service_name
  version      = var.service_version
  type         = var.service_type
  workspace_id = try(local.resource_pool.metadata[0].labels["os.modelarts/workspace.id"], null)
  deploy_type  = var.service_deploy_type
  description  = var.service_description

  group_configs {
    framework = var.service_group_framework
    name      = var.service_group_name
    pool_id   = try(local.resource_pool.metadata[0].name, null)
    weight    = var.service_group_weight
    count     = var.service_group_count

    dynamic "unit_configs" {
      for_each = var.service_unit_configs

      content {
        image {
          source   = unit_configs.value.image.source
          swr_path = unit_configs.value.image.swr_path
          id       = unit_configs.value.image.id
        }

        dynamic "custom_spec" {
          for_each = unit_configs.value.custom_spec != null ? [unit_configs.value.custom_spec] : []

          content {
            memory = custom_spec.value.memory
            cpu    = custom_spec.value.cpu
            gpu    = custom_spec.value.gpu
            ascend = custom_spec.value.ascend
          }
        }

        dynamic "models" {
          for_each = unit_configs.value.models

          content {
            source     = models.value.source
            mount_path = models.value.mount_path
            address    = models.value.address
            source_id  = models.value.source_id
          }
        }

        dynamic "codes" {
          for_each = unit_configs.value.codes

          content {
            source     = codes.value.source
            mount_path = codes.value.mount_path
            address    = codes.value.address
            source_id  = codes.value.source_id
          }
        }

        dynamic "readiness_health" {
          for_each = unit_configs.value.readiness_health != null ? [unit_configs.value.readiness_health] : []

          content {
            initial_delay_seconds = readiness_health.value.initial_delay_seconds
            timeout_seconds       = readiness_health.value.timeout_seconds
            period_seconds        = readiness_health.value.period_seconds
            failure_threshold     = readiness_health.value.failure_threshold
            check_method          = readiness_health.value.check_method
            command               = readiness_health.value.command
            url                   = readiness_health.value.url
          }
        }

        dynamic "startup_health" {
          for_each = unit_configs.value.startup_health != null ? [unit_configs.value.startup_health] : []

          content {
            initial_delay_seconds = startup_health.value.initial_delay_seconds
            timeout_seconds       = startup_health.value.timeout_seconds
            period_seconds        = startup_health.value.period_seconds
            failure_threshold     = startup_health.value.failure_threshold
            check_method          = startup_health.value.check_method
            command               = startup_health.value.command
            url                   = startup_health.value.url
          }
        }

        dynamic "liveness_health" {
          for_each = unit_configs.value.liveness_health != null ? [unit_configs.value.liveness_health] : []

          content {
            initial_delay_seconds = liveness_health.value.initial_delay_seconds
            timeout_seconds       = liveness_health.value.timeout_seconds
            period_seconds        = liveness_health.value.period_seconds
            failure_threshold     = liveness_health.value.failure_threshold
            check_method          = liveness_health.value.check_method
            command               = liveness_health.value.command
            url                   = liveness_health.value.url
          }
        }

        role     = unit_configs.value.role
        flavor   = unit_configs.value.flavor
        count    = unit_configs.value.count
        cmd      = unit_configs.value.cmd
        recovery = unit_configs.value.recovery
        envs     = unit_configs.value.envs
        port     = unit_configs.value.port
      }
    }
  }

  runtime_config = var.service_runtime_config
  upgrade_config = var.service_upgrade_config

  dynamic "log_configs" {
    for_each = var.service_log_configs

    content {
      type          = log_configs.value.type
      log_group_id  = log_configs.value.log_group_id
      log_stream_id = log_configs.value.log_stream_id
    }
  }

  tags = var.service_tags
}
