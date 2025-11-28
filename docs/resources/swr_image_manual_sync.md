---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_manual_sync"
description: |-
  Manages a SWR image manual sync resource within HuaweiCloud.
---

# huaweicloud_swr_image_manual_sync

Manages a SWR image manual sync resource within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {}
variable "repository_name" {}
variable "image_tag" {}
variable "target_region" {}
variable "target_organization" {}

resource "huaweicloud_swr_image_manual_sync" "test" {
  organization        = var.organization_name
  repository          = var.repository_name
  image_tag           = var.image_tag
  target_region       = var.target_region
  target_organization = var.target_organization
  override            = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `organization` - (Required, String, NonUpdatable) Specifies the name of the organization.

* `repository` - (Required, String, NonUpdatable) Specifies the name of the repository.

* `image_tag` - (Required, List, NonUpdatable) Specifies the iamge tags.

* `target_organization` - (Required, String, NonUpdatable) Specifies the target organization name.

* `target_region` - (Required, String, NonUpdatable) Specifies the target region name.

* `override` - (Optional, Bool, NonUpdatable) Specifies whether to overwrite. Default to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
