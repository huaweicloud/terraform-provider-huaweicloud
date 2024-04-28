---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_replication_attach"
description: ""
---

# huaweicloud_sdrs_replication_attach

Manages an SDRS replication attach resource within HuaweiCloud.

-> Please make sure the status of the protection group must be available or protected, the status of the
protected instance must be available or protected and the status of the replication pair must be available or protected.
The non-shared replication pair has not been attached to any protected instance.

## Example Usage

```hcl
variable "protected_instance_id" {}
variable "replication_id" {}
variable "device" {}

resource "huaweicloud_sdrs_replication_attach" "test" {
  instance_id    = var.protected_instance_id
  replication_id = var.replication_id
  device         = var.device
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a protected instance.

  Changing this parameter will create a new resource.

* `replication_id` - (Required, String, ForceNew) Specifies the ID of a replication pair.

  Changing this parameter will create a new resource.

* `device` - (Required, String, ForceNew) Specifies the disk device name of a replication pair. There are several
  restrictions on this field as follows：

  + The new disk device name cannot be the same as an existing one.

  + Set the parameter value to /dev/sda for the system disks of protected instances created using Xen servers and to
  /dev/sdx for data disks, where x is a letter in alphabetical order. For example, if there are two data disks, set the
  device names of the two data disks to /dev/sdb and /dev/sdc, respectively. If you set a device name starting with
  /dev/vd, the system uses /dev/sd by default.

  + Set the parameter value to /dev/vda for the system disks of protected instances created using KVM servers and
  to /dev/vdx for data disks, where x is a letter in alphabetical order. For example, if there are two data disks,
  set the device names of the two data disks to /dev/vdb and /dev/vdc, respectively. If you set a device name starting
  with /dev/sd, the system uses /dev/vd by default.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the SDRS protected instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SDRS replication attach can be imported using the `protected_instance_id` and `replication_id`, separated
by a slash , e.g.

```bash
$ terraform import huaweicloud_sdrs_replication_attach.test <protected_instance_id>/<replication_id>
```
