---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_idle_sessions_query_and_kill"
description: |-
  Use this resource to query and kill idle sessions of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_idle_sessions_query_and_kill

Use this resource to query and kill idle sessions of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}

resource "huaweicloud_gaussdb_idle_sessions_query_and_kill" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID.
  Only nodes with CN or DN (primary, standby) components are supported.

* `component_id` - (Required, String, NonUpdatable) Specifies the component ID.
  Only components on the specified node are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `success` - Whether the kill session request is successful.
