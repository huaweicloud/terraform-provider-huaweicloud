---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_volume_type_detail"
description: |-
  Use this data source to get the list of EVS volume type detail (V3) within HuaweiCloud.
---

# huaweicloud_evsv3_volume_type_detail

Use this data source to get the list of EVS volume type detail (V3) within HuaweiCloud.

## Example Usage

```hcl
variable "type_id" {}

data "huaweicloud_evsv3_volume_type_detail" "test" {
  type_id = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type_id` - (Required, String) Specifies the disk type ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `volume_type` - The returned disk type.
  The [volume_type](#volume_type_struct) structure is documented below.

<a name="volume_type_struct"></a>
The `volume_type` block supports:

* `id` - The disk type ID.

* `name` - The disk type name.

* `extra_specs` - The disk type flavor.
  The [extra_specs](#extra_specs_struct) structure is documented below.

* `description` - The disk type description.

<a name="extra_specs_struct"></a>
The `extra_specs` block supports:

* `reskey_availability_zones` - The list of AZs where the disk type is supported. Elements in the list are separated
    by commas (,). If this parameter is not specified, the disk type is supported in all AZs.

* `os_vendor_extended_sold_out_availability_zones` - The list of AZs where the disk type has been sold out. Elements
    in the list are separated by commas (,).
