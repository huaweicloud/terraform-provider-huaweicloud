---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_servergroup"
sidebar_current: "docs-huaweicloud-resource-compute-servergroup"
description: |-
  Manages a V2 Server Group resource within HuaweiCloud.
---

# huaweicloud\_compute\_servergroup

Manages Server Group resource within HuaweiCloud.
This is an alternative to `huaweicloud_compute_servergroup_v2`

## Example Usage

```hcl
resource "huaweicloud_compute_servergroup" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the server group. Changing this creates
    a new server group.

* `policies` - (Required) The set of policies for the server group. Only two
    policies are available right now, and both are mutually exclusive. Possible values are "affinity" and "anti-affinity". 
    "affinity": All instances/servers launched in this group will be hosted on the same compute node.
    "anti-affinity": All instances/servers launched in this group will be hosted on different compute nodes.
    Changing this creates a new server group.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `policies` - See Argument Reference above.
* `members` - The instances that are part of this server group.

## Import

Server Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_compute_servergroup.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
