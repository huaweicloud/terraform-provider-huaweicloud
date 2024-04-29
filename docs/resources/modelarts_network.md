---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_network"
description: ""
---

# huaweicloud_modelarts_network

Manages a Modelarts network resource within HuaweiCloud.  
A maximum of 15 networks can be created.

## Example Usage

```hcl
variable "cidr" {}

resource "huaweicloud_modelarts_network" "test" {
  name = "demo"
  cidr = var.cidr
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of network.  
  The name can contain `4` to `32` characters, only lowercase letters, digits and hyphens (-) are allowed.
  The name must start with a lowercase letter and end with a lowercase letter or digit.

  Changing this parameter will create a new resource.

* `cidr` - (Required, String, ForceNew) Network CIDR.
  Valid CIDR blocks are 10.0.0.0/8-24, 172.16.0.0/12-24, and 192.168.0.0/16-24.  

  Changing this parameter will create a new resource.

* `workspace_id` - (Optional, String, ForceNew) Workspace ID, which defaults to 0.  

  Changing this parameter will create a new resource.

* `peer_connections` - (Optional, List) List of networks that can be connected in peer mode.
The [peer_connections](#ModelartsNetwork_PeerConnection) structure is documented below.

<a name="ModelartsNetwork_PeerConnection"></a>
The `peer_connections` block supports:

* `vpc_id` - (Required, String) Interconnect VPC ID.  

* `subnet_id` - (Required, String) Interconnect subnet ID.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of network.  
  Value options are as follows:
    + **Creating**: The network is being created.
    + **Active**: The network is available.
    + **Abnormal**: The network is in an error state.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The modelarts network can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_network.test 0ce123456a00f2591fabc00385ff1234
```
