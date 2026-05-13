---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_drivers"
description: |-
  Use this data source to get the list of driver files for DRS within HuaweiCloud.
---

# huaweicloud_drs_drivers

Use this data source to get the list of driver files for DRS within HuaweiCloud.

## Example Usage

```hcl
variable "driver_type" {}

data "huaweicloud_drs_drivers" "test" {
  driver_type = var.driver_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `driver_type` - (Required, String) Specifies the type of the driver file to be queried.
  The valid values are as follows:
  + **db2**: DB2 for LUW.
  + **informix**: Informix.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of driver file details.

  The [items](#items) structure is documented below.

<a name="items"></a>
The `items` block supports:

* `driver_name` - The name of the driver file.

* `last_modified` - The last modification time of the driver file.

* `size` - The size of the driver file, in bytes.
