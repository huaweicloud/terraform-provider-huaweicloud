---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_sql_auto_throttling_records"
description: |-
  Use this data source to query the auto throttling records of a TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_sql_auto_throttling_records

Use this data source to query the auto throttling records of a TaurusDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_taurusdb_sql_auto_throttling_records" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the auto throttling records.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `node_id` - (Required, String) Specifies the node ID. The node role must be the primary node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of auto throttling records.
  The [logs](#sql_auto_throttling_logs_attr) structure is documented below.

<a name="sql_auto_throttling_logs_attr"></a>
The `logs` block supports:

* `node_id` - The node ID.

* `pattern` - The SQL throttling rule.

* `sql_type` - The throttling type.

* `max_concurrency` - The maximum number of concurrent requests.

* `create_at` - The timestamp when throttling starts.

* `expire_at` - The timestamp when throttling expires.
