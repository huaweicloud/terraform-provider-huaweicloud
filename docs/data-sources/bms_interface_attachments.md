---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_interface_attachments"
description: |-
  Use this data source to get the information about NICs bound to a BMS.
---

# huaweicloud_bms_interface_attachments

Use this data source to get the information about NICs bound to a BMS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_bms_interface_attachments" "demo" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the BMS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `interface_attachments` - Indicates the BMS NICs.

  The [interface_attachments](#interface_attachments_struct) structure is documented below.

<a name="interface_attachments_struct"></a>
The `interface_attachments` block supports:

* `pci_address` - Indicates the BDF number of the NIC in Linux Guest OS.

* `port_state` - Indicates the NIC port status.

* `fixed_ips` - Indicates private IP addresses of NICs.

  The [fixed_ips](#interface_attachments_fixed_ips_struct) structure is documented below.

* `net_id` - Indicates the ID of the subnet (network_id) to which the NIC ports belong.

* `port_id` - Indicates the ID of the NIC port.

* `mac_addr` - Indicates the MAC address of the NIC.

* `driver_mode` - Indicates the NIC driver type in Guest OS.

<a name="interface_attachments_fixed_ips_struct"></a>
The `fixed_ips` block supports:

* `subnet_id` - Indicates the ID of the subnet (subnet_id) corresponding to the private IP address of the NIC.

* `ip_address` - Indicates the NIC private IP address.
