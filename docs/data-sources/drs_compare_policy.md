---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_policy"
description: |-
  Use this data source to get the compare policy information of a DRS job.
---

# huaweicloud_drs_compare_policy

Use this data source to get the compare policy information of a DRS job.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_compare_policy" "test" {
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

* `interval_hour` - The interval hours between comparisons.

* `period` - The comparison period schedule.

* `status` - The status of the compare policy.
  The valid values are as follows:
  + **OPEN**: Enabled.
  + **CLOSED**: Disabled, no comparison policy is set.
  + **NO_SUPPORT**: No data available.

* `begin_time` - The start time of the comparison window.

* `end_time` - The end time of the comparison window.

* `compare_type` - The list of comparison types.
  The valid values are as follows:
  + **object**: Object comparison.
  + **lines**: Row count comparison.
  + **account**: User comparison.

* `next_compare_time` - The scheduled time for the next comparison in UTC format, e.g., **2023-06-12T08:00:00Z**.

* `compare_policy` - The comparison policy type.
  The valid values are as follows:
  + **normal**: Normal comparison.
  + **manyToOne**: Many-to-one comparison.
