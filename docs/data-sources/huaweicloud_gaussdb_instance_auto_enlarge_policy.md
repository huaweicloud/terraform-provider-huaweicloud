---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_auto_enlarge_policy"
description: |-
  Use this data source to query the storage auto scaling policy of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_auto_enlarge_policy

Use this data source to query the storage auto-scaling policy of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_instance_auto_enlarge_policy" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String) Specifies the GaussDB instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `switch_option` - Whether the storage auto-scaling is enabled.

* `limit_volume_size` - The upper limit of storage auto scaling (GB).

* `min_volume_size` - The minimum disk capacity to expand (GB).

* `max_volume_size` - The maximum disk capacity to expand (GB).

* `trigger_available_percent` - The available storage space percentage threshold.

* `percents` - The list of available storage space percentages.

* `step_size` - The expansion step size for fixed size expansion (GB).

* `step_percent` - The expansion step size for percentage expansion.
