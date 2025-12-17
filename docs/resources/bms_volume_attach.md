---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_volume_attach"
description: |-
  Manages a BMS volume attach resource within HuaweiCloud.
---

# huaweicloud_bms_volume_attach

Manages a BMS volume attach resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}
variable "volume_id" {}

resource "huaweicloud_bms_volume_attach" "test" {
  server_id = var.server_id
  volume_id = var.volume_id
  device    = "/dev/sdb"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the BMS ID.

* `volume_id` - (Required, String, NonUpdatable) Specifies the ID of the disk to be attached to a BMS.

* `device` - (Optional, String, NonUpdatable) Specifies the mount point, such as **/dev/sda** and **/dev/sdb**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<server_id>/<volume_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

This resource can be imported using the `server_id` and `volume_id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_bms_volume_attach.test <server_id>/<volume_id>
```
