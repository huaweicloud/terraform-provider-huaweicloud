---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_registered_region"
description: |-
  Manages an Identity Center registered region resource within HuaweiCloud.
---

# huaweicloud_identitycenter_registered_region

Manages an Identity Center registered region resource within HuaweiCloud.

## Example Usage

```hcl
variable "region_id" {}

resource "huaweicloud_identitycenter_registered_region" "test" {
  region_id = var.region_id
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Required, String, ForceNew) Specifies the ID of the region.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The Identity Center registered region can be imported using the `region_id`, e.g.

```bash
$ terraform import huaweicloud_identitycenter_registered_region.test <region_id>
```
