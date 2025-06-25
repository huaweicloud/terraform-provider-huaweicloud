---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_restore_read_replica_database"
description: |-
  Manages the operation to restore a database from a read replica to the primary RDS
  PostgreSQL instance within HuaweiCloud.
---

# huaweicloud_rds_restore_read_replica_database

Manages the operation to restore a database from a read replica to the primary RDS
PostgreSQL instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_restore_read_replica_database" "test" {
  instance_id = var.instance_id
  
  databases {
    old_name = "test"
    new_name = "test_terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL read replica instance.

* `databases` - (Required, List, NonUpdatable) Specifies the databases to be restored.
  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `old_name` - (Required, String, NonUpdatable) Specifies the name of the original database to be restored.

* `new_name` - (Required, String, NonUpdatable) Specifies the name of the new database after the restoration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.
