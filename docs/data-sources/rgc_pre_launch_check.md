---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_pre_launch_check"
description: |-
  Use this data source to get the pre-launch check result in Resource Governance Center.
---

# huaweicloud_rgc_pre_launch_check

Use this data source to get the pre-launch check result in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_pre_launch_check" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pre_launch_check` - Check if the current region can be set up for landing zone.
