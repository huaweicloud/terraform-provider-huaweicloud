---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_rpos_and_rtos"
description: |-
  Use this data source to batch query RPO and RTO of specified DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_batch_rpos_and_rtos

Use this data source to batch query RPO and RTO of specified DRS jobs within HuaweiCloud.

## Example Usage

```hcl
variable "jobs" {
  type = list(string)
}

data "huaweicloud_drs_batch_rpos_and_rtos" "test" {
  jobs = var.jobs
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `jobs` - (Required, List) Specifies the list of DRS job detail IDs to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The batch query RPO and RTO result list.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `job_id` - The job ID.

* `rpo_info` - The RPO information.

  The [rpo_info](#rpo_rto_info_struct) structure is documented below.

* `rto_info` - The RTO information.

  The [rto_info](#rpo_rto_info_struct) structure is documented below.

* `error_code` - The error code.

* `error_msg` - The error message.

<a name="rpo_rto_info_struct"></a>
The `rpo_info` and `rto_info` blocks support:

* `check_point` - The check point.

* `delay` - The delay in milliseconds.

* `gtid_set` - The GTID set.

* `time` - The current time in the format **yyyy-MM-dd HH:mm:ss**.
