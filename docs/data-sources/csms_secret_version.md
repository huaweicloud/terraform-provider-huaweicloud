---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secret_version"
description: ""
---

# huaweicloud_csms_secret_version

Use this data source to query the version and plaintext of the CSMS(Cloud Secret Management Service) secret.

## Example Usage

```hcl
data "huaweicloud_csms_secret_version" "version_1" {
  secret_name = "your_secret_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CSMS secrets.
  If omitted, the provider-level region will be used.

* `secret_name` - (Required, String) The name of the CSMS secret to query.

* `version` - (Optional, String) The version ID of the CSMS secret version to query.
  If omitted, the latest version will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `secret_text` - The plaintext of a secret in text format.

* `kms_key_id` - The ID of the KMS CMK used for secret encryption.

* `status` - The status of the CSMS secret version.

* `created_at` - Time when the CSMS secret version created, in UTC format.
