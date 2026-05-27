---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_health_compare_overview"
description: |-
  Use this data source to get the health compare object-level overview of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_health_compare_overview

Use this data source to get the health compare object-level overview of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}
variable "compare_job_id" {}

data "huaweicloud_drs_health_compare_overview" "test" {
  job_id         = var.job_id
  compare_job_id = var.compare_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the job ID.

* `compare_job_id` - (Required, String) Specifies the compare job ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compare_result` - The health compare object-level comparison result list.

  The [compare_result](#compare_result_struct) structure is documented below.

<a name="compare_result_struct"></a>
The `compare_result` block supports:

* `type` - The object type.  
  The valid values are as follows:
  + **DB**
  + **TABLE**
  + **VIEW**
  + **EVENT**
  + **ROUTINE**
  + **INDEX**
  + **TRIGGER**
  + **SYNONYM**
  + **FUNCTION**
  + **PROCEDURE**
  + **TYPE**
  + **RULE**
  + **DEFAULT_TYPE**
  + **PLAN_GUIDE**
  + **CONSTRAINT**
  + **FILE_GROUP**
  + **PARTITION_FUNCTION**
  + **PARTITION_SCHEME**
  + **TABLE_COLLATION**
  + **EXTENSIONS**

* `source_count` - The source count.

* `target_count` - The target count.

* `status` - The comparison result.  
  The valid values are as follows:
  + **0**: Inconsistent.
  + **2**: Consistent.
  + **3**: Incomplete.
