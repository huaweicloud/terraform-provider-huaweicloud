---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_monitor_data"
description: |-
  Use this data source to get monitoring data for a specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_job_monitor_data

Use this data source to get monitoring data for a specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_job_monitor_data" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth` - The EIP bandwidth. Unit: MB/s.

* `is_src_normal` - Whether the source database connection is normal.

* `is_dst_normal` - Whether the destination database connection is normal.

* `src_offset` - The source database offset position.

* `node_offset` - The migration instance offset position.

* `dst_offset` - The destination database offset position.

* `src_delay` - The source database delay.

* `dst_delay` - The destination database delay.

* `src_rps` - The source database RPS.

* `src_io` - The source database IO.

* `dst_rps` - The destination database RPS.

* `dst_io` - The destination database IO.

* `trans_data` - The amount of migrated data. Unit: MB.

* `trans_lines` - The number of migrated data rows.

* `used_volumes` - The disk usage. Unit: GB.

* `used_memory` - The memory usage. Unit: MB.

* `used_cpu_percent` - The CPU usage percentage.

* `node_volume_size` - The total node disk size. Unit: GB.

* `node_memory_size` - The total node memory size. Unit: MB.

* `update_time` - The update time.

* `apply_rate` - The synchronization speed. Unit: byte/s.
