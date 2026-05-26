---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_database_instance_connection"
description: |-
  Manages DAS Database instance connection resource within HuaweiCloud.
---

# huaweicloud_das_database_instance_connection

Manages DAS Database instance connection resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "username" {}
variable "password" {}

resource "huaweicloud_das_database_instance_connection" "test" {
  instance_id      = var.instance_id
  engine_type      = "mysql"
  network_type     = "rds"
  username         = var.username
  password         = var.password
  is_save_password = true
  description      = "Created by terraform script"
  sql_record_flag  = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the database instance connection is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to which 
  the database instance connection belongs.

* `engine_type` - (Required, String, NonUpdatable) Specifies the engine type of the database instance connection.

* `network_type` - (Required, String, NonUpdatable) Specifies the network type of the database instance connection.

* `username` - (Required, String) Specifies the username of the database instance connection.

* `password` - (Required, String) Specifies the password of the database instance connection.

* `is_save_password` - (Required, Bool) Specifies whether to save the password for the database instance connection.

* `node_ids` - (Optional, List) Specifies the unique identifiers of the instance nodes.

* `description` - (Optional, String) Specifies the description of the database instance connection.

* `port` - (Optional, Int) Specifies the port of the database instance connection.

* `database_name` - (Optional, String) Specifies the database name of the database instance connection.

* `sql_record_flag` - (Optional, Bool) Specifies whether SQL recording is enabled for the database instance connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `instance_name` - The instance name of the database instance connection.

* `datastore_version` - The datastore version of the database instance connection.

* `ip_address` - The ip address of the database instance connection.

* `created_at` - The timestamp when the database instance connection was created, in RFC3339 format.

* `status` - The status of the database instance connection.

* `conn_share_type` - The conn share type of the database instance connection.

* `shared_user_name` - The shared user name of the database instance connection.

* `shared_user_id` - The shared user ID of the database instance connection.

* `expired_at` - The timestamp when the database instance connection expires, in RFC3339 format.

## Import

The DAS database instance connection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_das_database_instance_connection.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `node_ids`,`password` and `sql_record_flag`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the imported state. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_das_database_instance_connection" "test" {
  ...

  lifecycle {
    ignore_changes = [
      node_ids,
      password,
      sql_record_flag,
    ]
  }
}
```
