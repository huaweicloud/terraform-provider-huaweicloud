---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_volume_attachments"
description: |-
  Use this data source to get the EVS disks attached to a BMS.
---

# huaweicloud_bms_volume_attachments

Use this data source to get the EVS disks attached to a BMS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_bms_volume_attachments" "demo" {
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

* `volume_attachments` - Indicates disks attached to a BMS.

  The [volume_attachments](#volume_attachments_struct) structure is documented below.

<a name="volume_attachments_struct"></a>
The `volume_attachments` block supports:

* `id` - Indicates the ID of the attached resource.

* `server_id` - Indicates the ID of the BMS that disks are attached to.

* `volume_id` - Indicates the ID of the disk attached to the BMS.

* `device` - Indicates the mount directory, for example, dev/sdd.
