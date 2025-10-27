---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_batch_export_private_key"
description: |-
  Manages a KPS batch export private key resource within HuaweiCloud.
---

# huaweicloud_kps_batch_export_private_key

Manages a KPS batch export private key resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the exported keys,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
data "huaweicloud_kps_keypairs" "test" {}

locals {
  name              = data.huaweicloud_kps_keypairs.test.keypairs.0.name
  type              = data.huaweicloud_kps_keypairs.test.keypairs.0.type
  scope             = data.huaweicloud_kps_keypairs.test.keypairs.0.scope
  public_key        = data.huaweicloud_kps_keypairs.test.keypairs.0.public_key
  fingerprint       = data.huaweicloud_kps_keypairs.test.keypairs.0.fingerprint
  is_key_protection = data.huaweicloud_kps_keypairs.test.keypairs.0.is_managed
  frozen_state      = data.huaweicloud_kps_keypairs.test.keypairs.0.frozen_state
}

resource "huaweicloud_kps_batch_export_private_key" "test" {
  export_file_name = "./keypairs.zip"

  keypairs {
    name              = local.name
    type              = local.type
    scope             = local.scope == "account" ? "domain" : local.scope
    public_key        = local.public_key
    fingerprint       = local.fingerprint
    is_key_protection = local.is_key_protection
    frozen_state      = local.frozen_state
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `export_file_name` - (Required, String, NonUpdatable) Specifies the directory file for storing exported key pairs,
  and requires the file ending in `.zip`.

* `keypairs` - (Required, List, NonUpdatable) Specifies the list of keypairs to export.
  The [keypairs](#keypairs) structure is documented below.

<a name="keypairs"></a>
The `keypairs` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the SSH keypair name.

* `type` - (Optional, String, NonUpdatable) Specifies the SSH keypair type. Valid values are **ssh** and **x509**.

* `scope` - (Optional, String, NonUpdatable) Specifies the scope of the keypair. Valid values are **domain** and **user**.

* `public_key` - (Optional, String, NonUpdatable) Specifies the imported OpenSSH-formatted public key.

* `fingerprint` - (Optional, String, NonUpdatable) Specifies the fingerprint of the keypair.

* `is_key_protection` - (Optional, Bool, NonUpdatable) Specifies whether the key is protected.

* `frozen_state` - (Optional, String, NonUpdatable) Specifies the frozen state of the keypair. Valid values are:
  + **0**: normal, not frozen
  + **1**: frozen due to common causes
  + **2**: frozen by the public security bureau
  + **3**: frozen due to common causes and by the public security bureau
  + **4**: frozen due to violations
  + **5**: frozen due to common causes and violations
  + **6**: frozen by the public security bureau and due to violations
  + **7**: frozen by the public security bureau and due to common causes and violations
  + **8**: frozen due to lack of real-name authentication
  + **9**: frozen due to common causes and lack of real-name authentication
  + **10**: frozen by the public security bureau and due to lack of real-name authentication

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
