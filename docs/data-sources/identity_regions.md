---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_regions"
description: |-
  Use this data source to query the region list within HuaweiCloud.
---

# huaweicloud_identity_regions

Use this data source to query the region list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identity_regions" "test1" {}

data "huaweicloud_identity_regions" "test2" {
  region_id = "cn-north-4"
}
```

## Argument Reference

* `region_id` - (Optional, String) Specifies the region id.

## Attribute Reference

* `regions` - Indicates the region info list.
  The [regions](#IdentityRegions_Regions) structure is documented below.

<a name="IdentityRegions_Regions"></a>
The `regions` block contains:

* `id` - Indicates the region id.

* `type` - Indicates the region type.

* `description` - Indicates the region description.

* `link` - Indicates the resource link.

* `locales` - Indicates the map of region name.
