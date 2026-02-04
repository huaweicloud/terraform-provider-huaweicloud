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
variable "region_id" {
  type    = string
  default = "cn-north-4"
}

data "huaweicloud_identity_regions" "test" {
  region_id = var.region_id
}
```

## Argument Reference

* `region_id` - (Optional, String) Specifies the id of the region to be queried.

## Attribute Reference

* `regions` - The information of region list.  
  The [regions](#identity_regions_regions) structure is documented below.

<a name="identity_regions_regions"></a>
The `regions` block contains:

* `id` - The ID of the region.

* `type` - The type of the region.

* `description` - The description of the region.

* `link` - The resource link of the region.

* `locales` - The map of localized region names, where the key is the language code and the value is the
  region name in that language.
