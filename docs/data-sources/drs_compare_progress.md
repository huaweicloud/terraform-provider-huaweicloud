---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_progress"
description: |-
  Use this data source to get the compare progress information of a DRS job.
---

# huaweicloud_drs_compare_progress

Use this data source to get the compare progress information of a DRS job.

## Example Usage

```hcl
variable "job_id" {} 
variable "compare_job_id" {}

data "huaweicloud_drs_compare_progress" "test" { 
  job_id         = var.job_id
  compare_job_id = var.compare_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `compare_job_id` - (Required, String) Specifies the ID of the compare job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `full_info` - The full comparison progress information. This field is returned for row comparison and content comparison.

  The [full_info](#full_info_struct) structure is documented below.

* `incre_info` - The incremental comparison progress information. This field is returned for dynamic content comparison.

  The [incre_info](#incre_info_struct) structure is documented below.

* `global_info` - The global comparison progress information.

  The [global_info](#global_info_struct) structure is documented below.

<a name="full_info_struct"></a>
The `full_info` block supports:

* `progress` - The full data comparison progress, in percentage (%).

* `src_speed` - The full data comparison speed.

* `recheck_entities` - The number of rows pending recheck for differences.

<a name="incre_info_struct"></a>
The `incre_info` block supports:

* `delay` - The incremental comparison delay. A value of 0 indicates that all incremental data has been compared.

* `src_speed` - The incremental data comparison speed.

* `rps` - The number of rows compared per second.

* `log_point` - The incremental log point.

* `recheck_entities` - The number of rows pending recheck for differences.

<a name="global_info_struct"></a>
The `global_info` block supports:

* `src_speed` - The global comparison speed.
