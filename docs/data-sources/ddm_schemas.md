---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_schemas"
description: ""
---

# huaweicloud_ddm_schemas

Use this data source to get the list of DDM schemas.

## Example Usage

```hcl
variable "ddm_instance_id" {}
variable "ddm_schema_name" {}

data "huaweicloud_ddm_schemas" "test" {
  instance_id = var.ddm_instance_id
  name        = var.ddm_schema_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of DDM instance.

* `name` - (Optional, String) Specifies the name of the DDM schema.

* `status` - (Optional, String) Specifies the status of the DDM schema.

* `shard_mode` - (Optional, String) Specifies the sharding mode of the schema. Values option: **cluster**, **single**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schemas` - Indicates the list of DDM schema.
  The [Schema](#DdmSchemas_Schema) structure is documented below.

<a name="DdmSchemas_Schema"></a>
The `Schema` block supports:

* `name` - Indicates the name of the DDM schema.

* `status` - Indicates the status of the DDM schema.

* `shard_mode` - Indicates the sharding mode of the schema.

* `shard_number` - Indicates the number of shards in the same working mode.

* `data_nodes` - Indicates the RDS instances associated with the schema.
  The [DataNode](#DdmSchemas_SchemaDataNode) structure is documented below.

<a name="DdmSchemas_SchemaDataNode"></a>
The `DataNode` block supports:

* `id` - Indicates the node ID of the associated RDS instance.

* `name` - Indicates the name of the associated RDS instance.

* `status` - Indicates the status of the associated RDS instance
