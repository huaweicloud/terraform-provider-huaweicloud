---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_cloudlog_regions"
description: |-
  Use this data source to get the list of regions.
---

# huaweicloud_secmaster_collector_cloudlog_regions

Use this data source to get the list of regions.

## Example Usage

```hcl
data "huaweicloud_secmaster_collector_cloudlog_regions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of the regions.
