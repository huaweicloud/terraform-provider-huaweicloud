---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_storage_types"
description: |-
  Use this data source to get the list of RDS storage types.
---

# huaweicloud_rds_storage_types

Use this data source to get the list of RDS storage types.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_storage_types" "test" {
  db_type    = "MySQL"
  db_version = "8.0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `db_type` - (Required, String) DB engine. The valid values are **MySQL**, **PostgreSQL**, **SQLServer**, **MariaDB**.

* `db_version` - (Required, String) DB version number.

* `instance_mode` - (Optional, String) HA mode. The valid values are **single**, **ha**, **replica**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_types` - Storage type list. For details, see Data structure of the storage_type field.
  The [storage_type](#Storagetype_storageType) structure is documented below.

<a name="Storagetype_storageType"></a>
The `storage_type` block supports:

* `name` - Storage type.  
  The options are as follows:
    - **ULTRAHIGH**: SSD storage.
    - **LOCALSSD**: Local SSD storage.
    - **CLOUDSSD**: Cloud SSD storage.
        This storage type is supported only with general-purpose and dedicated DB instances.
    - **ESSD**: extreme SSD storage.
        This storage type is supported only with dedicated DB instances.

* `az_status` - The status details of the AZs to which the specification belongs.
  Key indicates the AZ ID, and value indicates the specification status in the AZ.
  The options of value are as follows:
    - **normal**: The specifications in the AZ are available.
    - **unsupported**: The specifications are not supported by the AZ.
    - **sellout**: The specifications in the AZ are sold out.

* `support_compute_group_type` - Performance specifications.
  The options are as follows:
    - **normal**: General-enhanced.
    - **normal2**: General-enhanced II.
    - **armFlavors**: Kunpeng general-enhanced.
    - **dedicicatenormal**: Exclusive x86.
    - **armlocalssd**: Standard Kunpeng.
    - **normallocalssd**: Standard x86.
    - **general**: General-purpose.
    - **dedicated**: Dedicated, which is only supported for cloud SSDs.
    - **rapid**: Dedicated, which is only supported for extreme SSDs.
    - **bigmen**: Large-memory.
