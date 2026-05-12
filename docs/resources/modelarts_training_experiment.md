---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_experiment"
description: |-
  Manages a ModelArts training experiment resource within HuaweiCloud.
---

# huaweicloud_modelarts_training_experiment

Manages a ModelArts training experiment resource within HuaweiCloud.

## Example Usage

```hcl
variable "training_experiment_name" {}

resource "huaweicloud_modelarts_training_experiment" "test" {
  metadata {
    name = var.training_experiment_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the training experiment is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `metadata` - (Required, List) Specifies the configuration of the training experiment.  
  The [metadata](#training_experiment_metadata_arg) structure is documented below.

<a name="training_experiment_metadata_arg"></a>  
  The `metadata` block supports:

* `name` - (Required, String) Specifies the name of the training experiment.  
  The valid length is limited from `1` to `64`, only letters, digits, underscores (_) and hyphens (-) are allowed.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the ID of the workspace to which the training
  experiment belongs.  
  If omitted, the default workspace is used.

* `description` - (Optional, String) Specifies the description of the training experiment.  
  The maximum length is limited to `256` characters.  
  Only letters, Chinese characters, digits, spaces, underscores (_), hyphens (-), dots (.) and commas (,) are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the training experiment, in RFC3339 format.

* `update_time` - The latest update time of the training experiment, in RFC3339 format.

## Import

The training experiment can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_training_experiment.test <id>
```
