---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_backup_migration_jobs"
description: |-
  Use this data source to query the list of DRS backup migration jobs within HuaweiCloud.
---

# huaweicloud_drs_backup_migration_jobs

Use this data source to query the list of DRS backup migration jobs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_backup_migration_jobs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the job name. Fuzzy search is supported.

* `status` - (Optional, String) Specifies the backup migration job status.  
  The valid values are as follows:
  + **TRANSFERRING**: Restoring.
  + **SUCCESS**: Successful.
  + **FAILED**: Failed.
  + **PRECHECK FAILED**: Pre-check failed.

* `dbs_instance_ids` - (Optional, List) Specifies the database instance IDs. A maximum of `10` IDs are supported.

* `description` - (Optional, String) Specifies the description.

* `create_at` - (Optional, String) Specifies the creation time.

* `completed_at` - (Optional, String) Specifies the completed time.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `tags` - (Optional, String) Specifies the tags.

* `sort_key` - (Optional, String) Specifies the sort field.  
  The valid values are as follows:
  + **name**
  + **db_type**
  + **inst_id**
  + **ip**
  + **created_at**
  + **description**
  + **tag_value**

* `sort_dir` - (Optional, String) Specifies the sort method.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `count` - The total number of backup migration jobs.

* `jobs` - The backup migration job list.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - The job ID.

* `name` - The job name.

* `status` - The job status.  
  The valid values are as follows:
  + **TRANSFERRING**: Restoring.
  + **SUCCESS**: Successful.
  + **FAILED**: Failed.
  + **PRECHECK FAILED**: Pre-check failed.

* `engine_type` - The database engine.  
  The valid values are as follows:
  + **sqlserver**: RDS for SQL Server engine.

* `error_log` - The error log.

* `description` - The description.

* `create_time` - The job create time.

* `finish_time` - The job finish time.

* `enterprise_project_id` - The enterprise project ID.
