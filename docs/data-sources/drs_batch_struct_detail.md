---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_struct_detail"
description: |-
  Use this data source to get the details of DR initialization objects within HuaweiCloud.
---

# huaweicloud_drs_batch_struct_detail

Use this data source to get the details of DR initialization objects within HuaweiCloud.

## Example Usage

```hcl
variable "type" {}
variable "jobs" {
  type = list(string)
}

data "huaweicloud_drs_batch_struct_detail" "test" {
  type = var.type
  jobs = var.jobs
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the database migration object type.
  The valid values are as follows:
  + **database**
  + **schema**
  + **table**
  + **view**
  + **procedure**
  + **trigger**
  + **index**
  + **table_indexs**
  + **table_structure**

* `jobs` - (Required, List) Specifies the list of DRS job detail IDs to query.

* `cur_page` - (Optional, Int) Specifies the current page number.

* `per_page` - (Optional, Int) Specifies the number of items per page.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The list of batch struct details.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `job_id` - The task ID.

* `error_code` - The error code.

* `error_message` - The error message.

* `struct_detail` - The details of the DR initialization object.

  The [struct_detail](#struct_detail_struct) structure is documented below.

<a name="struct_detail_struct"></a>
The `struct_detail` block supports:

* `total_record` - The total number of tasks.

* `create_time` - The time when the data was generated.

* `list` - The list of comparison results.

  The [list](#list_struct) structure is documented below.

<a name="list_struct"></a>
The `list` block supports:

* `progress` - The progress of the task.

* `src_db` - The name of the source database.
   If the source database has a three-level structure, the format is: **database.schema**.

* `src_tb` - The name of the source object.

* `dst_db` - The name of the target database.

* `dst_tb` - The name of the target object.
