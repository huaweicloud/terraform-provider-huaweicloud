---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_progress_data"
description: |-
  Use this data source to get the data-level streaming comparison list of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_progress_data

Use this data source to get the data-level streaming comparison list of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}
variable "type" {}

data "huaweicloud_drs_progress_data" "test" {
  job_id = var.job_id
  type   = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the job ID.

* `type` - (Required, String) Specifies the migration object type.  
  The valid values are as follows:
  + **table**
  + **event**
  + **table_structure**
  + **procedure**
  + **view**
  + **function**
  + **database**
  + **trigger**
  + **table_indexs**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `create_time` - The data generation time.

* `flow_compare_data` - The comparison result list.

  The [flow_compare_data](#flow_compare_data_struct) structure is documented below.

<a name="flow_compare_data_struct"></a>
The `flow_compare_data` block supports:

* `src_db` - The source database name.

* `src_tb` - The source object name.

* `dst_db` - The destination database name.

* `dst_tb` - The destination object name.

* `progress` - The progress percentage.
