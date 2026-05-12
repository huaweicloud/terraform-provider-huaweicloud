---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_network"
description: ""
---

# huaweicloud_modelarts_network

Manages a Modelarts network resource within HuaweiCloud.  

## Example Usage

```hcl
variable "network_name" {}
variable "peering_vpc_id" {}
variable "peering_subnet_id" {}

resource "huaweicloud_modelarts_network" "test" {
  name = var.network_name
  cidr = "10.168.0.0/16"

  peer_connections {
    vpc_id    = var.peering_vpc_id
    subnet_id = var.peering_subnet_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the network is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the network.  
  The valid value is limited from `4` to `32`, only lowercase letters, digits and hyphens (-) are allowed.
  The name must start with a lowercase letter and end with a lowercase letter or digit.

* `cidr` - (Required, String, NonUpdatable) Specifies the CIDR of the network.  
  It's recommand to configure CIDR block as following ranges:
  + **10.0.0.0/8-24**
  + **172.16.0.0/12-24**
  + **192.168.0.0/16-24**

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the ID of the workspace to which the network belongs.  
  Defaults to **0** (default workspace).  

* `peer_connections` - (Optional, List) Specifies the list of networks that can be connected in peering mode.  
The [peer_connections](#modelarts_network_peer_connections) structure is documented below.

<a name="modelarts_network_peer_connections"></a>
The `peer_connections` block supports:

* `vpc_id` - (Required, String) Specifies the ID of the VPC to which the peering connection belongs.

* `subnet_id` - (Required, String) Specifies the ID of the subnet to which the peering connection belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is **{name}-{project_id}**.

* `status` - The status of the network.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The network can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_network.test <id>
```
