---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volumes"
description: ""
---

# huaweicloud_evs_volumes

Use this data source to query the detailed information list of the EVS disks within HuaweiCloud.

## Example Usage

```hcl
variable "target_server" {}

data "huaweicloud_evs_volumes" "test" {
  server_id = var.target_server
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the disk list.
  If omitted, the provider-level region will be used.

* `volume_id` - (Optional, String) Specifies the ID for the disk.

* `name` - (Optional, String) Specifies the name for the disks. This field will undergo a fuzzy matching query, the
  query result is for all disks whose names contain this value.

* `volume_type_id` - (Optional, String) Specifies the type ID for the disks.

* `availability_zone` - (Optional, String) Specifies the availability zone for the disks.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID for filtering.

* `shareable` - (Optional, Bool) Specifies whether the disk is shareable.

* `server_id` - (Optional, String) Specifies the server ID to which the disks are attached.

* `status` - (Optional, String) Specifies the disk status. The valid values are as following:
  + **FREEZED**
  + **BIND_ERROR**
  + **BINDING**
  + **PENDING_DELETE**
  + **PENDING_CREATE**
  + **NOTIFYING**
  + **NOTIFY_DELETE**
  + **PENDING_UPDATE**
  + **DOWN**
  + **ACTIVE**
  + **ELB**
  + **ERROR**
  + **VPN**

* `tags` - (Optional, Map) Specifies the included key/value pairs which associated with the desired disk.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A data source ID in hashcode format.

* `volumes` - The detailed information of the disks. Structure is documented below.

The `volumes` block supports:

* `id` - The data source ID of EVS disk, in UUID format.

* `attachments` - The disk attachment information. Structure is documented below.

* `availability_zone` - The availability zone of the disk.

* `bootable` - Whether the disk is bootable.

* `description` - The disk description.

* `volume_type` - The disk type. Valid values are as follows:
  + **SAS**: High I/O type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **ESSD**: Extreme SSD type.
  + **GPSSD2**: General purpose SSD V2 type.
  + **ESSD2**: Extreme SSD V2 type.

* `iops` - the IOPS(Input/Output Operations Per Second) of the volume. Only valid when `volume_type` is **GPSSD2** or
  **ESSD2**.

* `throughput` - The throughput of the system disk. Only valid when `volume_type` is **GPSSD2**.

* `enterprise_project_id` - The ID of the enterprise project associated with the disk.

* `name` - The disk name.

* `service_type` - The service type, such as EVS, DSS or DESS.

* `shareable` - Whether the disk is shareable.

* `size` - The disk size, in GB.

* `status` - The disk status.

* `create_at` - The time when the disk was created.

* `update_at` - The time when the disk was updated.

* `tags` - The disk tags.

* `wwn` - The unique identifier used when attaching the disk.

The `attachments` block supports:

* `id` - The ID of the attached resource in UUID format.

* `attached_at` - The time when the disk was attached.

* `attached_mode` - The ID of the attachment information.

* `device_name` - The device name to which the disk is attached.

* `server_id` - The ID of the server to which the disk is attached.
