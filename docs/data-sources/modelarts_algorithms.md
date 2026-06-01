---
subcategory: "ModelArts"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_algorithms"
description: |-
  Use this data source to query the ModelArts algorithm list within HuaweiCloud.
---

# huaweicloud_modelarts_algorithms

Use this data source to query the ModelArts algorithm list within HuaweiCloud.

## Example Usage

### Query all algorithms

```hcl
data "huaweicloud_modelarts_algorithms" "test" {}
```

### Filter by name fuzzy matching

```hcl
variable "algorithm_name" {}

data "huaweicloud_modelarts_algorithms" "test" {
  searches = "name:${var.algorithm_name}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the algorithms are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Optional, String) Specifies the ID of the workspace to which the algorithms belong.  
  If omitted, all algorithms in the current region will be queried.

* `searches` - (Optional, String) Specifies the filter condition for searching algorithms.
  The format is: `key:value`.
  For example, `name:${name}` means fuzzy matching by algorithm name.

* `sort_by` - (Optional, String) Specifies the field by which the algorithms are sorted.  
  The valid values are as follows:
  + **name**
  + **create_time**

  Defaults to `create_time`.

* `order` - (Optional, String) Specifies the sort order of the algorithms.  
  The valid values are as follows:
  + **asc**
  + **desc**

  Defaults to `desc`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `algorithms` - The list of algorithms that matched filter parameters.  
  The [algorithms](#modelarts_algorithms) structure is documented below.

<a name="modelarts_algorithms"></a>
The `algorithms` block supports:

* `metadata` - The metadata of the algorithm.  
  The [metadata](#modelarts_algorithms_metadata) structure is documented below.

* `job_config` - The configuration of the algorithm.  
  The [job_config](#modelarts_algorithms_job_config) structure is documented below.

* `resource_requirements` - The resource constraint list of the algorithm.  
  The [resource_requirements](#modelarts_algorithms_resource_requirements) structure is documented below.

* `advanced_config` - The advanced configuration of the algorithm.  
  The [advanced_config](#modelarts_algorithms_advanced_config) structure is documented below.

<a name="modelarts_algorithms_metadata"></a>
The `metadata` block supports:

* `id` - The ID of the algorithm.

* `name` - The name of the algorithm.

* `description` - The description of the algorithm.

* `workspace_id` - The ID of the workspace to which the algorithm belongs.

* `user_name` - The user name of the algorithm.

* `source` - The source type of the algorithm.

* `is_valid` - The availability of the algorithm.

* `state` - The state of the algorithm.

* `tags` - The tags of the algorithm.  
  The [tags](#modelarts_algorithms_metadata_tags) structure is documented below.

* `attr_list` - The attribute list of the algorithm.

* `version_num` - The version number of the algorithm.

* `size` - The size of the algorithm.

* `create_time` - The creation time of the algorithm, in RFC3339 format.

* `update_time` - The update time of the algorithm, in RFC3339 format.

<a name="modelarts_algorithms_metadata_tags"></a>
The `tags` block supports:

* `key` - The key of the tag.

<a name="modelarts_algorithms_job_config"></a>
The `job_config` block supports:

* `code_dir` - The code directory of the algorithm.

* `boot_file` - The boot file of the algorithm.

* `command` - The container start command for custom image algorithm.

* `parameters_customization` - Whether the hyperparameter can be customized.

* `parameters` - The running parameters of the algorithm.  
  The [parameters](#modelarts_algorithms_job_config_parameters) structure is documented below.

* `inputs` - The data input list of the algorithm.  
  The [inputs](#modelarts_algorithms_job_config_inputs) structure is documented below.

* `outputs` - The data output list of the algorithm.  
  The [outputs](#modelarts_algorithms_job_config_outputs) structure is documented below.

* `engine` - The engine of the algorithm.  
  The [engine](#modelarts_algorithms_job_config_engine) structure is documented below.

<a name="modelarts_algorithms_job_config_parameters"></a>
The `parameters` block supports:

* `name` - The name of the parameter.

* `value` - The value of the parameter.

* `description` - The description of the parameter.

<a name="modelarts_algorithms_job_config_inputs"></a>
The `inputs` block supports:

* `name` - The name of the data input channel.

* `description` - The description of the data input channel.

<a name="modelarts_algorithms_job_config_outputs"></a>
The `outputs` block supports:

* `name` - The name of the data output channel.

* `description` - The description of the data output channel.

<a name="modelarts_algorithms_job_config_engine"></a>
The `engine` block supports:

* `engine_id` - The ID of the engine.

* `engine_name` - The name of the engine.

* `engine_version` - The version of the engine.

* `image_url` - The custom image URL of the algorithm.

<a name="modelarts_algorithms_resource_requirements"></a>
The `resource_requirements` block supports:

* `key` - The key of the resource constraint.

* `values` - The values corresponding to the key.

* `operator` - The relationship between the key and values.

<a name="modelarts_algorithms_advanced_config"></a>
The `advanced_config` block supports:

* `auto_search` - The hyperparameter search strategy.  
  The [auto_search](#modelarts_algorithms_advanced_config_auto_search) structure is documented below.

<a name="modelarts_algorithms_advanced_config_auto_search"></a>
The `auto_search` block supports:

* `skip_search_params` - The hyperparameter combinations to be excluded from search.

* `reward_attrs` - The metric list of search.  
  The [reward_attrs](#modelarts_algorithms_advanced_config_auto_search_reward_attrs) structure is documented below.

* `search_params` - The search parameters.  
  The [search_params](#modelarts_algorithms_advanced_config_auto_search_search_params) structure is documented below.

* `algo_configs` - The search algorithm configurations.  
  The [algo_configs](#modelarts_algorithms_advanced_config_auto_search_algo_configs) structure is documented below.

<a name="modelarts_algorithms_advanced_config_auto_search_reward_attrs"></a>
The `reward_attrs` block supports:

* `name` - The metric name.

* `mode` - The search direction.

* `regex` - The metric regular expression.

<a name="modelarts_algorithms_advanced_config_auto_search_search_params"></a>
The `search_params` block supports:

* `name` - The hyperparameter name.

* `param_type` - The parameter type.

* `lower_bound` - The lower bound of the hyperparameter.

* `upper_bound` - The upper bound of the hyperparameter.

* `discrete_points_num` - The number of discrete samples for continuous hyperparameters.

* `discrete_values` - The discrete values for discrete hyperparameters.

<a name="modelarts_algorithms_advanced_config_auto_search_algo_configs"></a>
The `algo_configs` block supports:

* `name` - The search algorithm name.

* `params` - The search algorithm parameters.  
  The [params](#modelarts_algorithms_advanced_config_auto_search_algo_configs_params) structure is documented below.

<a name="modelarts_algorithms_advanced_config_auto_search_algo_configs_params"></a>
The `params` block supports:

* `key` - The key of the parameter.

* `value` - The value of the parameter.

* `type` - The type of the parameter.
