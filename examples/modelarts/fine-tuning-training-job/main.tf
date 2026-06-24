resource "huaweicloud_modelarts_workspace" "test" {
  count = var.workspace_name != "" ? 1 : 0

  name = var.workspace_name
}

resource "huaweicloud_smn_topic" "test" {
  count = var.topic_name != "" ? 1 : 0

  name = var.topic_name
}

resource "huaweicloud_modelarts_training_job" "test" {
  kind       = "job"
  train_type = var.training_job_train_type

  metadata {
    name         = var.training_job_name
    workspace_id = var.workspace_id != "" ? var.workspace_id : try(huaweicloud_modelarts_workspace.test[0].id, null)
    annotations  = var.training_job_annotations
    description  = var.training_job_description
  }

  algorithm {
    dynamic "inputs" {
      for_each = var.training_job_inputs

      content {
        local_dir = inputs.value.local_dir

        remote {
          dataset {
            id                 = inputs.value.dataset.id
            name               = inputs.value.dataset.name
            version_id         = inputs.value.dataset.version_id
            service_type       = inputs.value.dataset.service_type
            dataset_proportion = inputs.value.dataset.dataset_proportion
          }
        }
      }
    }

    environments = var.training_job_environments
  }

  spec {
    resource {
      flavor_id  = var.resource_flavor_id
      node_count = var.resource_node_count
    }

    asset_id = var.training_job_asset_id

    dynamic "asset_model" {
      for_each = [var.training_job_asset_model]

      content {
        name    = asset_model.value.name
        version = asset_model.value.version
        type    = asset_model.value.type
        code    = asset_model.value.code
        desc    = asset_model.value.desc
        series  = asset_model.value.series
      }
    }

    dynamic "notification" {
      for_each = var.training_job_notification_topic_urn != "" || var.topic_name != "" ? [1] : []

      content {
        topic_urn = var.training_job_notification_topic_urn != "" ? var.training_job_notification_topic_urn : var.topic_name != "" ? try(huaweicloud_smn_topic.test[0].id, null) : ""
        events    = var.training_job_notification_events
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

  ftjob_config {
    task_env {
      dynamic "envs" {
        for_each = var.training_job_ftjob_config.envs

        content {
          label       = envs.value.label
          des         = envs.value.des
          env_name    = envs.value.env_name
          env_type    = envs.value.env_type
          value       = envs.value.value
          modifiable  = envs.value.modifiable
          displayable = envs.value.displayable
        }
      }
    }

    dynamic "checkpoint_config" {
      for_each = var.training_job_ftjob_config.checkpoint_config != null ? [var.training_job_ftjob_config.checkpoint_config] : []

      content {
        checkpoint_id        = checkpoint_config.value.checkpoint_id
        save_checkpoints_max = checkpoint_config.value.save_checkpoints_max
        skipped_steps        = checkpoint_config.value.skipped_steps
        restore_training     = checkpoint_config.value.restore_training
      }
    }
  }

  tags = var.training_job_tags
}
