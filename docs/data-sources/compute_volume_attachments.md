---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_volume_attachments"
description: |-
  Use this data source to get the disk attachments of an ECS.
---

# huaweicloud_compute_volume_attachments

Use this data source to get the disk attachments of an ECS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_compute_volume_attachments" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ECS ID in UUID format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `volume_attachments` - Indicates the disks attached to an ECS.

  The [volume_attachments](#volume_attachments_struct) structure is documented below.

<a name="volume_attachments_struct"></a>
The `volume_attachments` block supports:

* `device` - Indicates the drive letter of the EVS disk, displayed as the device name on the console.

* `id` - Indicates the mount ID, which is the same as the EVS disk ID.

* `server_id` - Indicates the ECS ID in UUID format.

* `volume_id` - Indicates the EVS disk ID in UUID format.
