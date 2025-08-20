---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_volume"
description: |-
  Manages an EVS v3 volume resource within HuaweiCloud.
---

# huaweicloud_evsv3_volume

Manages an EVS v3 volume resource within HuaweiCloud.

## Example Usage

```hcl
variable "volume_type" {}
variable "availability_zone" {}
variable "image_id" {}
variable "name" {}
variable "size" {}

resource "huaweicloud_evsv3_volume" "test" {
  volume_type       = var.volume_type
  availability_zone = var.availability_zone
  image_id          = var.image_id
  name              = var.name
  size              = var.size

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted,
  the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_type` - (Required, String, NonUpdatable) Specifies the disk type.  
  The Valid values are as follows:
  + **SATA**: Ordinary IO type.
  + **SAS**: High I/O type.
  + **GPSSD**: General purpose SSD type.
  + **SSD**: Ultra-high I/O type.
  + **ESSD**: Extreme SSD type.
  + **GPSSD2**: General purpose SSD V2 type.
  + **ESSD2**: Extreme SSD V2 type.

  -> 1. When the specified disk type does not exist in the availability zone, creating the disk fails.
  <br/>2. When creating a disk from a snapshot, the `volume_type` field must be consistent with the source disk of the
  snapshot.

* `availability_zone` - (Optional, String, NonUpdatable) Specifies the availability zone for the disk.
  If the specified AZ does not exist or the specified AZ is different from the AZ where the backup is located,
  the disk creation fails.

* `disaster_recovery_azs` - (Optional, List, NonUpdatable) Specifies to create multiple AZs that support disaster
  recovery cloud hard drives.

-> Specifies either `availability_zone` or `disaster_recovery_azs`, not both.

* `description` - (Optional, String) Specifies the disk description. You can enter up to `85` characters.

* `image_id` - (Optional, String, NonUpdatable) Specifies the image ID from which to create the disk.  
  + Only one of `image_id` or `snapshot_id` can be specified, they cannot be used together.
  + Not supported to create a BMS system disk from a BMS image.
  + You can obtain this value by calling the dataSource `huaweicloud_images_image`.

* `metadata` - (Optional, Map, NonUpdatable) Specifies the key-value pair disk metadata. The length of the key and value
  should not exceed `255` bytes.  
  The valid key-value pairs are as follows:
  + **__system__cmkid**: The encryption CMK ID in metadata. This attribute is used together with **__system__encrypted**
    for encryption.
  + **__system__encrypted**: The encryption field in metadata. The value can be `0` (no encryption) or `1` (encryption).
    If this attribute is not specified, the encryption attribute of the disk is the same as that of the data source.
    If the disk is not created from a data source, the disk is not encrypted by default.
  + **hw:passthrough**: If this attribute value is **true**, the disk device type is SCSI, which allows ECS OSs to
    directly access the underlying storage media and supports SCSI reservation commands. If this attribute is set to
    **false**, the disk device type is VBD, which is also the default type. VBD supports only simple SCSI read/write
    commands. If this attribute is not specified, the disk device type is VBD.

  -> 1. You can also enter other key-value pairs according to the requirements for creating a disk.
  <br/>2. There cannot be null value key value pairs in metadata.

* `multiattach` - (Optional, Bool, NonUpdatable) Specifies whether the disk is shareable. Defaults to **false**.

* `name` - (Optional, String) Specifies the disk name. You can enter up to `64` characters.

* `size` - (Optional, Int) Specifies the disk size, in GB.
  For system disk, the valid value ranges from `1` GB to `1,024` GB.
  For data disk, the valid value ranges from `10` GB to `32,768` GB.

  -> There are the following restrictions for configuring this field:
  <br/>1. This parameter is required when creating an empty disk.
  <br/>2. This parameter is required when creating a disk from a snapshot. The disk size must be greater than or equal
  to the snapshot size.
  <br/>3. This parameter is required when creating a disk from an image. The disk size must be greater than or equal to
  the minimum disk capacity required by min_disk in the image attributes.
  <br/>4. This parameter is optional when you create the disk from a backup. If this parameter is not specified, the
  disk size is equal to the backup size.

  -> Editing this field has the following restrictions:
  <br/>1. Shrinking the disk is not supported.
  <br/>2. If the status of the to-be-expanded disk is **available**, there are no restrictions.
  <br/>3. If the status of the to-be-expanded disk is **in-use**, a shared disk cannot be expanded, which means that
  the value of `multiattach` must be **false**.
  <br/>4. If the status of the to-be-expanded disk is **in-use**, the status of the server to which the disk attached
  must be **ACTIVE**, **PAUSED**, **SUSPENDED**, or **SHUTOFF**.
  <br/>5. Please refer to [Notes and Constraints](https://support.huaweicloud.com/intl/en-us/productdesc-evs/evs_01_0085.html)
  to view the limitations of disk capacity expansion.

* `snapshot_id` - (Optional, String, NonUpdatable) Specifies the snapshot ID from which to create the disk.  
  + Only one of `image_id` or `snapshot_id` can be specified, they cannot be used together.

* `iops` - (Optional, Int, NonUpdatable) Specifies the IOPS(Input/Output Operations Per Second) for the volume.
  The field is valid and required when `volume_type` is set to **GPSSD2** or **ESSD2**.

* `throughput` - (Optional, Int, NonUpdatable) Specifies the throughput for the volume. The Unit is MiB/s.
  The field is valid and required when `volume_type` is set to **GPSSD2**, other types cannot be set.
  This field can be changed only when the disk status is Available or In-use.

-> Before configuring the parameters `volume_type`, `iops`, and `throughput`, please refer to
[Disk Types and Performance](https://support.huaweicloud.com/intl/en-us/productdesc-evs/en-us_topic_0014580744.html).

* `dedicated_storage_id` - (Optional, String, NonUpdatable) Specifies the ID of the DSS storage pool accommodating
  the disk.

* `cascade` - (Optional, Bool) Specifies whether to delete all snapshots associated with disk.
  The value can be **true** or **false**. Defaults to **false**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the disk.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `links` - The disk URI.  
  The [links](#links_struct) structure is documented below.

* `status` - The disk status.  
  Please refer to [EVS Disk Status](https://support.huaweicloud.com/intl/en-us/api-evs/evs_04_0040.html).

* `attachments` - The attachment information of the disk.  
  The [attachments](#attachments_struct) structure is documented below.

* `bootable` - Whether the disk is bootable. **true**: The disk is bootable. **false**: The disk is not bootable.

* `created_at` - The time when the disk was created.

* `volume_image_metadata` - The metadata of the disk image.

* `iops_attribute` - The disk IOPS information. This field is returned only when the disk type is **ESSD2** or
  **GPSSD2**.  
  The [iops_attribute](#iops_attribute_struct) structure is documented below.

* `throughput_attribute` - The disk throughput information. This field is returned only when the disk type is
  **GPSSD2**.  
  The [throughput_attribute](#throughput_attribute_struct) structure is documented below.

* `updated_at` - The time when the disk was updated.

* `snapshot_policy_id` - The snapshot policy ID bound to the disk.

<a name="links_struct"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.

<a name="attachments_struct"></a>
The `attachments` block supports:

* `attached_at` - The time when the disk was attached.

* `attachment_id` - The ID of the attachment information.

* `device` - The device name.

* `host_name` - The name of the physical host housing the cloud server to which the disk is attached.

* `id` - The ID of the attached resource.

* `server_id` - The ID of the server to which the disk is attached.

* `volume_id` - The disk ID.

<a name="iops_attribute_struct"></a>
The `iops_attribute` block supports:

* `frozened` - The frozen tag.

* `id` - The ID of the disk IOPS.

* `total_val` - The IOPS.

* `volume_id` - The disk ID.

<a name="throughput_attribute_struct"></a>
The `throughput_attribute` block supports:

* `frozened` - The frozen tag.

* `id` - The throughput ID.

* `total_val` - The throughput.

* `volume_id` - The disk ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

The EVS v3 volume can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evsv3_volume.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `volume_type`, `disaster_recovery_azs`,
`dedicated_storage_id`, and `cascade`. It is generally recommended running terraform plan after importing a disk.
You can then decide if changes should be applied to the disk, or the resource definition should be updated to align
with the disk. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_evsv3_volume" "test" {
    ...

  lifecycle {
    ignore_changes = [
      volume_type, disaster_recovery_azs, dedicated_storage_id, cascade,
    ]
  }
}
```
