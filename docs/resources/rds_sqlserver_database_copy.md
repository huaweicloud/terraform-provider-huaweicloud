---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_database_copy"
description: |-
  Manages an RDS SQLServer database copy resource within HuaweiCloud.
---

# huaweicloud_rds_sqlserver_database_copy

Manages an RDS database copy resource within HuaweiCloud.

-> **NOTE:** Deleting RDS SQLServer database copy is not supported. If you destroy a resource of RDS SQLServer database
  copy, it is only removed from the state, but still remains in the cloud. And the instance doesn't return to the state
  before modifying.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_sqlserver_database_copy" "test" {
  instance_id    = var.instance_id
  procedure_name = "copy_database"
  db_name_source = "test_db_source"
  db_name_target = "test_db_target"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `procedure_name` - (Required, String, NonUpdatable) Specifies the operation name. Value options: **copy_database**.

* `db_name_source` - (Required, String, NonUpdatable) Specifies the name of the source database.

* `db_name_target` - (Required, String, NonUpdatable) Specifies the name of the target database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is `<instance_id>/<procedure_name>`.
