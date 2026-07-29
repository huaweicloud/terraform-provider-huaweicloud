---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_batch_add_data_mask"
description: |-
  Use this resource to batch add data mask within HuaweiCloud.
---

# huaweicloud_dsc_batch_add_data_mask

Use this resource to batch add data mask within HuaweiCloud.

-> This resource is a one-time action resource used to batch add DSC data mask.
  Deleting this resource will not restore the masked data or undo the mask action, but will only
  remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_dsc_batch_add_data_mask" "test" {
  mask_strategies {
    name      = "col"
    algorithm = "SHA256"
  }

  data = [
    {
      col = "test1111"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `mask_strategies` - (Required, List, NonUpdatable) Specifies the list of mask strategies, each corresponding
  to a field. The number of mask strategies cannot exceed 100.
  The [mask_strategies](#mask_strategies_struct) structure is documented below.

* `data` - (Required, List, NonUpdatable) Specifies the data list to be masked.

<a name="mask_strategies_struct"></a>
The `mask_strategies` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the field name to be masked.
  The maximum length is 256 characters.

* `algorithm` - (Required, String, NonUpdatable) Specifies the masking algorithm name.

* `parameters` - (Optional, Map, NonUpdatable) Specifies the masking algorithm parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `masked_data` - The masked data in JSON format.
