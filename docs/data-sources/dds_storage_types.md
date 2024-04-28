---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_storage_types"
description: ""
---

# huaweicloud_dds_storage_types

Use this data source to get the list of DDS storage types.

## Example Usage

```hcl
data "huaweicloud_dds_storage_types" "test" {
  engine_name = "DDS-Community"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_name` - (Optional, String) Specifies the database type. The valid value is **DDS-Community**.
  For details, see [documentation](https://support.huaweicloud.com/api-dds/dds_database_version.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_types` - Indicates the database disk information list.
  The [storage_types](#Dds_storage_types) structure is documented below.

<a name="Dds_storage_types"></a>
The `storage_types` block supports:

* `name` - Indicates the storage type. The values are as follows:
  + **ULTRAHIGH**: SSD storage.
  + **EXTREMEHIGH**: extreme SSD storage.

* `az_status` - The status details of the AZs to which the specification belongs. Key indicates the AZ ID, and value
  indicates the specification status in the AZ. The values are as follows:
  + **normal**: The specifications in the AZ are available.
  + **unsupported**: The specifications are not supported by the AZ.
  + **sellout**: The specifications in the AZ are sold out.
