---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_job"
description: |-
  Manages a ModelArts training job resource within HuaweiCloud.
---

# huaweicloud_modelarts_training_job

Manages a ModelArts training job resource within HuaweiCloud.

## Example Usage

### Create a training job with existing algorithm

```hcl
variable "training_job_name" {}
variable "algorithm_id" {}
variable "inputs" {
  type = list(object({
    remote = object({
      obs = optional(object({
        obs_url = string
      }), null)

      dataset = optional(object({
        id           = string
        name         = optional(string)
        version_id   = optional(string)
        service_type = optional(string)
      }))
    })

    name          = optional(string)
    local_dir     = optional(string)
    access_method = optional(string)
    description   = optional(string)
  }))

  default = []
}
variable "outputs" {
  type = list(object({
    name           = string
    remote_obs_url = string
    local_dir      = optional(string)
    access_method  = optional(string)
    description    = optional(string)
  }))

  default = []
}
variable "parameters" {
  type = list(object({
    name  = string
    value = string

    constraint = optional(object({
      type        = string
      editable    = optional(bool)
      required    = optional(bool)
      valid_type  = optional(string)
      valid_range = optional(list(string))
    }))
  }))

  default = []
}
variable "resource_pool_flavor_id" {}
variable "resource_pool_node_count" {
  type = number
}

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name = var.training_job_name
  }

  algorithm {
    id = var.algorithm_id

    dynamic "inputs" {
      for_each = var.inputs

      content {
        name          = inputs.value.name
        description   = inputs.value.description
        access_method = inputs.value.access_method
        local_dir     = inputs.value.local_dir

        dynamic "remote" {
          for_each = [inputs.value.remote]

          content {
            dynamic "obs" {
              for_each = inputs.value.remote.obs != null ? [inputs.value.remote.obs] : []

              content {
                obs_url = obs.value.obs_url
              }
            }

            dynamic "dataset" {
              for_each = inputs.value.remote.dataset != null ? [inputs.value.remote.dataset] : []

              content {
                id           = dataset.value.id
                name         = dataset.value.name
                version_id   = dataset.value.version_id
                service_type = dataset.value.service_type
              }
            }
          }
        }
      }
    }

    dynamic "outputs" {
      for_each = var.outputs

      content {
        name          = outputs.value.name
        local_dir     = outputs.value.local_dir
        access_method = outputs.value.access_method
        description   = outputs.value.description

        remote {
          obs {
            obs_url = outputs.value.remote_obs_url
          }
        }
      }
    }

    dynamic "parameters" {
      for_each = var.parameters

      content {
        name  = parameters.value.name
        value = parameters.value.value

        dynamic "constraint" {
          for_each = parameters.value.constraint != null ? [parameters.value.constraint] : []

          content {
            type        = constraint.value.type
            editable    = constraint.value.editable
            required    = constraint.value.required
            valid_type  = constraint.value.valid_type
            valid_range = constraint.value.valid_range
          }
        }
      }
    }
  }

  spec {
    resource {
      flavor_id  = var.resource_pool_flavor_id
      node_count = var.resource_pool_node_count
    }
  }
}
```

### Create a training job with custom algorithm and dedicated resource pool

```hcl
variable "training_job_name" {}
variable "code_dir" {}
variable "boot_file" {}
variable "command" {}
variable "image_url" {}
variable "resource_pool_id" {}
variable "resource_node_count" {
  type = number
}
variable "main_container_customized_flavor" {
  type = optional(object({
    cpu_core_num    = optional(number)
    mem_size        = optional(number)
    accelerator_num = optional(number)
  }))
}

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name = var.training_job_name
  }

  algorithm {
    code_dir  = var.code_dir
    boot_file = var.boot_file
    command   = var.command

    engine {
      image_url = var.image_url
    }
  }

  spec {
    resource {
      pool_id    = var.resource_pool_id
      node_count = var.resource_node_count

      dynamic "main_container_customized_flavor" {
        for_each = var.main_container_customized_flavor != null ? [var.main_container_customized_flavor] : []

        content {
          cpu_core_num    = var.main_container_customized_flavor.cpu_core_num
          mem_size        = var.main_container_customized_flavor.mem_size
          accelerator_num = var.main_container_customized_flavor.accelerator_num
        }
      }

    }
  }
}
```

### Create a fine-tuning job and publish to asset

```hcl
variable "training_job_name" {}
variable "inputs" {
  type = list(object({
    remote = object({
      dataset = object({
        id                 = string
        name               = string
        dataset_proportion = optional(number)
      })
    })
  }))
}
variable "resource_pool_flavor_id" {}
variable "resource_pool_node_count" {}
variable "asset_id" {}
variable "output_model" {
  type = object({
    obs = object({
      obs_path   = string
      local_path = optional(string)
    })
  })
}
variable "asset_model" {
  type = object({
    name    = string
    version = string
    type    = string
    code    = optional(string)
    desc    = optional(string)
    series  = optional(string)
  })
}
variable "ftjob_config" {
  type = object({
    task_env = object({
      envs = list(object({
        label       = optional(string)
        des         = optional(string)
        env_name    = optional(string)
        env_type    = optional(string)
        value       = optional(string)
        modifiable  = optional(bool)
        displayable = optional(bool)
      }))
    })
    checkpoint_config = optional(object({
      save_checkpoints_max = number
    }))
  })
}

resource "huaweicloud_modelarts_training_job" "test" {
  kind       = "job"
  train_type = "SFT"

  metadata {
    name = var.training_job_name
  }

  algorithm {
    dynamic "inputs" {
      for_each = var.inputs

      content {
        remote {
          dataset {
            id                 = inputs.value.remote.dataset.id
            name               = inputs.value.remote.dataset.name
            dataset_proportion = inputs.value.remote.dataset.dataset_proportion
          }
        }
      }
    }
  }

  spec {
    resource {
      flavor_id  = var.resource_pool_flavor_id
      node_count = var.resource_pool_node_count
    }

    asset_id = var.asset_id

    output_model {
      obs {
        obs_path   = var.output_model.obs.obs_path
        local_path = var.output_model.obs.local_path
      }
    }

    asset_model {
      name    = var.asset_model.name
      version = var.asset_model.version
      type    = var.asset_model.type
      desc    = var.asset_model.desc
      series  = var.asset_model.series
    }
  }

  ftjob_config {
    dynamic "checkpoint_config" {
      for_each = var.ftjob_config.checkpoint_config != null ? [var.ftjob_config.checkpoint_config] : []

      content {
        save_checkpoints_max = checkpoint_config.value.save_checkpoints_max
      }
    }

    task_env {
      dynamic "envs" {
        for_each = var.ftjob_config.task_env.envs

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
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the training job is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `kind` - (Required, String, NonUpdatable) Specifies the type of the training job.  
  The valid values are as follows:
  + **job**
  + **federated_pool_job**
  + **edge_job**
  + **hetero_job**
  + **mrs_job**
  + **autosearch_job**
  + **diag_job**
  + **visualization_job**

* `metadata` - (Required, List) Specifies the metadata configuration of the training job.  
  The [metadata](#modelarts_training_job_metadata) structure is documented below.

* `algorithm` - (Required, List, NonUpdatable) Specifies the algorithm configuration of the training job.  
  The [algorithm](#modelarts_training_job_algorithm) structure is documented below.

* `spec` - (Required, List, NonUpdatable) Specifies the specification configuration of the training job.  
  The [spec](#modelarts_training_job_spec) structure is documented below.

* `endpoints` - (Optional, List, NonUpdatable) Specifies the remote access configuration of the training job.  
  The [endpoints](#modelarts_training_job_endpoints) structure is documented below.

* `train_type` - (Optional, String, NonUpdatable) Specifies the training type of the fine-tuning job.  
  The valid values are as follows:
  + **SFT** - Full fine-tuning
  + **PRETRAIN** - Pre-training
  + **LORA** - LoRA fine-tuning
  + **DPO** - DPO reinforcement learning
  + **RFT** - RFT reinforcement learning

* `ftjob_config` - (Optional, List, NonUpdatable) Specifies the fine-tuning training job configuration.  
  The [ftjob_config](#modelarts_training_job_ftjob_config) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the training job.

<a name="modelarts_training_job_metadata"></a>
The `metadata` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the training job.  
  The maximum length is `64` characters, only letters, digits, underscores (_) and hyphens (-) are allowed.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the ID of the workspace to which the training
  job belongs.  
  If omitted, the default workspace is used.

* `training_experiment_reference` - (Optional, List, NonUpdatable) Specifies the experiment configuration to be
  associated with the training job.  
  The [training_experiment_reference](#modelarts_training_job_metadata_training_experiment_reference)
  structure is documented below.

* `description` - (Optional, String) Specifies the description of the training job.  
  The maximum length is `256` characters.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the advanced feature configuration of the training job.  
  For details, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-modelarts/CreateTrainingJob.html#EN-US_TOPIC_0000002411101796__request_JobMetadata).

  + **job_template**: Heterogeneous job, can be set to **Template RL**.
  + **fault-tolerance/job-retry-num**: Number of retries upon a fault, can be set to a positive integer.
  + **fault-tolerance/job-unconditional-retry**: Unconditional restart, can be set to the string value of **true**.
    or **false**.
  + **fault-tolerance/hang-retry**: Restart upon suspension, can be set to the string value of **true** or **false**.
  + **jupyter-lab/enable**: JupyterLab training application, can be set to the string value of **true** or **false**.
  + **tensorboard/enable**: TensorBoard training application, can be set to the string value of **true** or **false**.
  + **mindstudio-insight/enable**: MindStudio Insight training application, can be set to the string value of **true**.
    or **false**.
  + **fault-tolerance/hccl_op_retry**: Operator retry, can be set to the string value of **true** or **false**.
  + **performance_diagnosis_enabled**: Performance diagnosis, can be set to the string value of **true** or **false**.
  
<a name="modelarts_training_job_metadata_training_experiment_reference"></a>
The `training_experiment_reference` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the ID of the experiment.

<a name="modelarts_training_job_algorithm"></a>
The `algorithm` block supports:

* `id` - (Optional, String, NonUpdatable) Specifies the ID of the algorithm.

* `subscription_id` - (Optional, String, NonUpdatable) Specifies the subscription ID of the subscribed algorithm.  
  This parameter must be used together with `algorithm.item_version_id`.

* `item_version_id` - (Optional, String, NonUpdatable) Specifies the version ID of the subscribed algorithm.  
  This parameter must be used together with `algorithm.subscription_id`.

* `code_dir` - (Optional, String, NonUpdatable) Specifies the code directory of the training job.  
  For custom algorithm, this parameter is required.

* `boot_file` - (Optional, String, NonUpdatable) Specifies the boot file of the training job code.  
  For preset algorithm, this parameter is available and required.

* `autosearch_config_path` - (Optional, String, NonUpdatable) Specifies the YAML configuration path of the auto search
  job.

* `autosearch_framework_path` - (Optional, String, NonUpdatable) Specifies the framework code directory of the auto
  search job.

* `command` - (Optional, String, NonUpdatable) Specifies the container startup command for
  custom image scenarios.  
  For custom image, this parameter is available and required.

* `local_code_dir` - (Optional, String, NonUpdatable) Specifies the local path where the algorithm code is downloaded
  in the training container.  
  + Must be a directory under `/home`.
  + In `v1` compatibility mode, this parameter is invalid.
  + Incompatible with `algorithm.code_dir` starting with `file://`.
  + Not supported to configure `/home/ma-user/modelarts`, `/home/ma-user/modelarts-dev`, `/home/ma-user/infer`
    and their subdirectories, nor to configure `/home/ma-user`.  

  When `subscription_id` parameter is specified, this parameter is not specified.

* `working_dir` - (Optional, String, NonUpdatable) Specifies the working directory when running the algorithm.  
  In `v1` compatibility mode, this parameter is invalid.  
  When `subscription_id` parameter is specified, this parameter is not specified.

* `engine` - (Optional, List, NonUpdatable) Specifies the engine configuration of the training job.  
  The [engine](#modelarts_training_job_algorithm_engine) structure is documented below.  
  When `subscription_id` parameter is specified, this parameter is not specified.
  
* `inputs` - (Optional, List, NonUpdatable) Specifies the data input list of the training job.  
  The [inputs](#modelarts_training_job_algorithm_inputs) structure is documented below.

* `outputs` - (Optional, List, NonUpdatable) Specifies the data output list of the training job.  
  The [outputs](#modelarts_training_job_algorithm_outputs) structure is documented below.

* `parameters` - (Optional, List, NonUpdatable) Specifies the runtime parameter list of the training job.  
  The [parameters](#modelarts_training_job_algorithm_parameters) structure is documented below.

* `environments` - (Optional, Map, NonUpdatable) Specifies the environment variables of the training job.  
  + The maximum length of the key is `8,192` characters.
  + The maximum length of the value is `4,096` characters.
  + The maximum allow setting of environment variables is `100`.
  + The variable name should only contain letters, digits, underscores (_), and start with a letter or underscore.
  + Does not support using the symbol `$` to reference variables.

<a name="modelarts_training_job_algorithm_engine"></a>
The `engine` block supports:

* `engine_id` - (Optional, String, NonUpdatable) Specifies the engine specification ID.
  For preset image, this parameter is available and required.

* `engine_name` - (Optional, String, NonUpdatable) Specifies the engine specification name.
  For preset image and custom image with preset image, this parameter is available and required.

* `engine_version` - (Optional, String, NonUpdatable) Specifies the engine specification version.
  For preset image, this parameter is available and required.

* `image_url` - (Optional, String, NonUpdatable) Specifies the custom image URL obtained from SWR.  
  For custom image, this parameter is available and required.

* `install_sys_packages` - (Optional, Bool, NonUpdatable) Specifies whether to install the moxing version specified by
  the training platform.  
  Defaults to **false**.  
  Only supported when `engine_name`, `engine_version`, and `image_url` are configured.

<a name="modelarts_training_job_algorithm_inputs"></a>
The `inputs` block supports:

* `remote` - (Required, List, NonUpdatable) Specifies the actual input data configuration.  
  The [remote](#modelarts_training_job_algorithm_inputs_remote) structure is documented below.

* `name` - (Optional, String, NonUpdatable) Specifies the input channel name.

* `local_dir` - (Optional, String, NonUpdatable) Specifies the local path mapped by the input channel in the
  container.
  The format is `${local_code_path}/input/${parameter_name}_${index}`.
  e.g. `/home/ma-user/modelarts/inputs/data_url_0`.
  
* `access_method` - (Optional, String, NonUpdatable) Specifies the delivery method of the input channel path.  
  The valid values are as follows:
  + **parameter**
  + **env**

* `description` - (Optional, String, NonUpdatable) Specifies the input channel description.

<a name="modelarts_training_job_algorithm_inputs_remote"></a>
The `remote` block supports:

* `dataset` - (Optional, List, NonUpdatable) Specifies the dataset input configuration.  
  The [dataset](#modelarts_training_job_algorithm_inputs_remote_dataset) structure is documented below.

* `obs` - (Optional, List, NonUpdatable) Specifies the OBS input configuration.  
  The [obs](#modelarts_training_job_algorithm_inputs_remote_obs) structure is documented below.

  -> Exactly one of `dataset` or `obs` parameter must be specified.

<a name="modelarts_training_job_algorithm_inputs_remote_dataset"></a>
The `dataset` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the dataset ID.

* `name` - (Optional, String, NonUpdatable) Specifies the dataset name.

* `version_id` - (Optional, String, NonUpdatable) Specifies the dataset version ID.  
  Required when `service_type` is not set to **V3**.

* `service_type` - (Optional, String, NonUpdatable) Specifies the dataset service type.
  + **V3**: Asset service provided dataset

* `dataset_proportion` - (Optional, Int, NonUpdatable) Specifies the dataset proportion used for fine-tuning training
  jobs.

<a name="modelarts_training_job_algorithm_inputs_remote_obs"></a>
The `obs` block supports:

* `obs_url` - (Required, String, NonUpdatable) Specifies the OBS path URL of the dataset.

<a name="modelarts_training_job_algorithm_outputs"></a>
The `outputs` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the output channel name.

* `remote` - (Required, List, NonUpdatable) Specifies the actual output data configuration.  
  The [remote](#modelarts_training_job_algorithm_outputs_remote) structure is documented below.

* `local_dir` - (Optional, String, NonUpdatable) Specifies the local path mapped by the output channel in the
  container.  
  The format is `${local_code_path}/output/${parameter_name}_${index}`.
  e.g. `/home/ma-user/modelarts/outputs/data_url_0`.

* `access_method` - (Optional, String, NonUpdatable) Specifies the delivery method of the output channel path.  
  The valid values are as follows:
  + **parameter**
  + **env**

* `description` - (Optional, String, NonUpdatable) Specifies the output channel description.

<a name="modelarts_training_job_algorithm_outputs_remote"></a>
The `remote` block supports:

* `obs` - (Required, List, NonUpdatable) Specifies the OBS output configuration.  
  The [obs](#modelarts_training_job_algorithm_outputs_remote_obs) structure is documented below.

<a name="modelarts_training_job_algorithm_outputs_remote_obs"></a>
The `obs` block supports:

* `obs_url` - (Required, String, NonUpdatable) Specifies the OBS path URL where the data is output.

<a name="modelarts_training_job_algorithm_parameters"></a>
The `parameters` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the parameter name.

* `value` - (Optional, String, NonUpdatable) Specifies the parameter value.

* `description` - (Optional, String, NonUpdatable) Specifies the parameter description.

* `constraint` - (Optional, List, NonUpdatable) Specifies the parameter constraint configuration.  
  The [constraint](#modelarts_training_job_algorithm_parameters_constraint) structure is documented below.

<a name="modelarts_training_job_algorithm_parameters_constraint"></a>
The `constraint` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the parameter type.

* `editable` - (Optional, Bool, NonUpdatable) Whether the parameter is editable.  
  Defaults to **false**.

* `required` - (Optional, Bool, NonUpdatable) Whether the parameter is required.  
  Defaults to **false**.

* `valid_type` - (Optional, String, NonUpdatable) Specifies the valid type of the parameter.  
  The valid values are as follows:
  + **Choice**: Enumeration value
  + **Range**
  + **None**

* `valid_range` - (Optional, List, NonUpdatable) Specifies the valid range of the parameter.

<a name="modelarts_training_job_spec"></a>
The `spec` block supports:

* `resource` - (Required, List, NonUpdatable) Specifies the resource specification of the training job.  
  The [resource](#modelarts_training_job_spec_resource) structure is documented below.

* `runtime_type` - (Optional, String, NonUpdatable) Specifies the runtime type of the training job.

* `log_export_path` - (Optional, List, NonUpdatable) Specifies the log export path configuration of the
  training job.  
  The [log_export_path](#modelarts_training_job_spec_log_export_path) structure is documented below.

* `log_export_config` - (Optional, List, NonUpdatable) Specifies the log export configuration of the training job.  
  The [log_export_config](#modelarts_training_job_spec_log_export_config) structure is documented below.

* `auto_stop` - (Optional, List, NonUpdatable) Specifies the auto stop configuration of the training job.  
  The [auto_stop](#modelarts_training_job_spec_auto_stop) structure is documented below.

* `schedule_policy` - (Optional, List, NonUpdatable) Specifies the schedule policy configuration of the
  training job.  
  The [schedule_policy](#modelarts_training_job_spec_schedule_policy) structure is documented below.  
  This parameter is available when only `spec.resource.pool_id` is specified.

* `notification` - (Optional, List, NonUpdatable) Specifies the notification configuration of the training job.  
  The [notification](#modelarts_training_job_spec_notification) structure is documented below.

* `custom_metrics` - (Optional, List, NonUpdatable) Specifies the custom metrics collection configuration of the
  training job.  
  The [custom_metrics](#modelarts_training_job_spec_custom_metrics) structure is documented below.

* `output_model` - (Optional, List, NonUpdatable) Specifies the output model configuration of the training job.  
  The [output_model](#modelarts_training_job_spec_output_model) structure is documented below.

* `asset_model` - (Optional, List, NonUpdatable) Specifies the asset model configuration of the training job.  
  The [asset_model](#modelarts_training_job_spec_asset_model) structure is documented below.

* `asset_id` - (Optional, String, NonUpdatable) Specifies the asset model ID for fine-tuning training jobs.

* `volumes` - (Optional, List, NonUpdatable) Specifies the volume mount configuration of the training job.  
  The [volumes](#modelarts_training_job_spec_volumes) structure is documented below.  
  This parameter is available when only `spec.resource.pool_id` is specified.

<a name="modelarts_training_job_spec_resource"></a>
The `resource` block supports:

* `node_count` - (Required, Int, NonUpdatable) Specifies the number of resource replicas used by the training job.

* `flavor_id` - (Optional, String, NonUpdatable) Specifies the resource flavor ID of the training job.  
  For public resource pool, this parameter is available and required.

* `pool_id` - (Optional, String, NonUpdatable) Specifies the dedicated resource pool ID.  
  For dedicated resource pool, this parameter is available and required.

* `pool_group_id` - (Optional, String, NonUpdatable) Specifies the federated resource pool ID.  
  This parameter is required only when `kind` is set to **federated_pool_job**.

* `main_container_customized_flavor` - (Optional, List, NonUpdatable) Specifies the customized flavor configuration of
  the main container.  
  The [main_container_customized_flavor](#modelarts_training_job_spec_resource_main_container_customized_flavor)
  structure is documented below.  
  This parameter is available only when `pool_id` is specified.

<a name="modelarts_training_job_spec_resource_main_container_customized_flavor"></a>
The `main_container_customized_flavor` block supports:

* `cpu_core_num` - (Optional, Float, NonUpdatable) Specifies the number of CPU cores.

* `mem_size` - (Optional, Float, NonUpdatable) Specifies the memory size.

* `accelerator_num` - (Optional, Float, NonUpdatable) Specifies the number of accelerator cards.

<a name="modelarts_training_job_spec_log_export_path"></a>
The `log_export_path` block supports:

* `obs_url` - (Optional, String, NonUpdatable) Specifies the OBS path where training job logs are saved.

* `host_path` - (Optional, String, NonUpdatable) Specifies the host path where training job logs are saved.

  -> At least one of `obs_url` or `host_path` parameter must be specified.

<a name="modelarts_training_job_spec_log_export_config"></a>
The `log_export_config` block supports:

* `version` - (Optional, String, NonUpdatable) Specifies the log version.  
  The valid values are as follows:
  + **v0**
  + **v1**

* `rotation_enabled` - (Optional, Bool, NonUpdatable) Whether to enable log rotation download.  
  Defaults to **false**.

<a name="modelarts_training_job_spec_auto_stop"></a>
The `auto_stop` block supports:

* `time_unit` - (Required, String, NonUpdatable) Specifies the time unit of the auto stop duration.  
  The valid values are as follows:
  + **HOURS**

* `duration` - (Required, Int, NonUpdatable) Specifies the running duration before the training job is automatically
  stopped.

<a name="modelarts_training_job_spec_schedule_policy"></a>
The `schedule_policy` block supports:

* `priority` - (Optional, Int, NonUpdatable) Specifies the priority of the training job.  
  This parameter is available only when `spec.resource.pool_id` is specified.

* `preemptible` - (Optional, Bool, NonUpdatable) Whether the training job can be preempted.  
  This parameter is available only when `spec.resource.pool_id` is specified.

* `required_affinity` - (Optional, List, NonUpdatable) Specifies the required affinity configuration of the training
  job.  
  The [required_affinity](#modelarts_training_job_spec_schedule_policy_required_affinity) structure is documented below.

* `preferred_affinity` - (Optional, List, NonUpdatable) Specifies the preferred affinity configuration of the training
  job.  
  The [preferred_affinity](#modelarts_training_job_spec_schedule_policy_preferred_affinity) structure is documented
  below.

  -> Only one of the `required_affinity` and `preferred_affinity` parameters can be specified.

<a name="modelarts_training_job_spec_schedule_policy_required_affinity"></a>
The `required_affinity` block supports:

* `affinity_type` - (Optional, String, NonUpdatable) Specifies the affinity scheduling policy type.  
  The valid values are as follows:
  + **cabinet**
  + **hyperinstance**

* `affinity_group_size` - (Optional, Int, NonUpdatable) Specifies the affinity group size.  
  This parameter is required when `affinity_type` is set to **hyperinstance**.

* `node_affinity` - (Optional, List, NonUpdatable) Specifies the node affinity configuration of the training job.  
  The [node_affinity](#modelarts_training_job_spec_schedule_policy_required_affinity_node_affinity)
  structure is documented below.

<a name="modelarts_training_job_spec_schedule_policy_required_affinity_node_affinity"></a>
The `node_affinity` block supports:

* `node_selector_terms` - (Required, List, NonUpdatable) Specifies the required node selector terms.  
  The [node_selector_terms](#modelarts_training_job_spec_schedule_policy_node_selector_terms)
  structure is documented below.

<a name="modelarts_training_job_spec_schedule_policy_preferred_affinity"></a>
The `preferred_affinity` block supports:

* `node_affinity` - (Optional, List, NonUpdatable) Specifies the preferred node affinity terms.  
  The [node_affinity](#modelarts_training_job_spec_schedule_policy_preferred_affinity_node_affinity)
  structure is documented below.

<a name="modelarts_training_job_spec_schedule_policy_preferred_affinity_node_affinity"></a>
The `node_affinity` block supports:

* `weight` - (Optional, Int, NonUpdatable) Specifies the weight associated with the preferred node selector term.  
  The range is from `0` to `100`.

* `preference` - (Optional, List, NonUpdatable) Specifies the preferred node selector term.  
  The [preference](#modelarts_training_job_spec_schedule_policy_node_selector_terms) structure is documented below.

<a name="modelarts_training_job_spec_schedule_policy_node_selector_terms"></a>
The `node_selector_terms` and `preference` block supports:

* `match_expressions` - (Optional, List, NonUpdatable) Specifies the node selector requirements based on node labels.  
  The [match_expressions](#modelarts_training_job_spec_schedule_policy_node_selector_requirement)
  structure is documented below.

* `match_fields` - (Optional, List, NonUpdatable) Specifies the node selector requirements based on node fields.  
  The [match_fields](#modelarts_training_job_spec_schedule_policy_node_selector_requirement)
  structure is documented below.

  -> Exactly one of `match_expressions` or `match_fields` parameter must be specified.

<a name="modelarts_training_job_spec_schedule_policy_node_selector_requirement"></a>
The `match_expressions` and `match_fields` block supports:

* `key` - (Required, String, NonUpdatable) Specifies the label key used by the node selector requirement.

* `operator` - (Required, String, NonUpdatable) Specifies the operator used by the node selector requirement.
  The valid values are as follows:
  + **In**: The value of the key must be in the given value list.
  + **NotIn**: The value of the key must not be in the given value list.
  + **Exists**: The key must exist, but its value is not specific.
  + **DoesNotExist**: The key must not exist.
  + **Gt**: The value of the key must be greater than the given value.
  + **Lt**: The value of the key must be less than the given value.

* `values` - (Optional, List, NonUpdatable) Specifies the label values used by the node selector requirement.
  + If the `operator` is **In** or **NotIn**, the value array cannot be empty.
  + If the `operator` is **Exists** or **DoesNotExist**, the value array must be empty.
  + If the `operator` is **Gt** or **Lt**, the value array must contain one element, which will be interpreted
    as an integer.

<a name="modelarts_training_job_spec_notification"></a>
The `notification` block supports:

* `topic_urn` - (Required, String, NonUpdatable) Specifies the URN of the SMN topic for training event notifications.

* `events` - (Optional, List, NonUpdatable) Specifies the training events that trigger notifications.  
  The valid values are as follows:
  + **JobStarted**
  + **JobCompleted**
  + **JobFailed**
  + **JobTerminated**
  + **JobRestarted**
  + **JobHanged**
  + **JobPreempted**

<a name="modelarts_training_job_spec_custom_metrics"></a>
The `custom_metrics` block supports:

* `exec` - (Optional, List, NonUpdatable) Specifies the command-based metrics collection configuration.  
  The [exec](#modelarts_training_job_spec_custom_metrics_exec) structure is documented below.

* `http_get` - (Optional, List, NonUpdatable) Specifies the HTTP-based metrics collection configuration.  
  The [http_get](#modelarts_training_job_spec_custom_metrics_http_get) structure is documented below.

  -> Exactly one of `exec` or `http_get` parameter must be specified.

<a name="modelarts_training_job_spec_custom_metrics_exec"></a>
The `exec` block supports:

* `command` - (Required, List, NonUpdatable) Specifies the command used to collect metrics.

<a name="modelarts_training_job_spec_custom_metrics_http_get"></a>
The `http_get` block supports:

* `path` - (Required, String, NonUpdatable) Specifies the HTTP path used to collect metrics.

* `port` - (Required, Int, NonUpdatable) Specifies the HTTP port used to collect metrics.

<a name="modelarts_training_job_spec_output_model"></a>
The `output_model` block supports:

* `obs` - (Required, List, NonUpdatable) Specifies the OBS output configuration of the model.  
  The [obs](#modelarts_training_job_spec_output_model_obs) structure is documented below.

<a name="modelarts_training_job_spec_output_model_obs"></a>
The `obs` block supports:

* `obs_path` - (Required, String, NonUpdatable) Specifies the OBS path where the output model is saved.

* `local_path` - (Optional, String, NonUpdatable) Specifies the local path where the output model is saved.

<a name="modelarts_training_job_spec_asset_model"></a>
The `asset_model` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the model name.

* `version` - (Required, String, NonUpdatable) Specifies the model version.

* `type` - (Required, String, NonUpdatable) Specifies the model type.

* `code` - (Optional, String, NonUpdatable) Specifies the model code.

* `desc` - (Optional, String, NonUpdatable) Specifies the model description.

* `series` - (Optional, String, NonUpdatable) Specifies the model series.

<a name="modelarts_training_job_spec_volumes"></a>
The `volumes` block supports:

* `nfs` - (Optional, List, NonUpdatable) Specifies the NFS volume mount configuration.  
  The [nfs](#modelarts_training_job_spec_volumes_nfs) structure is documented below.

* `pfs` - (Optional, List, NonUpdatable) Specifies the PFS volume mount configuration.  
  The [pfs](#modelarts_training_job_spec_volumes_pfs) structure is documented below.

* `obs` - (Optional, List, NonUpdatable) Specifies the OBS volume mount configuration.  
  The [obs](#modelarts_training_job_spec_volumes_obs) structure is documented below.

<a name="modelarts_training_job_spec_volumes_nfs"></a>
The `nfs` block supports:

* `nfs_server_path` - (Optional, String, NonUpdatable) Specifies the NFS server path.

* `local_path` - (Optional, String, NonUpdatable) Specifies the local mount path in the training container.

* `read_only` - (Optional, Bool, NonUpdatable) Whether the NFS volume is read-only in the container.  
  Defaults to **false**.

<a name="modelarts_training_job_spec_volumes_pfs"></a>
The `pfs` block supports:

* `pfs_path` - (Optional, String, NonUpdatable) Specifies the OBSFS path.

* `local_path` - (Optional, String, NonUpdatable) Specifies the local mount path in the training container.
  
  At least one of `pfs_path` and `local_path` must be specified.

<a name="modelarts_training_job_spec_volumes_obs"></a>
The `obs` block supports:

* `obs_path` - (Optional, String, NonUpdatable) Specifies the OBS path to be mounted.

* `local_path` - (Optional, String, NonUpdatable) Specifies the local mount path in the training container.

  At least one of `obs_path` and `local_path` must be specified.

<a name="modelarts_training_job_endpoints"></a>
The `endpoints` block supports:

* `ssh` - (Required, List, NonUpdatable) Specifies the SSH connection configuration.  
  The [ssh](#modelarts_training_job_endpoints_ssh) structure is documented below.

<a name="modelarts_training_job_endpoints_ssh"></a>
The `ssh` block supports:

* `key_pair_names` - (Required, List, NonUpdatable) Specifies the SSH key pair names.  
  The key pairs can be created and viewed on the ECS console.

<a name="modelarts_training_job_ftjob_config"></a>
The `ftjob_config` block supports:

* `ft_job_uuid` - (Optional, String, NonUpdatable) Specifies the model ID.

* `ft_train_type` - (Optional, String, NonUpdatable) Specifies the model training type.

* `model_type` - (Optional, String, NonUpdatable) Specifies the model type.

* `train_output_path` - (Optional, String, NonUpdatable) Specifies the output path of the training job.

* `train_process` - (Optional, Float, NonUpdatable) Specifies the training process progress.

* `checkpoint_id` - (Optional, String, NonUpdatable) Specifies the checkpoint ID.

* `task_env` - (Optional, List, NonUpdatable) Specifies the fine-tuning training job environment parameters.  
  The [task_env](#modelarts_training_job_ftjob_config_task_env) structure is documented below.

* `checkpoint_config` - (Optional, List, NonUpdatable) Specifies the checkpoint configuration.  
  The [checkpoint_config](#modelarts_training_job_ftjob_config_checkpoint_config) structure is documented below.

<a name="modelarts_training_job_ftjob_config_task_env"></a>
The `task_env` block supports:

* `envs` - (Required, List, NonUpdatable) Specifies the fine-tuning training environment variables.  
  The [envs](#modelarts_training_job_ftjob_config_task_env_envs) structure is documented below.

<a name="modelarts_training_job_ftjob_config_task_env_envs"></a>
The `envs` block supports:

* `label` - (Optional, String, NonUpdatable) Specifies the label of the environment variable.

* `des` - (Optional, String, NonUpdatable) Specifies the description of the environment variable.

* `env_name` - (Optional, String, NonUpdatable) Specifies the name of the environment variable.

* `env_type` - (Optional, String, NonUpdatable) Specifies the type of the environment variable.

* `value` - (Optional, String, NonUpdatable) Specifies the value of the environment variable.

* `modifiable` - (Optional, Bool, NonUpdatable) Whether the environment variable is modifiable.  
  Defaults to **false**.

* `displayable` - (Optional, Bool, NonUpdatable) Whether the environment variable is displayable.  
  Defaults to **false**.

* `used_steps` - (Optional, List, NonUpdatable) Specifies the steps where the environment variable is used.

<a name="modelarts_training_job_ftjob_config_checkpoint_config"></a>
The `checkpoint_config` block supports:

* `checkpoint_id` - (Optional, String, NonUpdatable) Specifies the checkpoint ID.

* `save_checkpoints_max` - (Optional, Int, NonUpdatable) Specifies the maximum number of checkpoints to save.  
  `0` means disabled, `-1` means unlimited.

* `skipped_steps` - (Optional, Int, NonUpdatable) Specifies the number of steps to skip.  
  `0` means no skip.

* `restore_training` - (Optional, Int, NonUpdatable) Specifies whether to restore training from a checkpoint.  
  `0` means no restore, `1` means restore.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the training job, in RFC3339 format.

* `status` - The current status of the training job.  
  + **Pending**
  + **Running**
  + **Completed**
  + **Failed**
  + **Terminated**
  + **Abnormal**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The training job can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_training_job.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `metadata.0.annotations`, `algorithm.0.inputs`, `spec.0.output_model`,
`spec.0.asset_model`, `spec.0.asset_id`, `train_type` and `ftjob_config`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_training_job" "test" {
  ...

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
      algorithm.0.inputs,
      spec.0.output_model,
      spec.0.asset_model,
      spec.0.asset_id,
      train_type,
      ftjob_config,
    ]
  }
}
```
