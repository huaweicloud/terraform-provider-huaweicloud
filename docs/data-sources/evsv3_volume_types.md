---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_volume_types"
description: |-
  Use this data source to query the list of EVS v3 volume types within HuaweiCloud.
---

# huaweicloud_evsv3_volume_types

Use this data source to query the list of EVS v3 volume types within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evsv3_volume_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `volume_types` - The list of volume types.  
  The [volume_types](#volume_types_struct) structure is documented below.

<a name="volume_types_struct"></a>
The `volume_types` block supports:

* `id` - The volume type ID.

* `name` - The volume type name.

* `extra_specs` - The volume type flavor.  
  The [extra_specs](#extra_specs_struct) structure is documented below.

* `description` - The volume type description.

<a name="extra_specs_struct"></a>
The `extra_specs` block supports:

* `availability_zones` - The list of availability zones where the volume type is supported.
  Multiple availability zones separated by commas (,).
  If this filed is empty, the volume type is supported all availability zones.

* `sold_out_availability_zones` - The list of availability zones where the volume type has been sold out.
  Multiple availability zones separated by commas (,).
