---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_sha256"
description: |-
  Use this data source to get SMS SHA256 hash value.
---

# huaweicloud_sms_sha256

Use this data source to get SMS SHA256 hash value.

## Example Usage

```hcl
variable "key" {}

data "huaweicloud_sms_sha256" "test" {
  key = var.key
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, String) Specifies the keyword.
  
  -> The value of the encrypted field must be in the UUID format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - Indicates the SHA256 hash value.
