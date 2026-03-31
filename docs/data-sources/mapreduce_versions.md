---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_versions"
description: |-
  Use this data source to get available cluster versions of MapReduce within HuaweiCloud.
---

# huaweicloud_mapreduce_versions

Use this data source to get available cluster versions of MapReduce within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_mapreduce_versions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cluster versions are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - The list of available cluster versions.
