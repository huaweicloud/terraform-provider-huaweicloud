---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_algorithm"
description: |-
  Manages a ModelArts algorithm resource within HuaweiCloud.
---

# huaweicloud_modelarts_algorithm

Manages a ModelArts algorithm resource within HuaweiCloud.

## Example Usage

### Create a preset image algorithm

```hcl
variable "algorithm_name" {}
variable "code_directory" {}
variable "boot_file" {}
variable "engine_id" {}
variable "engine_name" {}
variable "engine_version" {}
variable "inputs" {
  type = list(object({
    name          = string
    description   = optional(string)
    access_method = optional(string, "parameter")

    remote_constraints = optional(list(object({
      data_type  = string
      attributes = optional(string)
    })), [])
  }))
}
variable "outputs" {
  type = list(object({
    name          = string
    access_method = optional(string, "parameter")
    description   = optional(string)
  }))
}
variable "parameters" {
  type = list(object({
    name  = string
    value = optional(string)

    constraint = optional(object({
      type        = string
      required    = optional(bool)
      editable    = optional(bool)
      valid_range = optional(list(string))
      valid_type  = optional(string, "None")
    }), null)
  }))
}

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name = var.algorithm_name
  }

  job_config {
    code_dir  = var.code_directory
    boot_file = var.boot_file

    engine {
      engine_id      = var.engine_id
      engine_name    = var.engine_name
      engine_version = var.engine_version
    }

    parameters_customization = true

    dynamic "parameters" {
      for_each = var.parameters

      content {
        name  = parameters.value.name
        value = parameters.value.value

        dynamic "constraint" {
          for_each = parameters.value.constraint != null ? [1] : []

          content {
            type        = parameters.value.constraint.type
            required    = parameters.value.constraint.required
            editable    = parameters.value.constraint.editable
            valid_range = parameters.value.constraint.valid_range
            valid_type  = parameters.value.constraint.valid_type
          }
        }
      }
    }

    dynamic "inputs" {
      for_each = var.inputs

      content {
        name          = inputs.value.name
        description   = inputs.value.description
        access_method = inputs.value.access_method

        dynamic "remote_constraints" {
          for_each = inputs.value.remote_constraints

          content {
            data_type  = remote_constraints.value.data_type
            attributes = remote_constraints.value.attributes
          }
        }
      }
    }

    dynamic "outputs" {
      for_each = var.outputs

      content {
        name          = outputs.value.name
        description   = outputs.value.description
        access_method = outputs.value.access_method
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the algorithm is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `metadata` - (Required, List) Specifies the metadata of the algorithm.  
  The [metadata](#modelarts_algorithm_metadata) structure is documented below.

* `job_config` - (Required, List) Specifies the configuration of the algorithm.  
  The [job_config](#modelarts_algorithm_job_config) structure is documented below.

* `resource_requirements` - (Optional, List) Specifies the resource constraint list of the algorithm.  
  The [resource_requirements](#modelarts_algorithm_resource_requirements) structure is documented below.

* `advanced_config` - (Optional, List) Specifies the advanced configuration of the algorithm.  
  The [advanced_config](#modelarts_algorithm_advanced_config) structure is documented below.

<a name="modelarts_algorithm_metadata"></a>
The `metadata` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the algorithm.  
  The maximum length is `64`. Only letters, digits, underscores (_) and hyphens (-) are allowed.

* `description` - (Optional, String) Specifies the description of the algorithm.  
  The maximum length is `256`.

* `workspace_id` - (Optional, String) Specifies the ID of the workspace to which the algorithm belongs.  
  If omitted, the default workspace (`0`) is used.

* `tags` - (Optional, List) Specifies the tags of the algorithm.  
  The [tags](#modelarts_algorithm_metadata_tags) structure is documented below.

<a name="modelarts_algorithm_metadata_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.  
  The valid values are as follows:
  + **auto_search**
  + **turbo**
  + **economic**
  + **federated**

<a name="modelarts_algorithm_job_config"></a>
The `job_config` block supports:

* `engine` - (Required, List) Specifies the engine configuration of the algorithm.  
  The [engine](#modelarts_algorithm_job_config_engine) structure is documented below.

* `code_dir` - (Optional, String) Specifies the code directory of the algorithm.  
  For preset image algorithm, this parameter is **required**.

* `boot_file` - (Optional, String) Specifies the boot file path under the code directory.  
  For preset image algorithm, this parameter is available and **required**.

* `command` - (Optional, String) Specifies the container start command for custom image algorithm.  
  For custom image algorithm, this parameter is available and **required**.

* `inputs` - (Optional, List) Specifies the data input list of the algorithm.  
  The [inputs](#modelarts_algorithm_job_config_inputs) structure is documented below.

* `outputs` - (Optional, List) Specifies the data output list of the algorithm.  
  The [outputs](#modelarts_algorithm_job_config_outputs) structure is documented below.

* `parameters_customization` - (Optional, Bool) Specifies whether the hyperparameter can be customized when creating
  a training job.  
  Default to `false`.

* `parameters` - (Optional, List) Specifies the list of running parameters of the algorithm.  
  The [parameters](#modelarts_algorithm_job_config_parameters) structure is documented below.

<a name="modelarts_algorithm_job_config_engine"></a>
The `engine` block supports:

* `engine_id` - (Optional, String) Specifies the ID of the engine specification.  
  For preset image algorithm, this parameter is available and **required**.

* `engine_name` - (Optional, String) Specifies the name of the engine.  
  For preset image algorithm, this parameter is available and **required**.

* `engine_version` - (Optional, String) Specifies the version of the engine version.  
  For preset image algorithm, this parameter is available and **required**.

* `image_url` - (Optional, String) Specifies the custom image URL of the algorithm.  
  For custom image algorithm, this parameter is available and **required**.

<a name="modelarts_algorithm_job_config_inputs"></a>
The `inputs` block supports:

* `name` - (Optional, String) Specifies the name of the data input channel.

* `description` - (Optional, String) Specifies the description of the data input channel.

* `access_method` - (Optional, String) Specifies the access method of the data input channel.  
  The valid values are as follows:
  + **parameter**
  + **env**

* `remote_constraints` - (Optional, List) Specifies the constraint of the data input.  
  The [remote_constraints](#modelarts_algorithm_job_config_input_remote_constraints) structure is documented below.

<a name="modelarts_algorithm_job_config_input_remote_constraints"></a>
The `remote_constraints` block supports:

* `data_type` - (Optional, String) Specifies the type of the data input.  
  The valid values are as follows:
  + **obs**
  + **modelarts_dataset**

* `attributes` - (Optional, String) Specifies the attributes of the input data, in JSON format.  
  When `data_type` is **modelarts_dataset**, this parameter is available.
  
  For **data_format** key, the valid values are as follows:
  + **Default**: Indicates the manifest format.
  + **CarbonData**

  For **data_segmentation** key, the valid values are as follows:
  + **true**: Indicates that the data is segmented.
  + **false**: Indicates that the data is not segmented.
  + **no_limit**

  For **dataset_type** key, the valid values are as follows:
  + **0**: Image classification.
  + **1**: Object detection.
  + **2**: Image labeling.
  + **3**: Image segmentation.
  + **100**: Text classification.
  + **101**: Text labeling.
  + **102**: Text triplet.
  + **200**: Sound classification.
  + **201**: Speech labeling.
  + **202**: Speech paragraph labeling.
  + **300**: Point cloud.
  + **400**: Table.
  + **600**: Video labeling.
  + **900**: Free format.

<a name="modelarts_algorithm_job_config_outputs"></a>
The `outputs` block supports:

* `name` - (Required, String) Specifies the name of the data output channel.

* `access_method` - (Optional, String) Specifies the access method of the data output channel.  
  The valid values are as follows:
  + **parameter**
  + **env**

* `description` - (Optional, String) Specifies the description of the data output channel.

<a name="modelarts_algorithm_job_config_parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the name of the parameter.

* `constraint` - (Required, List) Specifies the constraint of the parameter.  
  The [constraint](#modelarts_algorithm_job_config_parameter_constraint) structure is documented below.

* `value` - (Optional, String) Specifies the value of the parameter.  
  When `constraint.editable` is **true**, this parameter is **required**.

* `description` - (Optional, String) Specifies the description of the parameter.

* `i18n_description` - (Optional, List) Specifies the internationalized description of the parameter.  
  The [i18n_description](#modelarts_algorithm_job_config_parameter_i18n_description) structure is documented below.

<a name="modelarts_algorithm_job_config_parameter_constraint"></a>
The `constraint` block supports:

* `type` - (Required, String) Specifies the type of the parameter.  
  The valid values are as follows:
  + **Integer**
  + **Float**
  + **String**
  + **Boolean**

* `required` - (Optional, Bool) Specifies whether the parameter is required.  
  Default to `false`.

* `editable` - (Optional, Bool) Specifies whether the parameter is editable.  
  Default to `false`.

* `valid_type` - (Optional, String) Specifies the valid type of the parameter value.  
  The valid values are as follows:
  + **Choice** - Enumeration value.
  + **Range** - Range value.
  + **None** - No restriction.

* `valid_range` - (Optional, List) Specifies the valid range list of the parameter value.

<a name="modelarts_algorithm_job_config_parameter_i18n_description"></a>
The `i18n_description` block supports:

* `language` - (Optional, String) Specifies the language code.

* `description` - (Optional, String) Specifies the description in the specified language.

<a name="modelarts_algorithm_resource_requirements"></a>
The `resource_requirements` block supports:

* `key` - (Optional, String) Specifies the key of the resource constraint.  
  The valid values are as follows:
  + **flavor_type** - Resource type.
  + **device_distributed_mode** - Whether multi-card training is supported.
  + **host_distributed_mode** - Whether distributed training is supported.

* `values` - (Optional, List) Specifies the list of values corresponding to the key.  
  For **flavor_type** key, the valid values are as follows:
    + **CPU**
    + **GPU**
    + **Ascend**

  For **device_distributed_mode** key, the valid values are as follows:
    + **multiple**: Supports multi-card training.
    + **singular**: Does not support multi-card training.

  For **host_distributed_mode** key, the valid values are as follows:
    + **multiple**: Supports distributed training.
    + **singular**: Does not support distributed training.

* `operator` - (Optional, String) Specifies the relationship between the key and values.  
  Currently, only **in** is supported.

<a name="modelarts_algorithm_advanced_config"></a>
The `advanced_config` block supports:

* `auto_search` - (Optional, List) Specifies the strategy configuration of hyperparameter search.  
  The [auto_search](#modelarts_algorithm_advanced_config_auto_search) structure is documented below.

<a name="modelarts_algorithm_advanced_config_auto_search"></a>
The `auto_search` block supports:

* `reward_attrs` - (Required, List) Specifies the metric list of search.  
  The [reward_attrs](#modelarts_algorithm_advanced_config_auto_search_reward_attrs) structure is documented below.

* `search_params` - (Required, List) Specifies the parameter list of search.  
  The [search_params](#modelarts_algorithm_advanced_config_auto_search_search_params) structure is documented below.

* `algo_configs` - (Required, List) Specifies the algorithm configuration of search.  
  The [algo_configs](#modelarts_algorithm_advanced_config_auto_search_algo_configs) structure is documented below.

* `skip_search_params` - (Optional, String) Specifies the hyperparameter combination to be excluded from search.

<a name="modelarts_algorithm_advanced_config_auto_search_reward_attrs"></a>
The `reward_attrs` block supports:

* `name` - (Optional, String) Specifies the metric name.

* `mode` - (Optional, String) Specifies the search direction.  
  The valid values are as follows:
  + **max** - A larger metric value is better.
  + **min** - A smaller metric value is better.

* `regex` - (Optional, String) Specifies the metric regular expression.

<a name="modelarts_algorithm_advanced_config_auto_search_search_params"></a>
The `search_params` block supports:

* `name` - (Optional, String) Specifies the hyperparameter name.

* `param_type` - (Optional, String) Specifies the parameter type.  
  The valid values are as follows:
  + **continuous** - Continuous hyperparameter.
  + **discrete** - Discrete hyperparameter.

* `lower_bound` - (Optional, String) Specifies the lower bound of the hyperparameter.

* `upper_bound` - (Optional, String) Specifies the upper bound of the hyperparameter.

* `discrete_points_num` - (Optional, String) Specifies the number of discrete samples for a continuous hyperparameter.

* `discrete_values` - (Optional, List) Specifies the list of discrete values for a discrete hyperparameter.

<a name="modelarts_algorithm_advanced_config_auto_search_algo_configs"></a>
The `algo_configs` block supports:

* `name` - (Optional, String) Specifies the name of the search algorithm.

* `params` - (Optional, List) Specifies the parameter list of the search algorithm.  
  The [params](#modelarts_algorithm_advanced_config_auto_search_algo_config_params) structure is documented below.

<a name="modelarts_algorithm_advanced_config_auto_search_algo_config_params"></a>
The `params` block supports:

* `key` - (Optional, String) Specifies the key of the parameter.

* `value` - (Optional, String) Specifies the value of the parameter.

* `type` - (Optional, String) Specifies the type of the parameter.  
  The valid values are as follows:
  + **Float**
  + **Integer**
  + **String**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Algorithms can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_algorithm.test <id>
```
