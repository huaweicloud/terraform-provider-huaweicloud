---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_types"
description: |-
  Use this data source to get the list of available DeH types in an AZ.
---

# huaweicloud_deh_types

Use this data source to get the list of available DeH types in an AZ.

## Example Usage

```hcl
data "huaweicloud_deh_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Required, String) Specifies the availability zone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dedicated_host_types` - Indicates the available DeH types.

  The [dedicated_host_types](#dedicated_host_types_struct) structure is documented below.

<a name="dedicated_host_types_struct"></a>
The `dedicated_host_types` block supports:

* `host_type` - Indicates the DeH type.

* `host_type_name` - Indicates the name of the DeH type.
