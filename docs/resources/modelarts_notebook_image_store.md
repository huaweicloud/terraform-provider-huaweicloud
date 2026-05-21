---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebook_image_store"
description: |-
  Use this resource to save a notebook instance image to ModelArts within HuaweiCloud.
---

# huaweicloud_modelarts_notebook_image_store

Use this resource to save a notebook instance image to ModelArts within HuaweiCloud.

-> This resource is only a one-time action resource for storing the notebook image in ModelArts side.
   Deleting this resource will not clear the corresponding image (version) record, but will only remove the resource
   information from the tfstate file.

~> Storing the image will map the original image used by the notebook to the current image, and will change the value of
   `image_id` parameter for the corresponding notebook resource. To avoid this change, it's a recommended to use the
   `lifecycle.ignore_changes` to ignore it.

## Example Usage

```hcl
variable "image_name" {}
variable "swr_organization_name" {}

resource "huaweicloud_modelarts_notebook" "test" {
  ...

  lifecycle {
    ignore_changes = [
      image_id,
    ]
  }
}

resource "huaweicloud_modelarts_notebook_image_store" "test" {
  notebook_id = huaweicloud_modelarts_notebook.test.id
  name        = var.image_name
  namespace   = var.swr_organization_name
  tag         = "v1.0.0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the notebook is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `notebook_id` - (Required, String, NonUpdatable) Specifies the ID of the notebook instance from which the image is saved.

* `name` - (Required, String, NonUpdatable) Specifies the name of the image.

* `namespace` - (Required, String, NonUpdatable) Specifies the SWR organization name of the image.

* `tag` - (Optional, String, NonUpdatable) Specifies the version tag of the image.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the image.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID to which the image belongs.  
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The image ID.

* `swr_path` - The SWR path of the image.

* `arch` - The processor architecture type supported by the image.

* `origin` - The origin of the image.

* `resource_categories` - The resource categories supported by the image.

* `service_type` - The service type supported by the image.

* `visibility` - The visibility of the image.

* `status` - The status of the image.

* `type` - The type of the image.

* `created_at` - The creation time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
