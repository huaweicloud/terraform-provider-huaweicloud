---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_versions"
description: ""
---

# huaweicloud_mapreduce_versions

Use this data source to get available cluster versions of MapReduce.

## Example Usage

```hcl
data "huaweicloud_mapreduce_versions" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - List of available cluster versions.
