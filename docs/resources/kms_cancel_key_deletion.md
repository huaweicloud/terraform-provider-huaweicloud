---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_cancel_key_deletion"
description: |-
  Manages a KMS cancel key deletion resource within HuaweiCloud.
---

# huaweicloud_kms_cancel_key_deletion

Manages a KMS cancel key deletion resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the cancellation of key deletion,
  but will only remove the resource information from the tfstate file.

-> This resource only supports canceling key deletion for keys with `key_state` set to **4** (pending deletion).

## Example Usage

```hcl
variable "key_id" {}

resource "huaweicloud_kms_cancel_key_deletion" "test" {
  key_id = var.key_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID.
  The valid length is `36` bytes, meeting regular match **^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$**.
  For example: **0d0466b0-e727-4d9c-b35d-f84bb474a37f**.

* `sequence` - (Optional, String, NonUpdatable) Specifies the sequence number of the request message, `36` bytes.
  For example: **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `key_state` - The current status of the KMS key.
  The valid values are as follows:
  + **2**: Enabled.
  + **3**: Disabled.
  + **4**: Pending deletion.
  + **5**: Pending import.
  + **7**: Frozen.
