---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_node_sessions"
description: |-
  Use this data source to get the list of user session threads on a TaurusDB node within HuaweiCloud.
---

# huaweicloud_taurusdb_node_sessions

Use this data source to get the list of user session threads on a TaurusDB node within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_taurusdb_node_sessions" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to query the resource.
   If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `node_id` - (Required, String) Specifies the ID of the node in the TaurusDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `processes` - The list of user session threads in the node in the TaurusDB instance.

  The [processes](#processes_struct) structure is documented below.

<a name="processes_struct"></a>
The `processes` block supports:

* `id` - The ID of the user session thread.

* `user` - The user who starts the session thread.

* `host` - The host and port that send requests.

* `db` - The name of the database that is being accessed.

* `command` - The command that is being executed.

* `time` - The time in seconds that the user session thread remains in the current state.

* `state` - The status of the SQL statement that is being executed.

* `info` - The additional information, which is usually the statement that is being executed.
