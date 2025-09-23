---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_schema"
description: |-
  Manages a DDM schema resource within HuaweiCloud.
---

# huaweicloud_ddm_schema

Manages a DDM schema resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "rds_instance_id" {}
variable "rds_password" {}

resource "huaweicloud_ddm_schema" "test"{
  instance_id  = var.instance_id
  name         = "test_schema"
  shard_mode   = "single"
  shard_number = 1

  data_nodes {
    id             = var.rds_instance_id
    admin_user     = "root"
    admin_password = var.rds_password
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DDM schema.
  An instance name starts with a letter, consists of `2` to `48` characters, and can contain only lowercase letters,
  digits, and underscores (_). Cannot contain keywords information_schema, mysql, performance_schema, or sys.

  Changing this parameter will create a new resource.

* `shard_mode` - (Required, String, ForceNew) Specifies the sharding mode of the schema. Values option:
  + **cluster**: indicates that the schema is in sharded mode.
  + **single**: indicates that the schema is in non-sharded mode.

  Changing this parameter will create a new resource.

* `shard_number` - (Required, Int, ForceNew) Specifies the number of shards in the same working mode.
  The value must be greater than or equal to the number of associated RDS instances and less than or equal
  to the number of associated instances multiplied by 64.

  Changing this parameter will create a new resource.

* `data_nodes` - (Required, List, ForceNew) Specifies the RDS instances associated with the schema.

  Changing this parameter will create a new resource.
  The [data_nodes](#data_nodes_struct) structure is documented below.

* `delete_rds_data` - (Optional, String) Specifies whether data stored on the associated DB instances is deleted.

<a name="data_nodes_struct"></a>
The `data_nodes` block supports:

* `id` - (Required, String, ForceNew) Specifies the ID of the RDS instance associated with the schema.

* `admin_user` - (Required, String, ForceNew) Specifies the username for logging in to the associated RDS instance.

* `admin_password` - (Required, String, ForceNew) Specifies the password for logging in to the associated RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the schema status.

* `shards` - Indicates the sharding information of the schema.
  The [shards](#shards_struct) structure is documented below.

* `data_nodes` - Indicates the RDS instances associated with the schema.
  The [data_nodes](#data_nodes_struct) structure is documented below.

* `data_vips` - Indicates the IP address and port number for connecting to the schema.

<a name="shards_struct"></a>
The `shards` block supports:

* `db_slot` - Indicates the number of shards.

* `name` - Indicates the shard name.

* `status` - Indicates the shard status.

* `id` - Indicates the ID of the RDS instance where the shard is located.

<a name="data_nodes_struct"></a>
The `data_nodes` block supports:

* `id` - Indicates the ID of the RDS instance associated with the schema.

* `name` - Indicates the name of the associated RDS instance.

* `status` - Indicates the status of the associated RDS instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 30 minutes.

## Import

The DDM schema can be imported using the `<instance_id>/<name>`, e.g.

```bash
$ terraform import huaweicloud_ddm_schema.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `data_nodes/admin_user`,
`data_nodes/admin_password`. It is generally recommended running `terraform plan` after importing a DDM schema. You can
then decide if changes should be applied to the DDM schema, or the resource definition should be updated to align with
the DDM schema. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ddm_schema" "test" {
  ...

  lifecycle {
    ignore_changes = [
      data_nodes.0.admin_user, data_nodes.0.admin_password
    ]
  }
}
```
