---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_keypairs"
description: |-
  Use this data source to get a list of keypairs.
---

# huaweicloud_kps_keypairs

Use this data source to get a list of keypairs.

## Example Usage

```hcl
variable "keypair_name" {}

data "huaweicloud_kps_keypairs" "test" {
  name = var.keypair_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the keypairs. If omitted, the provider-level
  region will be used.

* `name` - (Optional, String) Specifies the name of the keypair.

* `public_key` - (Optional, String) Specifies the imported OpenSSH-formatted public key.

* `fingerprint` - (Optional, String) Specifies the fingerprint of the keypair.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `keypairs` - The KPS keypairs list.

  The [keypairs](#kps-keypairs) structure is documented below.

<a name="kps-keypairs"></a>
The `keypairs` block supports:

* `name` - Indicates the name of the keypair.

* `scope` - Indicates the scope of keypair. The value can be **account**or **user**.

* `public_key` - Indicates the imported OpenSSH-formatted public key.

* `fingerprint` - Indicates the fingerprint information about a keypair.

* `is_managed` - Indicates whether the private key is managed by HuaweiCloud.
