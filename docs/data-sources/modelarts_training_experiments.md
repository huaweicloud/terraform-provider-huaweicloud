---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_experiments"
description: |-
  Use this data source to query the ModelArts training experiment list within HuaweiCloud.
---

# huaweicloud_modelarts_training_experiments

Use this data source to query the ModelArts training experiment list under a specified region within HuaweiCloud.

## Example Usage

### Query all training experiments

```hcl
data "huaweicloud_modelarts_training_experiments" "test" {}
```

### Filter by workspace ID

```hcl
variable "workspace_id" {}

data "huaweicloud_modelarts_training_experiments" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the training experiments are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Optional, String) Specifies the ID of the workspace to which the training experiments belong.  
  If omitted, all training experiments in the region will be queried.

* `sort_by` - (Optional, String) Specifies the field used for sorting the training experiments.  
  Defaults to `create_time`.  
  The valid values are as follows:
  + **name**
  + **create_time**
  + **update_time**

* `order` - (Optional, String) Specifies the sort order of the training experiments.  
  Defaults to `desc`.  
  The valid values are as follows:
  + **asc**
  + **desc**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `training_experiments` - The list of training experiments that matched the filter parameters.  
  The [training_experiments](#modelarts_training_experiments) structure is documented below.

<a name="modelarts_training_experiments"></a>
The `training_experiments` block supports:

* `metadata` - The metadata of the training experiment.  
  The [metadata](#modelarts_training_experiments_metadata) structure is documented below.

* `statistic` - The statistics of the training experiment.  
  The [statistic](#modelarts_training_experiments_statistic) structure is documented below.

<a name="modelarts_training_experiments_metadata"></a>
The `metadata` block supports:

* `id` - The ID of the training experiment.

* `name` - The name of the training experiment.

* `description` - The description of the training experiment.

* `workspace_id` - The ID of the workspace to which the training experiment belongs.

* `create_time` - The creation time of the training experiment, in RFC3339 format.

* `update_time` - The update time of the training experiment, in RFC3339 format.

<a name="modelarts_training_experiments_statistic"></a>
The `statistic` block supports:

* `job_count` - The number of training jobs under the training experiment.
