---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection_health_check"
description: ""
---

# huaweicloud_vpn_connection_health_check

Manages a VPN connection health check resource within HuaweiCloud.

## Example Usage

```hcl
variable "connection_id" {}

resource "huaweicloud_vpn_connection_health_check" "test" {
  connection_id = var.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connection_id` - (Required, String, ForceNew) Specifies the ID of the VPN connection to monitor.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `source_ip` - The source IP address of the VPN connection.

* `destination_ip` - The destination IP address of the VPN connection.

* `status` - The status of the connection health check.

## Import

The health check can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpn_connection_health_check.test <id>
```
