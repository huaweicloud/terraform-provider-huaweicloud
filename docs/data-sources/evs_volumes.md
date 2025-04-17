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
  ids       = urlencode("['XXX','XXX']")
  metadata  = urlencode("{\"hw:passthrough\": \"true\"}")
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

* `dedicated_storage_id` - (Optional, String) Specifies the dedicated storage pool ID. All disks in the dedicated storage
  pool can be filtered by exact match.

* `dedicated_storage_name` - (Optional, String) Specifies the dedicated storage pool name. All disks in the dedicated
  storage pool can be filtered by fuzzy match.

* `ids` - (Optional, String) Specifies the disk IDs. The value is in the ids=['id1','id2',...,'idx'] format.
  In the response, the `ids` value contains valid disk IDs only. Invalid disk IDs are ignored.
  The details about a maximum of `60` disks can be queried. If `volume_id` and `ids` are both specified in the request,
  `volume_id` will be ignored.
  Please pay attention to escape special characters before use. Please refer to the usage of example.

* `metadata` - (Optional, String) Specifies the disk metadata.
  Please pay attention to escape special characters before use. Please refer to the usage of example.

* `service_type` - (Optional, String) Specifies the service type. Supported services are **EVS**, **DSS**, and **DESS**.

* `sort_dir` - (Optional, String) Specifies the result sorting order. The default value is **desc**.
  + **desc**: The descending order.
  + **asc**: The ascending order.

* `sort_key` - (Optional, String) Specifies the keyword based on which the returned results are sorted.
  The value can be **id**, **status**, **size**, or **created_at**, and the default value is **created_at**.

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

* `dedicated_storage_id` - The ID of the dedicated storage pool housing the disk.

* `dedicated_storage_name` - The name of the dedicated storage pool housing the disk.

* `iops_attribute` - The disk IOPS information. This attribute appears only for a general purpose SSD V2 or an extreme
  SSD V2 disk. The [iops_attribute](#iops_attribute_struct) structure is documented below.

* `throughput_attribute` - The disk throughput information. This attribute appears only for a general purpose SSD V2 disk.
  The [throughput_attribute](#throughput_attribute_struct) structure is documented below.

* `links` - The disk URI. The [links](#links_struct) structure is documented below.

* `metadata` - The key-value pair disk metadata. Valid key-value pairs are as follows:
  + **__system__cmkid**: The encryption CMK ID in metadata. This attribute is used together with **__system__encrypted**
    for encryption.
  + **__system__encrypted**: The encryption field in metadata. The value can be `0` (no encryption) or `1` (encryption).
    If this attribute is not specified, the encryption attribute of the disk is the same as that of the data source.
    If the disk is not created from a data source, the disk is not encrypted by default.
  + **full_clone**: The creation method when the disk is created from a snapshot. `0`: linked clone. `1`: full clone.
  + **hw:passthrough**: If this attribute value is **true**, the disk device type is SCSI, which allows ECS OSs to directly
    access the underlying storage media and supports SCSI reservation commands. If this attribute is set to **false**,
    the disk device type is VBD, which is also the default type. VBD supports only simple SCSI read/write commands.
    If this attribute is not specified, the disk device type is VBD.
  + **orderID**: The attribute that describes the disk billing mode in metadata. If this attribute has a value, the disk
    is billed on a yearly/monthly basis. If this attribute is empty, the disk is billed on a pay-per-use basis.

* `serial_number` - The disk serial number. This field is returned only for non-HyperMetro SCSI disks and is used for
  disk mapping in the VM.

* `snapshot_id` - The snapshot ID. This attribute has a value if the disk is created from a snapshot.

* `volume_image_metadata` - The metadata of the disk image.

<a name="links_struct"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.

<a name="throughput_attribute_struct"></a>
The `throughput_attribute` block supports:

* `frozened` - The frozen tag.

* `id` - The throughput ID.

* `total_val` - The throughput.

<a name="iops_attribute_struct"></a>
The `iops_attribute` block supports:

* `frozened` - The frozen tag.

* `id` - The ID of the disk IOPS.

* `total_val` - The IOPS.

The `attachments` block supports:

* `id` - The ID of the attached resource in UUID format.

* `attached_at` - The time when the disk was attached.

* `attached_mode` - The ID of the attachment information.

* `device_name` - The device name to which the disk is attached.

* `server_id` - The ID of the server to which the disk is attached.

* `host_name` - The name of the physical host housing the cloud server to which the disk is attached.

* `attached_volume_id` - The ID of the attached disk.

* `volume_id` - The disk ID.
