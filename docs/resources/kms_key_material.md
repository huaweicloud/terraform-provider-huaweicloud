---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_material"
description: ""
---

# huaweicloud_kms_key_material

Manages a KMS key material resource within HuaweiCloud.

-> NOTE: Please confirm that the state of the imported key is pending import.

## Example Usage

variable "key_id" {}
variable "import_token" {}
variable "encrypted_key_material" {}

```hcl
resource "huaweicloud_kms_key_material" "test" {
  key_id                 = var.key_id
  import_token           = var.import_token
  encrypted_key_material = var.encrypted_key_material
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, ForceNew) Specifies the ID of the KMS key.
  Changing this creates a new resource.

* `import_token` - (Required, String, ForceNew) Specifies the key import token in Base64 format.
  The value contains `200` to `6144` characters, including letters, digits, slashes(/) and equals(=). This value is
  obtained through the interface [Obtaining Key Import Parameters](https://support.huaweicloud.com/intl/en-us/api-dew/CreateParametersForImport.html).
  Changing this creates a new resource.

* `encrypted_key_material` - (Required, String, ForceNew) Specifies the encrypted symmetric key material in Base64 format.
  The value contains `344` to `360` characters, including letters, digits, slashes(/) and equals(=).
  If an asymmetric key is imported, this parameter is a temporary intermediate key used to encrypt the private key.
  This value is obtained refer to
  [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-dew/dew_01_0089.html).
  Changing this creates a new resource.

* `encrypted_privatekey` - (Optional, String, ForceNew) Specifies the private key encrypted using a temporary
  intermediate key. The value contains `200` to `6144` characters, including letters, digits, slashes(/)
  and equals(=). This parameter is required for importing an asymmetric key.
  Changing this creates a new resource.

* `expiration_time` - (Optional, String, ForceNew) Specifies the expiration time of the key material.
  This field is only valid for symmetric keys. The time is in the format of timestamp, that is, the
  offset seconds from 1970-01-01 00:00:00 UTC to the specified time.
  The time must be greater than the current time. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `key_id`.

* `key_usage` - The key usage. The value can be **ENCRYPT_DECRYPT** or **SIGN_VERIFY**.

* `key_state` - The status of the kms key. The valid values are as follows:
  **1**: To be activated
  **2**: Enabled.
  **3**: Disabled.
  **4**: Pending deletion.
  **5**: Pending import.

## Import

The KMS key material can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_kms_key_material.test 7056d636-ac60-4663-8a6c-82d3c32c1c64
```

Note that the imported state may not be identical to your resource definition,
due to `import_token`, `encrypted_key_material` and `encrypted_privatekey` are missing from the API response.
It is generally recommended running `terraform plan` after importing a KMS key material.
You can then decide if changes should be applied to the KMS key material, or the resource
definition should be updated to align with the KMS key material. Also you can ignore changes as below.

```hcl
resource "huaweicloud_kms_key_material" "test" {
    ...

  lifecycle {
    ignore_changes = [ import_token, encrypted_key_material, encrypted_privatekey ]
  }
}
```
