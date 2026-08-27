---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_select_objects"
description: |-
  Manages a resource to batch select DRS job objects within HuaweiCloud.
---

# huaweicloud_drs_batch_select_objects

Manages a resource to batch select DRS job objects within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch select DRS job objects. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. Only real-time migration and real-time synchronization jobs support object selection.
  The job status must be **CONFIGURATION**. The connection tests to the source and destination databases must succeed,
  and the job modification API must be called successfully before using this resource.
  <br/>3. The execution result of this operation is based on the `status` field in the `results` block.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_batch_select_objects" "test" {
  jobs {
    job_id   = var.job_id
    selected = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the list of jobs for batch object selection.
  To ensure the API calling performance, the number of jobs in a batch call is recommended not to exceed `10`.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `selected` - (Optional, Bool, NonUpdatable) Specifies whether to select objects customarily.  
  The valid values are as follows:
  + **true**: Custom objects. The `job` parameter is required to specify the objects.
  + **false**: Migrate all objects from the source database.

* `sync_database` - (Optional, Bool, NonUpdatable) Specifies whether to synchronize objects at the database level.  
  The valid values are as follows:
  + **true**: Database-level synchronization.
  + **false**: Non-database-level synchronization.

* `job` - (Optional, List, NonUpdatable) Specifies the database objects to migrate or synchronize.
  This parameter is required when `selected` is **true**.

  The [job](#job_struct) structure is documented below.

<a name="job_struct"></a>
The `job` block supports:

* `id` - (Optional, String, NonUpdatable) Specifies the object ID.
  + When `object_type` is **database**, this is the database name.
  + When `object_type` is **table** or **view**, refer to the API example for the field value. Please refer to
    the [document](https://support.huaweicloud.com/intl/en-us/api-drs/drs_03_0109.html).

* `parent_id` - (Optional, String, NonUpdatable) Specifies the parent object ID.
  This parameter is required when `object_type` is **table** or **view**, and the value is the database name.

* `object_type` - (Optional, String, NonUpdatable) Specifies the object type.  
  The valid values are as follows:
  + **database**: Database.
  + **table**: Table.
  + **schema**: Schema.
  + **view**: View.

* `object_name` - (Optional, String, NonUpdatable) Specifies the object name, which can be a database name, table name,
  or view name.

* `select` - (Optional, String, NonUpdatable) Specifies whether the object is selected for migration.  
  The valid values are as follows:
  + **true**: The object will be migrated.
  + **false**: The object will not be migrated.
  + **partial**: Partially selected. Migrate some tables under the database.

  Defaults to **false**.

* `object_alias_name` - (Optional, String, NonUpdatable) Specifies the alias of the object after mapping.
  This parameter is supported only for synchronization tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `all_counts` - The total number of object selection tasks in this request.

* `results` - The results of the batch object selection.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `job_id` - The job ID.

* `status` - Whether the object selection task succeeded.
  The value can be **true** or **false**.

* `error_code` - The error code when the task fails.

* `error_msg` - The error message when the task fails.
