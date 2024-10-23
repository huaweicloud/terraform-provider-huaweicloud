---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_connection_bandwidth_associate"
description: ""
---

# huaweicloud_cc_central_network_connection_bandwidth_associate

Associate a global connection bandwidth to a central network connection within HuaweiCloud.

## Example Usage

```hcl
variable "central_network_id" {}
variable "connection_id" {}
variable "global_connection_bandwidth_id" {}

resource "huaweicloud_cc_central_network_connection_bandwidth_associate" "test" {
  central_network_id             = var.central_network_id
  connection_id                  = var.connection_id
  global_connection_bandwidth_id = var.global_connection_bandwidth_id
  bandwidth_size                 = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `central_network_id` - (Required, String, ForceNew) The ID of the central network to which the connection belongs.
  Changing this creates a new resource.

* `connection_id` - (Required, String, ForceNew) The ID of the connection.
  Changing this creates a new resource.

* `bandwidth_size` - (Required, Int) The bandwidth size of the connection.

* `global_connection_bandwidth_id` - (Required, String) The ID of the global connection bandwidth.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the central network connection ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

The central network connection bandwidth associate resource can be imported using the `central_network_id`
and `connection_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cc_central_network_connection_bandwidth_associate.test <central_network_id>/<connection_id>
```
