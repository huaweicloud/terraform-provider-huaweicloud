---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_regions"
description: ""
---

# huaweicloud_rms_regions

Use this data source to get the list of RMS regions.

## Example Usage

```hcl
data "huaweicloud_rms_regions" "regions" {}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) Specifies the region ID.

* `display_name` - (Optional, String) Specifies the region dispaly name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `regions` - The region list.

  The [regions](#regions_struct) structure is documented below.

<a name="regions_struct"></a>
The `regions` block supports:

* `region_id` - The region ID.

* `display_name` - The display name of the region.
