---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_images_image_share"
description: |-
  Manages an IMS image share resource within HuaweiCloud.
---

# huaweicloud_images_image_share

Manages an IMS image share resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_image_id" {}
variable "target_project_ids" {}

resource "huaweicloud_images_image_share" "test" {
  source_image_id    = var.source_image_id
  target_project_ids = var.target_project_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_image_id` - (Required, String, ForceNew) Specifies the ID of the source image. The source image must be in the
  same region as the current resource.

  Changing this parameter will create a new resource.

* `target_project_ids` - (Required, List) Specifies the IDs of the target projects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `source_image_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.
