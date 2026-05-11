---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_sql_control_history_rules"
description: |-
  Use this data source to query the historical SQL throttling rules of a TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_sql_control_history_rules

Use this data source to query the historical SQL throttling rules of a TaurusDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_taurusdb_sql_control_history_rules" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

### Filter by SQL Type

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_taurusdb_sql_control_history_rules" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  sql_type    = "SELECT"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the historical SQL throttling rules.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `node_id` - (Required, String) Specifies the node ID.

* `sql_type` - (Optional, String) Specifies the SQL statement type. The valid values are as follows:
  + **SELECT**
  + **UPDATE**
  + **DELETE**
  + **INSERT**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sql_filter_rules` - The list of historical SQL throttling rules.
  The [sql_filter_rules](#sql_control_history_rules_attr) structure is documented below.

<a name="sql_control_history_rules_attr"></a>
The `sql_filter_rules` block supports:

* `instance_id` - The instance ID.

* `node_id` - The node ID.

* `pattern` - The SQL throttling rule.

* `sql_type` - The SQL statement type.

* `max_concurrency` - The maximum number of concurrent SQL statements.

* `create_at` - The time when a SQL throttling rule was created.

* `expire_at` - The time when a SQL throttling rule expires.

* `delete_at` - The time when a SQL throttling rule was deleted.
