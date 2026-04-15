---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_pemission_set_privileges"
description: |-
  Use this data source to query the privileges of a DataArts Security permission set within HuaweiCloud.
---

# huaweicloud_dataarts_security_pemission_set_privileges

Use this data source to query the privileges of a DataArts Security permission set within HuaweiCloud.

## Example Usage

### Query all privileges under a specified permission set

```hcl
variable "workspace_id" {}
variable "permission_set_id" {}

data "huaweicloud_dataarts_security_pemission_set_privileges" "test" {
  workspace_id      = var.workspace_id
  permission_set_id = var.permission_set_id
}
```

### Query the privileges under a specified permission set and using the permission type filter

```hcl
variable "workspace_id" {}
variable "permission_set_id" {}

data "huaweicloud_dataarts_security_pemission_set_privileges" "test" {
  workspace_id      = var.workspace_id
  permission_set_id = var.permission_set_id
  privilege_type    = "ALLOW"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the permission set privileges are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the permission set belongs.

* `permission_set_id` - (Required, String) Specifies the ID of the permission set to which the granted privileges belong.

* `privilege_type` - (Optional, String) Specifies the privilege type used to filter privileges.

* `privilege_action` - (Optional, String) Specifies the privilege action used to filter privileges.

* `cluster_id` - (Optional, String) Specifies the cluster ID used to filter privileges.

* `cluster_name` - (Optional, String) Specifies the cluster name used to filter privileges.

* `datasource_type` - (Optional, String) Specifies the data source type used to filter privileges.

* `database_name` - (Optional, String) Specifies the database name used to filter privileges.

* `table_name` - (Optional, String) Specifies the table name used to filter privileges.

* `column_name` - (Optional, String) Specifies the column name used to filter privileges.

* `sync_status` - (Optional, String) Specifies the synchronization status used to filter privileges.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `privileges` - The list of privileges that matched filter parameters.  
  The [privileges](#dataarts_security_pemission_set_privileges_attr) structure is documented below.

<a name="dataarts_security_pemission_set_privileges_attr"></a>
The `privileges` block supports:

* `id` - The ID of the granted privilege.

* `permission_set_id` - The ID of the permission set to which the granted privilege belongs.

* `instance_id` - TThe ID of the instance to which the granted privilege belongs.

* `type` - The type of the permission to be configured.

* `actions` - The list of granted privileges.

* `cluster_id` - The ID of the cluster to which the granted privilege belongs.

* `cluster_name` - The name of the cluster to which the granted privilege belongs.

* `datasource_type` - The type of the granted data source.

* `database_name` - The name of the database to which the granted privilege belongs.

* `schema_name` - The name of the schema to which the granted privilege belongs.

* `namespace` - The name of the namespace to which the granted privilege belongs.

* `table_name` - The name of the table to which the granted privilege belongs.

* `column_name` - The name of the column to which the granted privilege belongs.

* `row_level_security` - The row level security of the granted privilege.

* `sync_status` - The synchronization status of the granted privilege.

* `sync_msg` - The synchronization message of the granted privilege.

* `url` - The URL path name of the granted privilege.
