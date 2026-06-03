---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_job_engines"
description: |-
  Use this data source to get AI preset framework engine list supported by the ModelArts training job
  within HuaweiCloud.
---

# huaweicloud_modelarts_training_job_engines

Use this data source to get AI preset framework engine list supported by the ModelArts training job within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_modelarts_training_job_engines" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the training job engines are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `engines` - The list of training job engines.  
  The [engines](#modelarts_training_job_engines_attr) structure is documented below.

<a name="modelarts_training_job_engines_attr"></a>
The `engines` block supports:

* `engine_id` - The ID of the engine.

* `engine_name` - The name of the engine.

* `engine_version` - The version of the engine.

* `v1_compatible` - Whether the engine is v1 compatible.

* `run_user` - The default startup user UID of the engine.

* `image_info` - The image information of the engine.  
  The [image_info](#modelarts_training_job_engines_image_info) structure is documented below.

<a name="modelarts_training_job_engines_image_info"></a>
The `image_info` block supports:

* `cpu_image_url` - The CPU image URL of the engine.

* `gpu_image_url` - The GPU or Ascend image URL of the engine.

* `image_version` - The image version of the engine.
