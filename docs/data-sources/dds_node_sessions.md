---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_node_sessions"
description: |-
  Use this data source to get the list of DDS instance node sessions.
---

# huaweicloud_dds_node_sessions

Use this data source to get the list of DDS instance node sessions.

## Example Usage

```hcl
variable "node_id" {}

data "huaweicloud_dds_node_sessions" "test" {
  node_id = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `node_id` - (Required, String) Specifies the ID of the node.
  For a cluster instance, you can select any mongos, shard, or config node.
  For a replica set instance, you can select the primary or secondary node.

* `plan_summary` - (Optional, String) Specifies the description of an execution plan.
  The valid values are as follows:
  + **COLLSCAN**
  + **IXSCAN**
  + **FETCH**
  + **SORT**
  + **LIMIT**
  + **SKIP**
  + **COUNT**
  + **COUNT_SCAN**
  + **TEXT**
  + **PROJECTION**

* `type` - (Optional, String) Specifies the operation type.
  The valid values are as follows:
  + **none**
  + **update**
  + **insert**
  + **query**
  + **command**
  + **getmore**
  + **remove**
  + **killcursors**

* `namespace` - (Optional, String) Specifies the namespace.

* `cost_time` - (Optional, Int) Specifies the running time, in μs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sessions` - Indicates the list of sessions.
  The [sessions](#sessions_struct) structure is documented below.

* `total_count` - Indicates the total number of the sessions.

<a name="sessions_struct"></a>
The `sessions` block supports:

* `id` - Indicates the session ID.

* `active` - Indicates whether the current session is active.

* `operation` - Indicates the operation.

* `type` - Indicates the operation type.

* `cost_time` - Indicates the running time.

* `plan_summary` - Indicates the description of an execution plan.

* `host` - Indicates the host.

* `client` - Indicates the client address.

* `description` - Indicates the connection description.

* `namespace` - Indicates the namespace.

* `db` - Indicates the name of the database that is being operated.

* `user` - Indicates the user name.
