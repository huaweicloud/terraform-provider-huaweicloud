---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_database_privilege"
description: ""
---

# huaweicloud_dli_database_privilege

Using this resource to manage the privileges for the DLI database or data table within HuaweiCloud.

## Example Usage

### Grant permissions of database to a specific user

```hcl
variable "authorized_iam_user_name" {}
variable "database_name" {}

resource "huaweicloud_dli_database_privilege" "test" {
  user_name = var.authorized_iam_user_name
  object    = format("databases.%s", var.database_name)

  privileges = [
    "CREATE_TABLE",
    "CREATE_VIEW",
    "DISPLAY_ALL_TABLES",
    "SELECT",
    "DESCRIBE_TABLE",
    "SHOW_CREATE_TABLE",
  ]
}
```

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies the name of the authorized (IAM) user.  
  The valid length is limited from `1` to `32`, only letters, digits, hyphens (-), underscores (_), dots (.) and spaces
  are allowed. The name cannot contain start with a digit.
  Changing this parameter will create a new resource.

* `object` - (Required, String, ForceNew) Specifies the authorization object definition.
  + If it's authorized to the database, the format is **databases.{database_name}**.
  + If it's authorized to the data table, the format is **databases.{database_name}.tables.{data_table_name}**.

  Changing this parameter will create a new resource.

* `privileges` - (Optional, List) Specifies the list of permissions granted to the database or data table.  
  The valid [permissions](#permissions_for_database_and_table) are documented below.
  Currently, only **SELECT** is supported and is the default value.  
  If you want to grant all permissions, please configure **ALL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consisting of `object` and `user_name`, the format is `<object>/<user_name>`.

## Import

The resource can be imported using the `object` and `user_name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dli_database_privilege.test <object>/<user_name>
```

## Appendix

<a name="permissions_for_database_and_table"></a>

| Number | Permissions for database | Permissions for data table |
| ------ | ------------------------ | -------------------------- |
| Non-inherited permissions | DISPLAY_ALL_TABLES<br>DISPLAY_DATABASE<br>DROP_DATABASE<br>CREATE_TABLE<br>CREATE_VIEW<br>EXPLAIN<br>CREATE_ROLE<br>DROP_ROLE<br>SHOW_ROLES<br>GRANT_ROLE<br>REVOKE_ROLE<br>SHOW_USERS<br>CREATE_FUNCTION<br>DROP_FUNCTION<br>SHOW_FUNCTIONS<br>DESCRIBE_FUNCTION<br> | DISPLAY_TABLE<br>SELECT<br>DESCRIBE_TABLE<br>SHOW_CREATE_TABLE<br>DROP_TABLE<br>TRUNCATE_TABLE<br>ALTER_TABLE_RENAME<br>INSERT_INTO_TABLE<br>INSERT_OVERWRITE_TABLE<br>ALTER_TABLE_ADD_COLUMNS<br>SPARK_APP_ACCESS_META<br>SHOW_PRIVILEGES<br>GRANT_PRIVILEGE<br>REVOKE_PRIVILEGE |
| Inherited permissions | SELECT<br>DESCRIBE_TABLE<br>SHOW_CREATE_TABLE<br>DROP_TABLE<br>TRUNCATE_TABLE<br>ALTER_TABLE_RENAME<br>INSERT_INTO_TABLE<br>INSERT_OVERWRITE_TABLE<br>ALTER_TABLE_ADD_COLUMNS<br>SPARK_APP_ACCESS_META<br>SHOW_PRIVILEGES<br>GRANT_PRIVILEGE<br>REVOKE_PRIVILEGE<br>ALTER_TABLE_ADD_PARTITION<br>ALTER_TABLE_DROP_PARTITION<br>ALTER_TABLE_SET_LOCATION<br>ALTER_TABLE_RENAME_PARTITION<br>ALTER_TABLE_RECOVER_PARTITION<br>SHOW_PARTITIONS |  |
