---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_batch_import_keypair"
description: |-
  Manages a KPS batch import keypair resource within HuaweiCloud.
---

# huaweicloud_kps_batch_import_keypair

Manages a KPS batch import keypair resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the imported keypairs,
but will only remove the resource information from the tfstate file.

## Example Usage

### Batch import keypair with private key

```hcl
variable "public_key" {}
variable "private_key" {}

resource "huaweicloud_kps_batch_import_keypair" "test" {
  keypairs {
    name       = "test-name"
    type       = "ssh"
    public_key = var.public_key
    scope      = "domain"

    key_protection {
      private_key = var.private_key
      encryption {
        type = "default"
      }
    }
  }
}
```

### Batch import keypair without private key

```hcl
variable "public_key" {}

resource "huaweicloud_kps_batch_import_keypair" "test" {
  keypairs {
    name       = "test-name"
    type       = "ssh"
    public_key = var.public_key
    scope      = "domain"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `keypairs` - (Required, List, NonUpdatable) Specifies the list of keypairs to import.
  The [keypairs](#Struct_keypairs) structure is documented below.

<a name="Struct_keypairs"></a>
The `keypairs` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the SSH keypair name.

* `type` - (Optional, String, NonUpdatable) Specifies the SSH keypair type. Valid values are **ssh** and **x509**.

* `public_key` - (Optional, String, NonUpdatable) Specifies the imported OpenSSH-formatted public key.

* `scope` - (Optional, String, NonUpdatable) Specifies the scope of the keypair. Valid values are **domain** and **user**.

* `user_id` - (Optional, String, NonUpdatable) Specifies the user ID of the keypair.

* `key_protection` - (Optional, List, NonUpdatable) Specifies the key protection configuration.
  This parameter can be configured with at most one element.
  The [key_protection](#Struct_key_protection) structure is documented below.

<a name="Struct_key_protection"></a>
The `key_protection` block supports:

* `private_key` - (Optional, String, NonUpdatable, Sensitive) Specifies the private key of the keypair.

* `encryption` - (Required, List, NonUpdatable) Specifies the encryption configuration.
  The [encryption](#Struct_encryption) structure is documented below.

<a name="Struct_encryption"></a>
The `encryption` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the encryption type. Valid values are **kms** and **default**.
  + **default**: The default encryption mode. Applicable to sites where KMS is not deployed.
  + **kms**: KMS encryption mode.

* `kms_key_name` - (Optional, String, NonUpdatable) Specifies the KMS key name.

* `kms_key_id` - (Optional, String, NonUpdatable) Specifies the KMS key ID.

-> If `type` is set to **kms**, you must enter the KMS key name or ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `succeeded_keypairs` - The list of successfully imported keypairs.
  The [succeeded_keypairs](#Struct_succeeded_keypairs) structure is documented below.

* `failed_keypairs` - The list of failed imported keypairs.
  The [failed_keypairs](#Struct_failed_keypairs) structure is documented below.

<a name="Struct_succeeded_keypairs"></a>
The `succeeded_keypairs` block supports:

* `name` - The SSH keypair name.

* `type` - The SSH keypair type.

* `public_key` - The public key of the keypair.

* `private_key` - The private key of the keypair.

* `fingerprint` - The fingerprint of the keypair.

* `user_id` - The user ID of the keypair.

<a name="Struct_failed_keypairs"></a>
The `failed_keypairs` block supports:

* `keypair_name` - The SSH keypair name.

* `failed_reason` - The failed reason.
