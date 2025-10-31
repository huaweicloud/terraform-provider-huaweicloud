---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_landing_zone_available_updates"
description: |-
  Use this data source to check if the Landing Zone currently requires an upgrade in Resource Governance Center.
---

# huaweicloud_rgc_landing_zone_available_updates

Use this data source to check if the Landing Zone currently requires an upgrade in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_landing_zone_available_updates" "test" {
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `baseline_update_available` - Check if the basic configuration of the Landing Zone is available.

* `control_update_available` - Check if there are any new control policies under the current account.

* `landing_zone_update_available` - Check if the Landing Zone is updatable.

* `service_landing_zone_version` - The latest version number of Landing Zone.

* `user_landing_zone_version` - The current version of the user's Landing Zone.
