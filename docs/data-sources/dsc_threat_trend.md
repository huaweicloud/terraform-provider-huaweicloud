---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_threat_trend"
description: |-
  Use this data source to get the threat trend information within HuaweiCloud.
---

# huaweicloud_dsc_threat_trend

Use this data source to get the threat trend information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_threat_trend" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `api_attacked_variation` - The variation list of the number of API attacks.

* `database_attacked_variation` - The variation list of the number of database attacks.

* `time_axis` - The time axis.
