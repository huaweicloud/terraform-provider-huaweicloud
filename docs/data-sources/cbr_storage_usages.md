---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_storage_usages"
description: |-
  Use this data source to query the storage usage statistics of CBR within HuaweiCloud.
---

# huaweicloud_cbr_storage_usages

Use this data source to query the storage usage statistics of CBR within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_storage_usages" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_id` - (Optional, String) Specifies the ID of the resource to filter the storage usage.

* `resource_type` - (Optional, String) Specifies the type of the resource to filter by.
  Valid values are **OS::Nova::Server**, **OS::Cinder::Volume**, **OS::Ironic::BareMetalServer**,
  **OS::Native::Server**, **OS::Sfs::Turbo**, **OS::Workspace::DesktopV2**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_usages` - The storage usage statistics list.

  The [storage_usages](#storage_usages_struct) structure is documented below.

<a name="storage_usages_struct"></a>
The `storage_usages` block supports:

* `resource_id` - The ID of the resource.

* `resource_name` - The name of the resource.

* `resource_type` - The type of the resource.

* `backup_count` - The number of backups.

* `backup_size` - The size of backups in bytes.

* `backup_size_multiaz` - The size of multi-AZ backups in bytes.
