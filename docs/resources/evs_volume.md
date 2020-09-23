---
subcategory: "Elastic Volume Service (EVS)"
---

# huaweicloud\_evs\_volume

Manages a volume resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_evs_volume" "volume" {
  name              = "volume"
  description       = "my volume"
  volume_type       = "SATA"
  size              = 20
  availability_zone = "cn-north-4a"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Example Usage with KMS encryption

```hcl
resource "huaweicloud_evs_volume" "volume" {
  name              = "volume"
  description       = "my volume"
  volume_type       = "SATA"
  size              = 20
  kms_id            = var.kms_id
  availability_zone = "cn-north-4a"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) The availability zone for the volume.
    Changing this creates a new volume.

* `volume_type` - (Required) Specifies the disk type.
    Currently, the value can be SSD, SAS, or SATA.
    - SSD: specifies the ultra-high I/O disk type.
    - SAS: specifies the high I/O disk type.
    - SATA: specifies the common I/O disk type.
    If the specified disk type is not available in the AZ, the disk will fail to create.

* `name` - (Optional) Specifies the disk name.
    If you create disks one by one, the name value is the disk name. The value can contain a maximum of 255 bytes.
    If you create multiple disks (the count value is greater than 1), the system automatically adds a hyphen followed 
    by a four-digit incremental number, such as -0000, to the end of each disk name. For example, 
    the disk names can be volume-0001 and volume-0002. The value can contain a maximum of 250 bytes.

* `size` - (Optional) Specifies the disk size, in GB. Its value can be as follows:
    - System disk: 1 GB to 1024 GB
    - Data disk: 10 GB to 32768 GB
    This parameter is mandatory when you create an empty disk. You can specify the parameter value as required within the value range.
    This parameter is mandatory when you create the disk from a snapshot. Ensure that the disk size is greater than or equal to the snapshot size.
    This parameter is mandatory when you create the disk from an image. Ensure that the disk size is greater than or equal to 
    the minimum disk capacity required by min_disk in the image attributes.
    This parameter is optional when you create the disk from a backup. If this parameter is not specified, the disk size is equal to the backup size.
    Changing this parameter will update the disk. You can extend the disk by setting this parameter to a new value, which must be between current size
    and the max size(System disk: 1024 GB; Data disk: 32768 GB). Shrinking the disk is not supported.

* `description` - (Optional) Specifies the disk description. The value can contain a maximum of 255 bytes.

* `image_id` - (Optional) The image ID from which to create the volume.
    Changing this creates a new volume.

* `backup_id` - (Optional) The backup ID from which to create the volume.
    Changing this creates a new volume.

* `snapshot_id` - (Optional) The snapshot ID from which to create the volume.
    Changing this creates a new volume.

* `tags` - (Optional) A maximum of 10 tags can be created for a disk.
    Tag keys of a tag must be unique. Deduplication will be performed for duplicate keys. 
    Therefore, only one tag key in the duplicate keys is valid.

    - Tag key: A tag key is a string of no more than 36 characters.
    It consists of letters, digits, underscores (_), hyphens (-), and Unicode characters (\u4E00-\u9FFF).

    - Tag value: A tag value is a string of no more than 43 characters and can be an empty string.
    It consists of letters, digits, underscores (_), periods (.), hyphens (-), and Unicode characters (\u4E00-\u9FFF).
	
* `multiattach` - (Optional, Default:false) Specifies the shared EVS disk information.
    Changing this creates a new volume.

* `kms_id` - (Optional) The Encryption KMS ID to create the volume.
    Changing this creates a new volume.

* `device_type` - (Optional) The device type of volume to create. Valid options are VBD and SCSI.
	Defaults to VBD. Changing this creates a new volume.

## Attributes Reference

The following attributes are exported:

* `availability_zone` - See Argument Reference above.
* `volume_type` - See Argument Reference above.
* `name` - See Argument Reference above.
* `size` - See Argument Reference above.
* `description` - See Argument Reference above.
* `image_id` - See Argument Reference above.
* `backup_id` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `multiattach` - See Argument Reference above.
* `kms_id` - See Argument Reference above.
* `device_type` - See Argument Reference above.
* `attachment` - If a volume is attached to an instance, this attribute will
    display the Attachment ID, Instance ID, and the Device as the Instance
    sees it.
* `wwn` - Specifies the unique identifier used for mounting the EVS disk.

## Import

Volumes can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_evs_volume.volume_1 14a80bc7-c12c-4fe0-a38a-cb77eeac9bd6
```
