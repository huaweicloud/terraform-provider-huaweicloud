---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_running_instance_statistics"
description: |-
  Use this data source to query the running instance statistics under a specified region within HuaweiCloud.
---

# huaweicloud_dcs_running_instance_statistics

Use this data source to query the running instance statistics under a specified region within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dcs_running_instance_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the running instance statistics.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - The list of running instance statistics.  
  The [statistics](#dcs_running_instance_statistics_statistics) structure is documented below.

<a name="dcs_running_instance_statistics_statistics"></a>
The `statistics` block supports:

* `instance_id` - The ID of the cache instance.

* `input_kbps` - The network input traffic of the cache instance, in kbit/s.

* `output_kbps` - The network output traffic of the cache instance, in kbit/s.

* `keys` - The number of cached data entries.

* `used_memory` - The used cache memory, in MB.

* `max_memory` - The total cache memory, in MB.

* `cmd_get_count` - The number of times the cache get command has been called.

* `cmd_set_count` - The number of times the cache set command has been called.

* `used_cpu` - The cumulative CPU time consumed by the cache in user and kernel states, in seconds.
