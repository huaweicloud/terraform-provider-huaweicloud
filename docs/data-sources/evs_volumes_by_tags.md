---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volumes_by_tags"
description: |-
  Use this data source to get EVS volumes filtered by tags within HuaweiCloud.
---

# huaweicloud_evs_volumes_by_tags

Use this data source to get EVS volumes filtered by tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_volumes_by_tags" "example" {
  action = "filter"

  tags = [
    {
      key    = "key_string"
      values = ["value_string"]
    }
  ]

  matches = [
    {
      key   = "resource_name"
      value = "shared01"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the action name. Possible value: **filter**.

* `tags` - (Required, List) Specifies the list of included tags.The key-value pairs of tags.
  A tag list can contain a maximum of `10` keys. Tag keys in a tag list must be unique.
  When multiple keys are specified in a tag list, only the disks having all specified keys are queried.
  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the search criteria supported by disks. Tag keys in a tag list must be unique.
  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the resource tag.

* `values` - (Required, List) Specifies the list of values corresponding to the key.
  A tag list can contain a maximum of `10` values. Tag values in a tag list must be unique.
  If the tag value list is empty, disks that contain any key can be queried.
  When there are multiple values and the key requirements are met, disks that have any of the specified values are queried.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Optional, String) Specifies the key of the resource match.
  Supported keys: **resource_name**, **service_type**.

* `value` - (Optional, String) Specifies the value of the resource match.
  The value, which can contain a maximum of `255` characters. If **resource_name** is specified for `key`,
  the tag value uses a fuzzy match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The number of disks that meet the query criteria.

* `resources` - The list of disks that meet the query criteria.
  The [resources](#evs_volumes_resources) structure is documented below.

<a name="evs_volumes_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The resource details.
  The [resource_detail](#evs_volumes_resource_detail) structure is documented below.

* `tags` - The tag list. The [tags](#evs_volumes_tags) structure is documented below.

<a name="evs_volumes_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `value` - The value of the resource tag.

<a name="evs_volumes_resource_detail"></a>
The `resource_detail` block supports:

* `id` - The disk ID.

* `links` - The disk URI. The [links](#evs_volumes_links) structure is documented below.

* `name` - The disk name.

* `status` - The disk status. For details, see "EVS Disk Status" in Elastic Volume Service User Guide.
  [reference](https://support.huaweicloud.com/intl/en-us/api-evs/evs_04_0040.html)

* `attachments` - The disk attachment information.
  The [attachments](#evs_volumes_attachments) structure is documented below.

* `availability_zone` - The AZ to which the disk belongs.

* `snapshot_id` - The snapshot ID. This parameter has a value if the disk is created from a snapshot.

* `description` - The disk description.

* `created_at` - The time when the disk was created.

* `os_vol_tenant_attr_tenant_id` - The ID of the tenant to which the disk belongs.
  The tenant ID is the same as the project ID.

* `volume_image_metadata` - The metadata of the disk image.

* `volume_type` - The disk type. The value can be **SSD** (ultra-high I/O), **SAS** (high I/O), or **SATA** (common I/O).

* `size` - The disk size, in GiB.

* `bootable` - Whether the disk is bootable. **true**: The disk is bootable. **false**: The disk is not bootable.

* `metadata` - The disk metadata map. Possible keys are:
  + **__system__cmkid**: The encryption CMK ID in metadata.
    This parameter is used together with `__system__encrypted` for encryption.
    The length of cmkid is fixed at `36` bytes. For details about how to obtain the key ID, see Querying the Key List.
  + **__system__encrypted**: The encryption field in metadata.
    The value can be `0` (no encryption) or `1` (encryption). If this parameter is not specified,
    the encryption attribute of the disk is the same as that of the data source.
    If the disk is not created from a data source, the disk is not encrypted by default.
  + **hw_passthrough**: If this parameter value is **true**, the disk device type is SCSI,
    which allows ECS OSs to directly access the underlying storage media.
    SCSI reservation commands are supported. If this parameter is set to **false**, the disk device type is **VBD**,
    which is also the default type. VBD supports only simple SCSI read/write commands.
    If this parameter is not specified, the disk device type is **VBD**.
  + **orderID**: The parameter that describes the disk billing mode in metadata. If this parameter has a value,
    the disk is billed on a yearly/monthly basis. If not, the disk is billed on a pay-per-use basis.

* `updated_at` - The time when the disk was updated.

* `service_type` - The service type. Supported values: **EVS**, **DSS**, **DESS**.

* `multiattach` - Whether the disk is shareable.

* `dedicated_storage_id` - The ID of the dedicated storage pool housing the disk.

* `dedicated_storage_name` - The name of the dedicated storage pool housing the disk.

* `tags` - The disk tags. This field has values if the disk has tags. Otherwise, this field is empty.

* `wwn` - The unique identifier used when attaching the disk.

* `enterprise_project_id` - The ID of the enterprise project that the disk has been added to.

<a name="evs_volumes_attachments"></a>
The `attachments` block supports:

* `server_id` - The ID of the server to which the disk is attached.

* `attachment_id` - The attachment ID.

* `attached_at` - The time when the disk was attached.

* `volume_id` - The disk ID.

* `device` - The device name.

* `id` - The ID of the attached disk.

* `host_name` - The name of the physical host housing the cloud server to which the disk is attached.

<a name="evs_volumes_links"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.
