---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_regions"
description: |-
  Use this data source to query the list of KMS regions supported by cross-regional keys within HuaweiCloud.
---

# huaweicloud_kms_key_regions

Use this data source to query the list of KMS regions supported by cross-regional keys within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_kms_key_regions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `regions` - The list of KMS regions supported by cross-regional keys.
