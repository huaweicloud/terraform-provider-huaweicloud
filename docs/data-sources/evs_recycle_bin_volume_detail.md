---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_recycle_bin_volume_detail"
description: |-
  Use this data source to get the detail of EVS recycle bin volume within HuaweiCloud.
---

# huaweicloud_evs_recycle_bin_volume_detail

Use this data source to get the detail of EVS recycle bin volume within HuaweiCloud.

## Example Usage

```hcl
variable "volume_id" {}

data "huaweicloud_evs_recycle_bin_volume_detail" "test" {
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `volume_id` - (Required, String) Specifies the disk ID.  
  For its values, can be obtained using `huaweicloud_evs_volumes` dataSource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `volume` - The list of volume detail in the recycle bin.

  The [volume](#volume_struct) structure is documented below.

<a name="volume_struct"></a>
The `volume` block supports:

* `id` - The disk ID.

* `name` - The disk name.

* `description` - The disk description.

* `status` - The disk status.

* `attachments` - The attachment information of the disk.

  The [attachments](#attachments_struct) structure is documented below.

* `multiattach` - Whether the disk is shared.  
  The valid values are as follows:
  + **true**: Indicated as a shared cloud disk.
  + **false**: Indicated as a common cloud disk.

* `size` - The disk size in GiB.

* `metadata` - The key-value pair disk metadata. Valid key-value pairs are as follows:
  + **__system__cmkid**: The encrypted cmkid field in metadata, when combined with `__system__encrypted`, indicates the
    need for encryption. The cmkid length is fixed at `36` bytes.
  + **__system__encrypted**: The field in metadata that represents encryption function, where `0` represents no
    encryption and `1` represents encryption. When this field is not specified, the encryption properties of the cloud
    disk remain consistent with the data source. If the scene is not created from the data source, it is not encrypted
    by default.
  + **hw:passthrough**: The value of **true** indicates that the device type of the cloud disk is SCSI, which allows the
    ECS operating system to directly access the underlying storage medium, supports SCSI lock command.
    The value of **false** indicates that the device type of the cloud disk is VBD (Virtual Block Device), which is the
    default type, VBD can only support simple SCSI read and write commands.
    When this field does not exist, the cloud disk defaults to VBD type.

* `bootable` - Whether the disk is a boot disk.  
  The valid values are as follows:
  + **true**: indicates a boot disk.
  + **false**: indicates a non-boot disk.

* `tags` - The disk tags.

* `availability_zone` - The availability zone to which the disk belongs.

* `created_at` - The time when the disk was created. The time format is UTC YYYY-MM-DDTHH:MM:SS.XXXXXX.

* `service_type` - The service type to which the disk belongs.  
  The valid values are as follows:
  + **EVS**: Elastic Volume Service.
  + **DSS**: Dedicated Storage Service.

* `updated_at` - The time when the disk information was updated. The time format is UTC YYYY-MM-DDTHH:MM:SS.XXXXXX.

* `volume_type` - The disk type. The valid values are **SATA**, **SAS**, **GPSSD**, **SSD**, **ESSD**, **GPSSD2**,
  and **ESSD2**.

* `enterprise_project_id` - The enterprise project ID.

* `plan_delete_at` - The expected time for cleaning up the disk. The time format is UTC YYYY-MM-DDTHH:MM:SS.XXXXXX.

* `pre_deleted_at` - The time when the disk was put into the recycle bin.
  The time format is UTC YYYY-MM-DDTHH:MM:SS.XXXXXX.

* `dedicated_storage_id` - The ID of the dedicated storage pool to which the disk belongs.

* `dedicated_storage_name` - The name of the dedicated storage pool to which the disk belongs.

* `volume_image_metadata` - The image metadata of the disk. Regarding the detailed explanation of this field,
  see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0703.html).

* `wwn` - The unique identifier of the disk.

<a name="attachments_struct"></a>
The `attachments` block supports:

* `attached_at` - The time when the disk was attached. The time format is UTC YYYY-MM-DDTHH:MM:SS.XXXXXX.

* `attachment_id` - The ID corresponding to the attachment information.

* `device` - The attachment point.

* `host_name` - The name of the physical host to which the disk is attached.

* `id` - The ID of the attached resource.

* `server_id` - The ID of the ECS to which the disk is attached.

* `volume_id` - The disk ID.
