---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_object_compare_detail"
description: |-
  Use this data source to get the object-level comparison details of a DRS task within HuaweiCloud.
---

# huaweicloud_drs_object_compare_detail

Use this data source to get the object-level comparison details of a DRS task within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_object_compare_detail" "test" {
  job_id       = "your_job_id"
  compare_type = "DB"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the task ID.

* `compare_type` - (Required, String) Specifies the object type.  
  The valid values are as follows:
  + **DB**: Database.
  + **TABLE**: Table.
  + **VIEW**: View.
  + **EVENT**: Event.
  + **ROUTINE**: Stored procedure and function.
  + **INDEX**: Index.
  + **TRIGGER**: Trigger.
  + **SYNONYM**: Synonym.
  + **FUNCTION**: Function.
  + **PROCEDURE**: Stored procedure.
  + **TYPE**: Custom type.
  + **RULE**: Rule.
  + **DEFAULT_TYPE**: Default value.
  + **PLAN_GUIDE**: Execution plan.
  + **CONSTRAINT**: Constraint.
  + **FILE_GROUP**: File group.
  + **PARTITION_FUNCTION**: Partition function.
  + **PARTITION_SCHEME**: Partition scheme.
  + **TABLE_COLLATION**: Table collation.

* `compare_job_id` - (Optional, String) Specifies the comparison task ID.
  If omitted, the latest comparison task information will be returned by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compare_details` - The list of object-level comparison details.

  The [compare_details](#compare_details_struct) structure is documented below.

<a name="compare_details_struct"></a>
The `compare_details` block supports:

* `source_db_name` - The source database name.

* `target_db_name` - The target database name.

* `source_db_value` - The value in the source database.

* `target_db_value` - The value in the target database.

* `status` - The comparison result.  
  The valid values are as follows:
  + **0**: Inconsistent.
  + **2**: Consistent.
  + **3**: Not completed.

* `error_message` - The error message.
