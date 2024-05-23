---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_upgrade_target_images"
description: |-
  Use this data source to get the list of the CSS upgrade target images.
---

# huaweicloud_css_upgrade_target_images

Use this data source to get the list of the CSS upgrade target images.

## Example Usage

```hcl
variable "cluster_id" {}
variable "upgrade_type" {}

data "huaweicloud_css_upgrade_target_images" "test" {
  cluster_id   = var.cluster_id
  upgrade_type = var.upgrade_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster to be upgraded.

* `upgrade_type` - (Required, String) Specifies the upgrade type.
  + **same:** the same version.
  + **cross:** the cross version.

* `image_id` - (Optional, String) Specifies the ID of the target image.

* `engine_type` - (Optional, String) Specifies the datastore type of the target image.

* `engine_version` - (Optional, String) Specifies the datastore version of the target image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `images` - The list of the upgrade target images.

  The [images](#images_struct) structure is documented below.

<a name="images_struct"></a>
The `images` block supports:

* `id` - The target image ID that can be upgraded.

* `name` - The name of the target image that can be upgraded.

* `engine_type` - The image datastore type.

* `engine_version` - The image datastore version.

* `priority` - The target image priority.

* `description` - The image description information.
