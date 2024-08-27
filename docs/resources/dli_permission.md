---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_permission"
description: ""
---

# huaweicloud_dli_permission

Manages the usage permissions of those resources: `huaweicloud_dli_queue`, `huaweicloud_dli_database`,
 `huaweicloud_dli_table`, `huaweicloud_dli_package`, `huaweicloud_dli_flinksql_job`, `huaweicloud_dli_flinkjar_job`
  within HuaweiCloud DLI.

## Example Usage

### Grant a permission of queue

```hcl
variable "user_name" {}
variable "queue_name" {}

resource "huaweicloud_dli_permission" "test" {
  user_name  = var.user_name
  object     = "queues.${var.queue_name}"
  privileges = ["SUBMIT_JOB","DROP_QUEUE"]
}
```

### Grant a permission of database

```hcl
variable "user_name" {}
variable "database_name" {}

resource "huaweicloud_dli_permission" "test" {
  user_name  = var.user_name
  object     = "databases.${var.database_name}"
  privileges = ["SELECT"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DLI permission resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies name of the user who is granted with usage permission.
 Changing this parameter will create a new resource.

* `object` - (Required, String, ForceNew) Specifies which object's data usage permissions will be shared.
  Its naming format is as follows:
  + **queues.`queues_name`**: the usage permissions of queue.
  + **databases.`database_name`**: the usage permissions of data in the database.
  + **databases.`database_name`.tables.`table_name`**: the usage permissions of data in the table.
  + **databases.`database_name`.tables.`table_name`.columns.`column_name`**: the usage permissions of data in the column.
  + **jobs.flink.`flink_job_id`**: the usage permissions of data in the flink job.
  + **groups.`package_group_name`**: the usage permissions of data in the package group.
  + **resources.`package_name`**: the usage permissions of data in the package.

  Changing this parameter will create a new resource.

* `privileges` - (Required, List) Specifies the usage permissions of data.
  + **Permissions on Queue, Database and Table**,
   please see [Permissions Management](https://support.huaweicloud.com/intl/en-us/productdesc-dli/dli_07_0006.html)

  + **Permissions on Flink job**. For more details, please see
  [Managing Flink Job Permissions](https://support.huaweicloud.com/intl/en-us/usermanual-dli/dli_01_0479.html) :
      * **GET**: This permission allows user to view the job details.
      * **UPDATE**: This permission allows user to modify the job.
      * **DELETE**: This permission allows user to delete the job.
      * **START**: This permission allows user to start the job.
      * **STOP**: This permission allows user to stop the job.
      * **EXPORT**: This permission allows user to export the job.
      * **GRANT_PRIVILEGE**: This permission allows user to grant job permissions to other users.
      * **REVOKE_PRIVILEGE**: This permission allows user to revoke the job permissions that other users have but cannot
      revoke the job creator's permissions.
      * **SHOW_PRIVILEGES**: This permission allows user to view the job permissions of other users.

  + **Permissions on Package Groups and Packages**. For more details, please see [Managing Permissions on Packages and
   Package Groups](https://support.huaweicloud.com/intl/en-us/usermanual-dli/dli_01_0477.html) :
      * **USE_GROUP**: This permission allows user to use the package of this group.
      * **UPDATE_GROUP**: This permission allows user to update the packages in the group, including creating a package
        in the group.
      * **GET_GROUP**: This permission allows user to query the details of a package in a group.
      * **DELETE_GROUP**: This permission allows user to delete the package of the group.
      * **GRANT_PRIVILEGE**: This permission allows user to grant group permissions to other users.
      * **REVOKE_PRIVILEGE**: This permission allows user to revoke the group permissions that other users have but
        cannot revoke the group owner's permissions.
      * **SHOW_PRIVILEGES**: This permission allows user to view the group permissions of other users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **object/user_name**. It is composed of `object` and `user_name`,
 separated by a slash.

* `is_admin` - Whether this user is an administrator.

## Import

The permission can be imported by `id`, it is composed of `object` and `user_name`, separated by a slash. e.g.:

```bash
terraform import huaweicloud_dli_permission.test databases.database_name/user_name
```
