---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_node_sessions_kill"
description: |-
  Use this resource to terminate specified user session threads on a TaurusDB node within HuaweiCloud.
---

# huaweicloud_taurusdb_node_sessions_kill

Use this resource to terminate specified user session threads on a TaurusDB node within HuaweiCloud.

-> This resource is only a one-time action resource to kill user session threads. Deleting this resource
will not clear the corresponding instance, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "processes_ids" {
  type = list(number)
}

resource "huaweicloud_taurusdb_node_sessions_kill" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  processes   = var.processes_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the ID of the TaurusDB instance.

* `node_id` - (Required, String, NoneUpdatable) Specifies the ID of the node in the TaurusDB instance.

* `processes` - (Required, List, NoneUpdatable) Specifies the IDs of user session threads to be terminated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `processes_killed` - The IDs of terminated user session threads in requested processes.

* `processes_not_found` - The IDs of user session threads that were not found in requested processes.
