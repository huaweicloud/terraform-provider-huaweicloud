---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_permission_set_privilege"
description: |-
  Use this resource to assign data source permissions to workspace permission set or permission set within HuaweiCloud.
---

# huaweicloud_dataarts_security_permission_set_privilege

Use this resource to assign data source permissions to workspace permission set or permission set within HuaweiCloud.

-> If you are assigning privileges to permission set, you can only select permission types included in the parent
   permission set.

## Example Usage

```hcl
variable "workspace_id" {}
variable "permission_set_id" {}
variable "permission_actions" {
  type = list(string)
}
variable "database_name" {}
variable "table_name" {}
variable "connection_id" {}

resource "huaweicloud_dataarts_security_permission_set_privilege" "test" {
  workspace_id      = var.workspace_id
  permission_set_id = var.permission_set_id
  datasource_type   = "DLI"
  type              = "ALLOW"
  actions           = var.permission_actions
  cluster_name      = "*"
  database_name     = var.database_name
  table_name        = var.table_name
  connection_id     = var.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the permission set belongs.
  Changing this creates a new resource.

* `permission_set_id` - (Required, String, ForceNew) Specifies the ID of the permission set to be granted.
  Changing this creates a new resource.

* `datasource_type` - (Required, String, ForceNew) Specifies the type of granted data source.
  The valid values are **HIVE**, **DWS** and **DLI**.
  Changing this creates a new resource.
  
* `type` - (Required, String, ForceNew) Specifies the type of permission to be configured.
  Currently, only **ALLOW** is supported.
  Changing this creates a new resource.

* `actions` - (Required, List) Specifies the list of granted permissions. The valid length is limited from `1` to `32`,
  The valid [permissions](#permissions_for_permission_set) are documented below.

* `cluster_name` - (Required, String, ForceNew) Specifies the cluster name corresponding to the granted data source.
  The valid ranges from `1` to `128`.
  Changing this creates a new resource.
  If `datasource_type` is set to `DLI`, the parameter is set to `*`.

* `cluster_id` - (Optional, String, ForceNew) Specifies the cluster ID corresponding to the granted data source.
  It is required when `datasource_type` is `HIVE` or `DWS`.
  Changing this creates a new resource.

* `connection_id` - (Optional, String) Specifies the data connection ID corresponding to the granted data source.

* `database_url` - (Optional, String, ForceNew) Specifies the URL of the database corresponding to the granted data source.
  Changing this creates a new resource.
  This parameter is only valid when `datasource_type` is set to `HIVE`.
  This parameter is conflict with `database_name`, `table_name` and `column_name`.

* `database_name` - (Optional, String, ForceNew) Specifies the name of the database corresponding to the granted data source.
  Changing this creates a new resource.
  It is required when `datasource_type` is `DWS` or `DLI`.

* `table_name` - (Optional, String, ForceNew) Specifies the name of the data table corresponding to the granted data source.
  Changing this creates a new resource.

* `column_name` - (Optional, String, ForceNew) Specifies the name of column corresponding to the granted data source.
  Changing this creates a new resource.

  -> 1. For `database_name`, `table_name` and `column_name` parameters, the valid length is limited from `1` to `128`,
     only letters, digits, underscores (_) and asterisk (*) are allowed.<br/>2. The permissions of databases, tables,
     and columns are managed by layer.
     For example, a user who has been granted database permissions does not have the permissions of tables and columns.
     Table and column permissions must be granted separately.
  
* `schema_name` - (Optional, String, ForceNew) Specifies the schema name corresponding to the DWS data source.
  Changing this creates a new resource.
  This parameter is only valid when `datasource_type` is set to `DWS`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The current synchronization status of the resource.
  The valid values are **UNKNOWN**, **NOT_SYNC**, **SYNC_SUCCESS** and **SYNC_FAIL**.

* `sync_msg` - The status information of the resource.

## Import

The resource can be imported using `workspace_id`, `permission_set_id` and `id`, separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_security_permission_set_privilege.test <workspace_id>/<permission_set_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `connection_id`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_security_permission_set_privilege" "test" {
  ...

  lifecycle {
    ignore_changes = [
      connection_id,
    ]
  }
}
```

## Appendix

<a name="permissions_for_permission_set"></a>

| Type | HIVE | DWS | DLI |
| ---- | ---- | --- | --- |
| Permissions | ALL<br>SELECT<br>UPDATE<br>CREATE<br>DROP<br>ALTER<br>INDEX<br>READ<br>WRITE<br> | ALL<br>SELECT<br>UPDATE<br>DROP<br>ALTER<br>INSERT<br>CREATE_TABLE<br>DELETE<br>CREATE_SCHEMA<br> | SELECT<br>DROP<br>ALTER<br>INSERT<br>CREATE_TABLE |
