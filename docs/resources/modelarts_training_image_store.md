---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_image_store"
description: |-
  Use this resource to save a training job image to SWR within HuaweiCloud.
---

# huaweicloud_modelarts_training_image_store

Use this resource to save a training job image to SWR within HuaweiCloud.

~> Deleting this resource will only delete the stored SWR image tags, but will not delete the repository where the tags
  are located.

## Example Usage

```hcl
variable "training_job_id" {}
variable "task_id" {}
variable "image_name" {}
variable "swr_organization_name" {}

resource "huaweicloud_modelarts_training_image_store" "test" {
  training_job_id = var.training_job_id
  task_id         = var.task_id
  name            = var.image_name
  namespace       = var.swr_organization_name
  tag             = "v1.0.0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the training job is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `training_job_id` - (Required, String, NonUpdatable) Specifies the ID of the training job.

* `task_id` - (Required, String, NonUpdatable) Specifies the task name of the training job.  
  Can be obtained from the `status.tasks` field of the data source `huaweicloud_modelarts_training_jobs`.  

* `name` - (Required, String, NonUpdatable) Specifies the name of the SWR repository to which the image is stored.  
  The maximum length is `512` characters, and only lowercase letters, digits, hyphens (-), underscores (_) and dots (.)
  are allowed.

* `namespace` - (Required, String, NonUpdatable) Specifies the name of the SWR organization to which the image is stored.

* `tag` - (Required, String, NonUpdatable) Specifies the tag of the image.  
  The maximum length is `64` characters, and only letters, digits, hyphens (-), underscores (_) and dots (.)
  are allowed.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the image.  
  The maximum length is `512` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
