---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_retire_grant"
description: |-
  Manages a resource to retire grant within HuaweiCloud.
---

# huaweicloud_kms_retire_grant

Manages a resource to retire grant within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "key_id" {}
variable "grant_id" {}

resource "huaweicloud_kms_retire_grant" "test" {
  key_id   = var.key_id
  grant_id = var.grant_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID.

* `grant_id` - (Required, String, NonUpdatable) Specifies the grant ID.

* `sequence` - (Optional, String, NonUpdatable) Specifies the sequence number of the request message, `36` bytes.
  For example: **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
