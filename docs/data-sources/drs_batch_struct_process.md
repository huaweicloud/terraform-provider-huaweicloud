---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_struct_process"
description: |-
  Use this data source to get the batch structure processing results of DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_batch_struct_process

Use this data source to get the batch structure processing results of DRS jobs within HuaweiCloud.

## Example Usage

```hcl
variable "job_ids" { 
  type = list(string)
}

data "huaweicloud_drs_batch_struct_process" "test" {
  job_ids = var.job_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_ids` - (Required, List) Specifies the list of DRS job IDs for batch query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The list of disaster recovery initialization progress returned by batch query.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `job_id` - The ID of the DRS job.

* `struct_process` - The disaster recovery initialization progress information.

  The [struct_process](#struct_process_struct) structure is documented below.

* `error_code` - The error code.

* `error_message` - The error message.

* `database_info` - The business/disaster recovery database information during the comparison task.

  The [database_info](#database_info_struct) structure is documented below.

<a name="struct_process_struct"></a>
The `struct_process` block supports:

* `result` - The comparison result.

  The [result](#result_struct) structure is documented below.

* `create_time` - The time when the data was generated.

<a name="result_struct"></a>
The `result` block supports:

* `type` - The type of database object.  
  The valid values are as follows:
  + **DATABASE**: Database.
  + **SCHEMA**: Schema.
  + **PACKAGE**: Package.
  + **TABLE**: Table.
  + **COLUMN**: Column.
  + **VIEW**: View.
  + **FUNCTION**: Function.
  + **PROCEDURE**: Stored procedure.
  + **ROUTINE**: Routine.
  + **TRIGGER**: Trigger.
  + **INDEX**: Index.
  + **TABLE_INDEX**: Normal index, aggregated by table.
  + **TABLE_RENAME_OR_COPY**: Table rename or copy.
  + **TABLE_STRUCTURE**: Table structure.
  + **EVENT**: Event.
  + **SYNONYM**: Synonym, specific to SQL Server.
  + **TYPE**: Custom type, specific to SQL Server.
  + **RULE**: Rule, specific to SQL Server.
  + **DEFAULT**: Default value, specific to SQL Server.
  + **PLAN_GUIDE**: Execution plan, specific to SQL Server.
  + **FILE_GROUP**: File group, specific to SQL Server.
  + **PARTITION_FUNCTION**: Partition function, specific to SQL Server.
  + **SHARD_KEY**: Shard key, specific to MongoDB.
  + **VALIDATOR**: Validator, specific to MongoDB.
  + **SEQUENCE**: Sequence.
  + **MATVIEW**: Materialized view.
  + **PARTITION_SCHEME**: Partition scheme, specific to SQL Server.
  + **ACCOUNT**: Account.
  + **EXTENSION**: Extension, specific to PostgreSQL.
  + **AGGREGATE**: Aggregate function, specific to PostgreSQL.
  + **MATERIALIZED_VIEW**: Materialized view, specific to PostgreSQL.
  + **TEXT_SEARCH_DICTIONARY**: Text search dictionary, specific to PostgreSQL.
  + **CONVERSION**: Type conversion, specific to PostgreSQL.
  + **DATA_TYPE**: Data type, specific to PostgreSQL.
  + **TEXT_SEARCH_CONFIGURATION**: Text search configuration, specific to PostgreSQL.
  + **STATISTICS_EXTENSION**: Statistics extension, specific to PostgreSQL.
  + **MEMBERSHIP**: User membership, specific to PostgreSQL.
  + **EVENT_TRIGGER**: Event trigger, specific to PostgreSQL.
  + **COLLATION**: Collation, specific to PostgreSQL.
  + **TEXT_SEARCH_PARSER**: Text search parser, specific to PostgreSQL.
  + **PRIVILEGES**: Privileges, specific to PostgreSQL.
  + **FOREIGN_KEY**: Foreign key, specific to PostgreSQL.
  + **ROLE**: Role.

* `status` - The status of structure processing.  
  The valid values are as follows:
  + **1**: Comparing.
  + **2**: Comparison completed.

* `src_count` - The number of source objects.

* `dst_count` - The number of destination objects.

* `start_time` - The start time (Unix timestamp in milliseconds).

* `end_time` - The end time (Unix timestamp in milliseconds).

<a name="database_info_struct"></a>
The `database_info` block supports:

* `service_database` - The business database information.

* `disaster_recovery_database` - The disaster recovery database information.
