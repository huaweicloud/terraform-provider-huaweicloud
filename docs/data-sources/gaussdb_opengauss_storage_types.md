---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_storage_types"
description: |-
  Use this data source to get the disk types of GaussDB OpenGauss.
---

# huaweicloud_gaussdb_opengauss_storage_types

Use this data source to get the disk types of GaussDB OpenGauss.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_storage_types" "test" {
  version = "8.201"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `version` - (Required, String) Specifies the DB version number.

* `ha_mode` - (Optional, String) Specifies the instance type.
  Value options:
  + **enterprise**: enterprise edition
  + **centralization_standard**: primary/standby, which is case insensitive.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_type` - Indicates the storage type information.

  The [storage_type](#storage_type_struct) structure is documented below.

<a name="storage_type_struct"></a>
The `storage_type` block supports:

* `name` - Indicates the storage type.
  The value can be:
  + **ULTRAHIGH**: SSD storage.
  + **ESSD**: extreme SSD storage.

* `az_status` - Indicates the status details of the AZs to which the specification belongs.
  Key indicates the AZ ID, and value indicates the specification status in the AZ.
  The value can be:
  + **normal**: on sale.
  + **unsupported**: not supported.
  + **sellout**: sold out.

* `support_compute_group_type` - Indicates the performance specifications.
  The value can be:
  + **normal:** dedicated (1:8).
  + **normal2**: dedicated (1:4).
  + **armFlavors**: Kunpeng dedicated (1:8).
  + **armFlavors2**: Kunpeng dedicated (1:4).
  + **armFlavors2Shared**: Kunpeng general computing-plus II (shared).
  + **general**: General-purpose (1:4).
  + **exclusive**: Dedicated (1:4) It is only suitable for primary/standby instances of the basic edition.
  + **armExclusive**: Kunpeng dedicated (1:4) It is only suitable for primary/standby instances of the basic edition.
  + **economical**: Favored (1:4).
  + **economical2**: Favored (1:8).
