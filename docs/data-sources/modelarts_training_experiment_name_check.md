---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_experiment_name_check"
description: |-
  Use this data source to check whether the ModelArts training experiment name is duplicate.
---

# huaweicloud_modelarts_training_experiment_name_check

Use this data source to check whether the specified training experiment name already exists
within HuaweiCloud ModelArts service.

## Example Usage

```hcl
variable "experiment_name" {}

data "huaweicloud_modelarts_training_experiment_name_check" "test" {
  experiment_name = var.experiment_name
}
```

### Check with a specific workspace

```hcl
variable "experiment_name" {}
variable "workspace_id" {}

data "huaweicloud_modelarts_training_experiment_name_check" "test" {
  experiment_name = var.experiment_name
  workspace_id    = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the ModelArts service is located.
  If omitted, the provider-level region will be used.

* `experiment_name` - (Required, String) Specifies the name of the training experiment to be checked.
  The name consists of `1` to `64` characters. Only letters, digits, and hyphens (-) are allowed.

* `workspace_id` - (Optional, String) Specifies the ID of the workspace to which the training experiment belongs.
  If omitted, the default workspace (`0`) will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_duplicate` - Whether the training experiment name is duplicate.
  + **true** - The name already exists.
  + **false** - The name does not exist.
