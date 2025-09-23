---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_plugin_update"
description: |-
  Manages an RDS plugin update resource within HuaweiCloud.
---

# huaweicloud_rds_pg_plugin_update

Manages an RDS plugin update resource within HuaweiCloud.

-> **NOTE:** Deleting RDS PostgreSQL plugin update modify is not supported. If you destroy a resource of RDS PostgreSQL
  plugin update, it is only removed from the state, but still remains in the cloud. And the instance doesn't return to
  the state before modifying.

## Example Usage

```hcl
variable "instance_id" {}
variable "database_name" {}
variable "extension_name" {}

resource "huaweicloud_rds_pg_plugin_update" "test" {
  instance_id    = var.instance_id
  database_name  = var.database_name
  extension_name = var.extension_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `database_name` - (Required, String, NonUpdatable) Specifies the database name.

* `extension_name` - (Required, String, NonUpdatable) Specifies the extension name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is `<instance_id>/<database_name>/<extension_name>`.
