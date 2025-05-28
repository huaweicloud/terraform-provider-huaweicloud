---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_public_key"
description: |-
  Use this data source to get the information of the public key which is a specified asymmetric key.
---

# huaweicloud_kms_public_key

Use this data source to get the information of the public key which is a specified asymmetric key.

## Example Usage

```hcl
variable "key_id" {}

data "huaweicloud_kms_public_key" "test" {
  key_id = var.key_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `key_id` - (Required, String) Specifies the ID of the key which is `36` bytes and meeting regular matching as
  `^ [0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`. For example,
  **0d0466b0-e727-4d9c-b35d-f84bb474a37f**.

* `sequence` - (Optional, String) Specifies the request sequence number which is a `36` byte serial number.  
  For example, **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `public_key` - The information of the public key.
