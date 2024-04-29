---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection_health_checks"
description: ""
---

# huaweicloud_vpn_connection_health_checks

Manages a VPN connection health checks data source within HuaweiCloud.

## Example Usage

```hcl
variable "connection_id" {}
variable "status" {}
variable "source_ip" {}
variable "destination_ip" {}

data "huaweicloud_vpn_connection_health_checks" "services" {
  connection_id  = var.connection_id
  status         = var.status
  source_ip      = var.source_ip
  destination_ip = var.destination_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the VPN connection health check.
  If omitted, the provider-level region will be used.

* `connection_id` - (Optional, String) Specifies the ID of the VPN connection.

* `status` - (Optional, String) Specifies the status of the VPN connection health check.

* `source_ip` - (Optional, String) Specifies the source IP of the VPN connection health check.

* `destination_ip` - (Optional, String) Specifies the destination IP of the VPN connection health check.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connection_health_checks` - All resource connection health checks that match the filter parameters.
  The [connection_health_checks](#connection_health_checks) structure is documented below.

<a name="connection_health_checks"></a>
The `connection_health_checks` block supports:

* `id` - The ID of the connection health check.

* `status` - The status of the connection health check.

* `connection_id` - The connection ID of the connection health check.

* `type` - The type of the connection health check.

* `source_ip` - The source IP address of the VPN connection.

* `destination_ip` - The destination IP address of the VPN connection.

* `proto_type` - The proto type of the connection health check.
