---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vaults_summary"
description: |-
  Use this data source to get CBR total capacity and used capacity of all vaults within HuaweiCloud.
---

# huaweicloud_cbr_vaults_summary

Use this data source to get CBR total capacity and used capacity of all vaults within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_vaults_summary" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `size` - The allocated capacity for the associated resource, in GB.

* `used_size` - The used capacity, in GB.
