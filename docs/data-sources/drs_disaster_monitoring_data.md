---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_disaster_monitoring_data"
description: |-
  Use this data source to get disaster recovery monitoring data for specified DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_disaster_monitoring_data

Use this data source to get disaster recovery monitoring data for specified DRS jobs within HuaweiCloud.

## Example Usage

```hcl
variable "job_ids" { 
  type = list(string)
}

data "huaweicloud_drs_disaster_monitoring_data" "test" {
  job_ids = var.job_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_ids` - (Required, List) Specifies the list of DRS job IDs to query monitoring data for.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The collection of disaster recovery monitoring data response bodies.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The job ID.

* `data_guard_monitor` - The disaster recovery task monitoring data.

  The [data_guard_monitor](#data_guard_monitor_struct) structure is documented below.

<a name="data_guard_monitor_struct"></a>
The `data_guard_monitor` block supports:

* `bandwidth` - The bandwidth.

* `cpu_used_percent` - The CPU usage percentage.

* `dst_delay` - The destination database delay.

* `dst_io` - The destination IO.

* `dst_normal` - The destination database connection status.

* `dst_offset` - The destination database offset position.

* `dst_rps` - The destination RPS.

* `mem_used_in_mb` - The memory usage.

* `node_mem_in_mb` - The total node memory size.

* `node_offset` - The migration instance offset position.

* `node_volume_in_gb` - The total node disk size.

* `sr_delay` - The source database delay.

* `sr_offset` - The source database offset position.

* `src_io` - The source IO.

* `src_normal` - The source database connection status.

* `src_rps` - The source RPS.

* `trans_in_mb` - The migration data volume.

* `trans_lines` - The number of migrated data rows.

* `volume_used_in_gb` - The disk usage.

* `migration_bytes_per_second` - The number of bytes migrated per second.
