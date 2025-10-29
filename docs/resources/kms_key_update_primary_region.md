---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_update_primary_region"
description: |-
  Manages a resource to update KMS key primary region within HuaweiCloud.
---

# huaweicloud_kms_key_update_primary_region

Manages a resource to update KMS key primary region within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "key_id" {}
variable "primary_region" {}

resource "huaweicloud_kms_key_update_primary_region" "test" {
  key_id         = var.key_id
  primary_region = var.primary_region
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID.

* `primary_region` - (Required, String, NonUpdatable) Specifies the area code of the new primary region to which
  the key belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
