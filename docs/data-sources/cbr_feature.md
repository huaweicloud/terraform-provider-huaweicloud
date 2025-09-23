---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_feature"
description: |-
  Use this data source to query a specific feature of CBR within HuaweiCloud.
---

# huaweicloud_cbr_feature

Use this data source to query a specific feature of CBR within HuaweiCloud.

-> The API used by this datasource is currently in public beta and is temporarily unavailable in some regions.

## Example Usage

```hcl
data "huaweicloud_cbr_feature" "test" {
  feature_key = "replication.enable"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `feature_key` - (Required, String) Specifies the key of the feature to query.
  Valid values are:
  + **replication.enable**
  + **replication.source_region**
  + **replication.destination_regions**
  + **replication.destination_dgw_regions**
  + **features.backup_double_az**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `feature_value` - The value of the specified feature in JSON format.
