---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_keystone_metadata_file"
description: |-
    Use this data source to query the Metadata file of Keystone within HuaweiCloud.
---

# huaweicloud_identity_keystone_metadata_file

Use this data source to query the Metadata file of Keystone within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identity_keystone_metadata_file" "test1" {}

data "huaweicloud_identity_keystone_metadata_file" "test2" {
  unsigned = true
}
```

## Argument Reference

* `unsigned` - (Optional, Bool) Specifies whether to sign the metadata according to the SAML2.0 specification, default
  is false.

## Attribute Reference

* `metadata_file` - Indicates the keystone metadata file.
