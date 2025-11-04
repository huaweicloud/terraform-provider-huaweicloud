---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_attachable_nics"
description: |-
  Use this data source to get the NICs that can be attached to an ECS.
---

# huaweicloud_compute_attachable_nics

Use this data source to get the NICs that can be attached to an ECS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_compute_attachable_nics" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ECS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attachable_quantity` - Indicates the number of NICs that can be attached to an ECS.

  The [attachable_quantity](#attachable_quantity_struct) structure is documented below.

* `interface_attachments` - Indicates the NIC information list.

  The [interface_attachments](#interface_attachments_struct) structure is documented below.

<a name="attachable_quantity_struct"></a>
The `attachable_quantity` block supports:

* `free_efi_nic` - Indicates the remaining number of EFI NICs that can be attached to an ECS.

* `free_scsi` - Indicates the number of SCSI disks that can be attached.

* `free_blk` - Indicates the number of virtio_blk disks that can be attached.

* `free_disk` - Indicates the total number of disks that can be attached.

* `free_nic` - Indicates the total number of NICs that can be attached.

<a name="interface_attachments_struct"></a>
The `interface_attachments` block supports:

* `driver_mode` - Indicates the NIC driver type, which is virtio by default.

* `pci_address` - Indicates the BDF number of the elastic network interface in Linux GuestOS.

* `port_state` - Indicates the NIC port status.

* `fixed_ips` - Indicates the private IP address list for NICs.

  The [fixed_ips](#interface_attachments_fixed_ips_struct) structure is documented below.

* `net_id` - Indicates the network ID (network_id) that the NIC port belongs to.

* `port_id` - Indicates the ID of the NIC port.

* `mac_addr` - Indicates the MAC address of the NIC.

* `delete_on_termination` - Indicates whether to delete a NIC when detaching it.
  The value can be:
  + **true**: Delete the NIC.
  + **false**: Do not delete the NIC.

* `preserve_on_delete` - Indicates whether to retain the NIC when it is deleted.
  The value can be:
  + **true**: Retain the NIC.
  + **false**: Do not retain the NIC.

* `min_rate` - Indicates the minimum NIC bandwidth.

* `multiqueue_num` - Indicates the number of queues.

<a name="interface_attachments_fixed_ips_struct"></a>
The `fixed_ips` block supports:

* `subnet_id` - Indicates the subnet of the NIC private IP address.

* `ip_address` - Indicates the NIC private IP address.
