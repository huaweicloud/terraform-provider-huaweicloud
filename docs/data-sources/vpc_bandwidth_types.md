---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_bandwidth_types"
description: |-
  Use this data source to get a list of share bandwidth types.
---

# huaweicloud_vpc_bandwidth_types

Use this data source to get a list of share bandwidth types.

## Example Usage

```hcl
data "huaweicloud_vpc_bandwidth_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `bandwidth_type` - (Optional, String) Specifies the bandwidth type.

* `name_en` - (Optional, String) Specifies the English description of the bandwidth type.

* `name_zh` - (Optional, String) Specifies the Chinese description of the bandwidth type.

* `public_border_group` - (Optional, String) Specifies the location of the bandwidth type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `share_bandwidth_types` - Indicates the shared bandwidth types.

  The [share_bandwidth_types](#share_bandwidth_types_struct) structure is documented below.

<a name="share_bandwidth_types_struct"></a>
The `share_bandwidth_types` block supports:

* `id` - Indicates the ID of the supported bandwidth type.

* `bandwidth_type` - Indicates the bandwidth type.

* `name_en` - Indicates the English description of the bandwidth type.

* `name_zh` - Indicates the Chinese description of the bandwidth type.

* `description` - Indicates the description of the bandwidth type.

* `public_border_group` - Indicates whether the bandwidth type is at central site or edge site.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.
