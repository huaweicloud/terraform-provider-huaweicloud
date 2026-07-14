---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_kms_keys"
description: |-
  Use this data source to query the KMS keys available for GaussDB transparent data encryption within HuaweiCloud.
---

# huaweicloud_gaussdb_kms_keys

Use this data source to query the KMS keys available for GaussDB transparent data encryption within HuaweiCloud.

## Example Usage

```hcl
variable "kms_project_name" {}

data "huaweicloud_gaussdb_kms_keys" "test" {
  kms_project_name = var.kms_project_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the KMS keys.
  If omitted, the provider-level region will be used.

* `kms_project_name` - (Required, String) Specifies the name of the resource space where the
  KMS master key ID used by GaussDB for transparent data encryption is located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `key_details` - The KMS key details.
  The [key_details](#kms_keys_key_details) structure is documented below.

* `authorized_id` - The authorized user ID. This must be set for the current user when enabling
  transparent data encryption.

* `authorized_name` - The authorized username.

<a name="kms_keys_key_details"></a>
The `key_details` block supports:

* `key_id` - The key ID.

* `default_key_flag` - The default master key flag.
  + **1**: default master key.
  + **0**: non-default master key.

* `key_alias` - The key alias.

* `key_spec` - The key generation algorithm.
  The valid values are as follows:
  + **AES_256**
  + **SM4**
  + **RSA_2048**
  + **RSA_3072**
  + **RSA_4096**
  + **EC_P256**
  + **EC_P384**
  + **SM2**
  + **ALL**

* `domain_id` - The user domain ID.

* `key_state` - The key status.
  The valid values are as follows:
  + **1**: to be activated.
  + **2**: enabled.
  + **3**: disabled.
  + **4**: pending deletion.
  + **5**: pending import.
