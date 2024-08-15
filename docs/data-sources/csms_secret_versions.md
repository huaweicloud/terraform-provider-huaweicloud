---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secret_versions"
description: |-
  Use this data source to get a list of the secret versions.
---

# huaweicloud_csms_secret_versions

Use this data source to get a list of the secret versions.

## Example Usage

```hcl
variable "secret_name" {}

data "huaweicloud_csms_secret_versions" "test" {
  secret_name = var.secret_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `secret_name` - (Required, String) Specifies the secret name to which the versions belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - The list of the versions.

  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `id` - The ID of the secret version.

* `kms_key_id` - The ID of the KMS key associated the secret.

* `secret_name` - The secret name to which the version belongs.

* `version_stages` - The secret version status list.

* `expire_time` - The expiration time of the secret version, in RFC3339 format.

* `created_at` - The creation time of the secret version, in RFC3339 format.
