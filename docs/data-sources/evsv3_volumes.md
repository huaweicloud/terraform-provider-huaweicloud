---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_volumes"
description: |-
  Use this data source to get the detailed information list of the EVS v3 disks within HuaweiCloud.
---

# huaweicloud_evsv3_volumes

Use this data source to get the detailed information list of the EVS v3 disks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evsv3_volumes" "test" {
  metadata  = urlencode("{\"hw:passthrough\": \"false\"}")
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name for the disks. This field will undergo a fuzzy matching query, the
  query result is for all disks whose names contain this value.

* `sort_key` - (Optional, String) Specifies the keyword based on which the returned results are sorted.
  The value can be **id**, **status**, **size**, or **created_at**, defaults to **created_at**.

* `sort_dir` - (Optional, String) Specifies the result sorting order.  
  The valid values are as follows:
  + **desc**: The descending order.
  + **asc**: The ascending order.

  Defaults to **desc**.

* `status` - (Optional, String) Specifies the disk status.  
  The valid values are as follows:
  + **creating**
  + **available**
  + **in-use**
  + **error**
  + **attaching**
  + **detaching**
  + **restoring-backup**
  + **backing-up**
  + **error_restoring**
  + **uploading**
  + **downloading**
  + **extending**
  + **error_extending**
  + **deleting**
  + **error_deleting**
  + **rollbacking**
  + **error_rollbacking**
  + **awaiting-transfer**

* `metadata` - (Optional, String) Specifies the disk metadata.
  Please pay attention to escape special characters before use. Please refer to the usage of example.

* `availability_zone` - (Optional, String) Specifies the availability zone for the disks.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `volumes` - The detailed information list of the disks.
  The [volumes](#v3_volumes_struct) structure is documented below.

<a name="v3_volumes_struct"></a>
The `volumes` block supports:

* `id` - The ID of the disk.

* `links` - The disk uri description information.
  The [links](#v3_volumes_links_struct) structure is documented below.

* `name` - The name of the disk.

* `status` - The status of the disk.

* `attachments` - The disk attachment information.
  The [attachments](#v3_volumes_attachments_struct) structure is documented below.

* `availability_zone` - The availability zone of the disk.

* `snapshot_id` - The snapshot ID. This attribute has a value if the disk is created from a snapshot.

* `description` - The description of the disk.

* `bootable` - Is it a boot disk. The valid value is **true** or **false**.

* `created_at` - The time when the disk was created.

* `volume_type` - The disk type.  
  The valid values are as follows:
  + **SATA**: Ordinary IO disk (sold out).
  + **SAS**: High I/O type.
  + **GPSSD**: General purpose SSD type.
  + **SSD**: Ultra-high I/O type.
  + **ESSD**: Extreme SSD type.
  + **GPSSD2**: General purpose SSD V2 type.
  + **ESSD2**: Extreme SSD V2 type.

* `metadata` - The key-value pair disk metadata.  
  The valid key-value pairs are as follows:
  + **__system__cmkid**: The encryption CMK ID in metadata. When used in conjunction with `__system__encrypted`,
    it indicates the need for encryption. Fixed length of `36` bytes.
  + **__system__encrypted**: The encryption field in metadata. The value can be `0` (no encryption) or `1` (encryption).
    When this field does not exist, the disk defaults to not encrypting.
  + **full_clone**: When creating a disk from a snapshot, if you need to use link cloning, please specify a value of `0`
    for this field.
  + **hw:passthrough**: If this attribute value is **true**, the disk device type is SCSI, which allows ECS operating
    system to directly access the underlying storage media and supports SCSI reservation commands. If this attribute is
    set to **false**, the disk device type is VBD, which is also the default type. VBD supports only simple SCSI
    read/write commands. If this attribute is not specified, the disk device type is VBD.

* `size` - The disk size, in GiB.

* `updated_at` - The time when the disk was updated.

* `multiattach` - Is it a shareable disk.

* `volume_image_metadata` - The metadata of disk image.

* `iops` - The iops information of the disk is only returned when the disk type is **ESSD2** or **GPSSD2**.
  The [iops](#v3_volumes_iops_struct) structure is documented below.

* `throughput` - The throughput information of the disk is only returned when the disk type is **GPSD2**.
  The [throughput](#v3_volumes_throughput_struct) structure is documented below.

* `snapshot_policy_id` - The snapshot policy ID bound to the disk.

<a name="v3_volumes_links_struct"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.

<a name="v3_volumes_attachments_struct"></a>
The `attachments` block supports:

* `attached_at` - The time when the disk was attached.

* `attachment_id` - The ID corresponding to the attachment information.

* `device` - The device name to which the disk is attached.

* `host_name` - The name of the physical host housing the cloud server to which the disk is attached.

* `id` - The ID of the attached resource.

* `server_id` - The ID of the server to which the disk is attached.

* `volume_id` - The disk ID.

<a name="v3_volumes_iops_struct"></a>
The `iops` block supports:

* `frozened` - The frozen tag.

* `id` - The ID of the disk IOPS.

* `total_val` - The IOPS size.

* `volume_id` - The disk ID.

<a name="v3_volumes_throughput_struct"></a>
The `throughput` block supports:

* `frozened` - The frozen tag.

* `id` - The disk throughput identification.

* `total_val` - The throughput size.

* `volume_id` - The disk ID.
