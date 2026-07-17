---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_kms_detail"
description: |-
  Use this data source to query the detail of a KMS key used by GaussDB for transparent data encryption within HuaweiCloud.
---

# huaweicloud_gaussdb_kms_detail

Use this data source to query the detail of a KMS key used by GaussDB for transparent data encryption within HuaweiCloud.

## Example Usage

```hcl
variable "kms_project_name" {}
variable "key_id" {}

data "huaweicloud_gaussdb_kms_detail" "test" {
  kms_project_name = var.kms_project_name
  key_id           = var.key_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the KMS key detail.
  If omitted, the provider-level region will be used.

* `kms_project_name` - (Required, String) Specifies the name of the resource space where the
  KMS master key used by GaussDB for transparent data encryption is located.

* `key_id` - (Required, String) Specifies the KMS key ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `key_alias` - The key alias.

* `domain_id` - The user domain ID.

* `key_state` - The key status.
  The valid values are as follows:
  + **1**: to be activated.
  + **2**: enabled.
  + **3**: disabled.
  + **4**: pending deletion.
  + **5**: pending import.
