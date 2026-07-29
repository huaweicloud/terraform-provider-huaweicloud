---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_device_mask_algorithms"
description: |-
  Use this data source to get the list of DSC mask algorithms for a specific device.
---

# huaweicloud_dsc_device_mask_algorithms

Use this data source to get the list of DSC mask algorithms for a specific device.

## Example Usage

```hcl
variable "device_id" {}

data "huaweicloud_dsc_device_mask_algorithms" "test" {
  device_id = var.device_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `device_id` - (Required, String) Specifies the device ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `mask_algorithms` - The mask algorithm list.

  The [mask_algorithms](#mask_algorithms_struct) structure is documented below.

<a name="mask_algorithms_struct"></a>
The `mask_algorithms` block supports:

* `id` - The algorithm ID.

* `name` - The algorithm name.
