---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_rts_top_sql_statistics"
description: |-
  Use this data source to query the real-time session Top SQL statistics of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_rts_top_sql_statistics

Use this data source to query the real-time session Top SQL statistics of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_rts_top_sql_statistics" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the session Top SQL statistics.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `top_sql_info` - The list of Top SQL statistics.
  The [top_sql_info](#session_top_sqls_attr) structure is documented below.

<a name="session_top_sqls_attr"></a>
The `top_sql_info` block supports:

* `node_name` - The node name.

* `unique_sql_id` - The normalized SQL ID.

* `query` - The query statement.

* `count` - The SQL execution count.
