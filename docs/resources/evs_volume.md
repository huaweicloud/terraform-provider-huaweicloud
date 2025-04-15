---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume"
description: |-
  Manages a volume resource within HuaweiCloud.
---

# huaweicloud_evs_volume

Manages a volume resource within HuaweiCloud.

## Example Usage

```hcl
variable "availability_zone" {}

resource "huaweicloud_evs_volume" "volume" {
  name              = "volume"
  description       = "my volume"
  volume_type       = "SAS"
  size              = 20
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Example Usage with KMS encryption

```hcl
variable "availability_zone" {}

resource "huaweicloud_evs_volume" "volume" {
  name              = "volume"
  description       = "my volume"
  volume_type       = "SAS"
  size              = 20
  kms_id            = var.kms_id
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Example Usage with server_id

```hcl
variable "image_id" {}
variable "flavor_id" {}
variable "key_pair" {}
variable "security_group_id" {}
variable "availability_zone" {}
variable "subnet_id" {}

resource "huaweicloud_compute_instance" "myinstance" {
  name               = "instance"
  image_id           = var.image_id
  flavor_id          = var.flavor_id
  key_pair           = var.key_pair
  security_group_ids = [var.security_group_id]
  availability_zone  = var.availability_zone

  network {
    uuid = var.subnet_id
  }
}

resource "huaweicloud_evs_volume" "volume" {
  name              = "volume"
  description       = "my volume"
  volume_type       = "SAS"
  size              = 20
  availability_zone = var.availability_zone
  server_id         = huaweicloud_compute_instance.myinstance.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the disk. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone for the disk.

  Changing this parameter will create a new resource.

* `volume_type` - (Required, String) Specifies the disk type. Valid values are as follows:
  + **SAS**: High I/O type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **ESSD**: Extreme SSD type.
  + **GPSSD2**: General purpose SSD V2 type.
  + **ESSD2**: Extreme SSD V2 type.

  If the specified disk type is not available in the AZ, the disk will fail to create. The volume type **ESSD2** only
  support in postpaid charging mode. When creating a cloud disk from a snapshot, the `volume_type` field must be
  consistent with the snapshot source cloud disk.

  -> There are some restrictions on changing the cloud disk type:
  <br/>1. Changing the cloud disk type is currently in the public beta stage. Please submit a work order in advance
  to apply for the public beta.
  <br/>2. The cloud disk type can be changed only when the disk status is Available or In-use.
  <br/>3. Changing the disk type may take several hours or even longer, and cannot be stopped.
  **It is strongly recommended that users proactively configure a reasonable change timeout before changing the disk type.**
  <br/>4. Refer to [Changing the EVS Disk Type](https://support.huaweicloud.com/intl/en-us/usermanual-evs/evs_01_0062.html)
  for more details.

* `iops` - (Optional, Int) Specifies the IOPS(Input/Output Operations Per Second) for the volume.
  The field is valid and required when `volume_type` is set to **GPSSD2** or **ESSD2**.
  This field can be changed only when the disk status is Available or In-use.

* `throughput` - (Optional, Int) Specifies the throughput for the volume. The Unit is MiB/s.
  The field is valid and required when `volume_type` is set to **GPSSD2**.
  This field can be changed only when the disk status is Available or In-use.

-> Before configuring the parameters `volume_type`, `iops`, and `throughput`, please refer to
[Disk Types and Performance](https://support.huaweicloud.com/intl/en-us/productdesc-evs/en-us_topic_0014580744.html).

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

* `description` - (Optional, String) Specifies the disk description. You can enter up to `85` characters.

* `image_id` - (Optional, String, ForceNew) Specifies the image ID from which to create the disk.

  Changing this parameter will create a new resource.

* `backup_id` - (Optional, String, ForceNew) Specifies the backup ID from which to create the disk.

  Changing this parameter will create a new resource.

* `snapshot_id` - (Optional, String, ForceNew) Specifies the snapshot ID from which to create the disk.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the disk.

* `multiattach` - (Optional, Bool, ForceNew) Specifies whether the disk is shareable. Defaults to **false**.

  Changing this parameter will create a new resource.

* `kms_id` - (Optional, String, ForceNew) Specifies the Encryption KMS ID to create the disk.

  Changing this parameter will create a new resource.

* `device_type` - (Optional, String, ForceNew) Specifies the device type of disk to create. Valid options are **VBD** and
  **SCSI**. Defaults to **VBD**.

  Changing this parameter will create a new resource.

* `dedicated_storage_id` - (Optional, String, ForceNew) Specifies the ID of the DSS storage pool accommodating the disk.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the disk.
  For enterprise users, if omitted, default enterprise project will be used.

  Changing this parameter will create a new resource.

* `cascade` - (Optional, Bool) Specifies the delete mode of snapshot. The default value is **false**. All snapshot
  associated with the disk will also be deleted when the parameter is set to **true**.

  -> This parameter is only valid for pay-as-you-go resources, and the snapshots bound to the package period resources
  will be removed while resources unsubscribed.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the disk.
  The valid values are as follows:
  + **prePaid**: the yearly/monthly billing mode.
  + **postPaid**: the pay-per-use billing mode.

  Defaults to **postPaid**. Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the disk.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the disk.
  + If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  + If `period_unit` is set to **year**, the valid value is `1`.

  This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

* `server_id` - (Optional, String, ForceNew) Specifies the server ID to which the cloud volume is to be mounted.
  After specifying the value of this field, the cloud volume will be automatically attached on the cloud server.
  The charging_mode of the created cloud volume will be consistent with that of the cloud server.
  Currently, only ECS cloud-servers are supported, and BMS bare metal cloud-servers are not supported yet.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `attachment` - If a disk is attached to an instance, this attribute will display the attachment ID, instance ID, and
  the device as the instance sees it. The [attachment](#attachment_struct) structure is documented below.

* `wwn` - The unique identifier used for mounting the EVS disk.

* `dedicated_storage_name` - The name of the DSS storage pool accommodating the disk.

* `status` - The disk status.
  Please refer to [EVS Disk Status](https://support.huaweicloud.com/intl/en-us/api-evs/evs_04_0040.html).

* `bootable` - Whether the disk is bootable. **true**: The disk is bootable. **false**: The disk is not bootable.

* `created_at` - The time when the disk was created.

* `updated_at` - The time when the disk was updated.

* `iops_attribute` - The disk IOPS information. This attribute appears only for a general purpose SSD V2 or an extreme
  SSD V2 disk. The [iops_attribute](#iops_attribute_struct) structure is documented below.

* `throughput_attribute` - The disk throughput information. This attribute appears only for a general purpose SSD V2 disk.
  The [throughput_attribute](#throughput_attribute_struct) structure is documented below.

* `links` - The disk URI.
  The [links](#links_struct) structure is documented below.

* `all_metadata` - The key-value pair disk metadata. Valid key-value pairs are as follows:
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

* `service_type` - The service type. Supported services are **EVS**, **DSS**, and **DESS**.

* `all_volume_image_metadata` - The metadata of the disk image.

<a name="links_struct"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.

<a name="iops_attribute_struct"></a>
The `iops_attribute` block supports:

* `frozened` - The frozen tag.

* `id` - The ID of the disk IOPS.

* `total_val` - The IOPS.

<a name="throughput_attribute_struct"></a>
The `throughput_attribute` block supports:

* `frozened` - The frozen tag.

* `id` - The throughput ID.

* `total_val` - The throughput.

<a name="attachment_struct"></a>
The `attachment` block supports:

* `id` - The ID of the attachment information.

* `instance_id` - The ID of the server to which the disk is attached.

* `device` - The device name.

* `attached_at` - The time when the disk was attached.

* `host_name` - The name of the physical host housing the cloud server to which the disk is attached.

* `attached_volume_id` - The ID of the attached disk.

* `volume_id` - The disk ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 180 minutes.
* `delete` - Default is 3 minutes.

## Import

Volumes can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evs_volume.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `cascade`, `period_unit`, `period`,
`server_id`, `auto_renew`, and `charging_mode`. It is generally recommended running terraform plan after importing a disk.
You can then decide if changes should be applied to the disk, or the resource definition should be updated to align
with the disk. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_evs_volume" "test" {
    ...

  lifecycle {
    ignore_changes = [
      cascade, period_unit, period, server_id, auto_renew, charging_mode, 
    ]
  }
}
```
