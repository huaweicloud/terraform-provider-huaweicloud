---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_connection"
description: |-
  Manages an ESW connection resource within HuaweiCloud.
---

# huaweicloud_esw_connection

Manages an ESW connection resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vpc_id" {}
variable "virsubnet_id" {}

resource "huaweicloud_esw_connection" "test" {
  instance_id  = var.instance_id
  name         = "test_name"
  vpc_id       = var.vpc_id
  virsubnet_id = var.virsubnet_id
  fixed_ips    = ["192.168.1.80", "192.168.1.100"]

  remote_infos {
    segmentation_id = 9999
    tunnel_ip       = "11.11.11.11"
    tunnel_port     = "8888"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ESW connection resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance.

* `name` - (Required, String) Specifies the name of the connection.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the vpc ID.

* `virsubnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

* `remote_infos` - (Required, List, NonUpdatable) Specifies the remote tunnel infos.
  The [remote_infos](#remote_infos_struct) structure is documented below.

* `fixed_ips` - (Optional, List, NonUpdatable) Specifies the downlink network port primary and standby IPs.

<a name="remote_infos_struct"></a>
The `remote_infos` block supports:

* `segmentation_id` - (Required, Int, NonUpdatable) Specifies the tunnel number for the connection corresponds to the
  VXLAN network identifier (VNI).

* `tunnel_ip` - (Required, String, NonUpdatable) Specifies the remote tunnel IP of the ESW instance.

* `tunnel_port` - (Optional, Int, NonUpdatable) Specifies the remote tunnel port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `remote_infos` - Indicates the remote tunnel infos.
  The [remote_infos](#remote_infos_attribute) structure is documented below.

* `status` - Indicates the status of the connection.

* `created_at` - Indicates the created time of the connection.

* `updated_at` - Indicates the updated time of the connection.

<a name="remote_infos_attribute"></a>
The `remote_infos` block supports:

* `tunnel_type` - Indicates the tunnel type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

This resource can be imported using the `instance_id` and `id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_esw_connection.test <instance_id>/<id>
```
