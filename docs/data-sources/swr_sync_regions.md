---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_sync_regions"
description: |-
  Use this data source to get the list of regions which is available to synchronize images.
---

# huaweicloud_swr_sync_regions

Use this data source to get the list of regions which are available to synchronize images.

## Example Usage

```hcl
data "huaweicloud_swr_sync_regions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sync_regions` - The region IDs that are available to synchronize images.
