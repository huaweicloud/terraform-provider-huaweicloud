---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_kill_op_rules"
description: |-
  Use this data source to get the list of killOp rules.
---

# huaweicloud_dds_kill_op_rules

Use this data source to get the list of killOp rules.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_kill_op_rules" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.  

* `instance_id` - (Required, String) Specifies the ID of the DDS instance.

* `operation_types` - (Optional, String) Specifies the SQL operation type.  
  The valid values are as follows:
  + **insert**: Indicates operation for inserting data.
  + **update**: Indicates operation for updating data.
  + **query**: Indicates operation for querying data.
  + **command**: Indicates command operation.
  + **remove**: Indicates operation for deleting data.
  + **getmore**: Indicates operation for obtaining more data.

* `namespaces` - (Optional, String) Specifies the namespace of a table.  
  The value format is database_name or database_name.table_name.
  This parameter can be left blank, indicating that this rule has no restrictions on table namespaces.
  If this parameter is set to a database name, this rule applies to operations on all collections in the database.
  If this parameter is set to a value in the format of database_name.collection_name, this rule only applies to
  operations on the collection.

* `status` - (Optional, String) Specifies the status of killOp rule.  
  The valid values are as follows:
  + **ENABLED**
  + **DISABLED**

* `plan_summary` - (Optional, String) Specifies the execution plan.  
  The valid values are as follows:
  + **COLLSCAN**
  + **SORT_KEY_GENERATOR**
  + **SKIP**
  + **LIMIT**
  + **GEO_NEAR_2DSPHERE**
  + **GEO_NEAR_2D**
  + **AGGREGATE**
  + **OR**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - Indicates the list of kill Op rules.
  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - Indicates the killOP rule ID.

* `operation_types` - Indicates the SQL operation type.

* `status` - Indicates the killOp rule status.

* `namespaces` - Indicates the namespace.

* `client_ips` - Indicates the client IP address.

* `plan_summary` - Indicates the execution plan.

* `node_type` - Indicates the node type.
  + **mongos_shard**: Indicates that this rule applies to both mongos and shard nodes.
  + **mongos**: Indicates that this rule only applies to the mongos node in a cluster.
  + **shard**: Indicates that this rule only applies to the shard node in a cluster.
  + **replica**: Indicates that this rule applies to replica sets.

* `max_concurrency` - Indicates the maximum number of concurrent SQL statements.

* `secs_running` - Indicates the maximum execution duration of a single SQL statement.
