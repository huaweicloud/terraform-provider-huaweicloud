---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_supported_areas"
description: |-
  Use this data source to get the supported geographic regions.
---

# huaweicloud_cc_supported_areas

Use this data source to get the supported geographic regions.

## Example Usage

```hcl
data "huaweicloud_cc_supported_areas" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `areas` - Indicates the geographic region list.

  The [areas](#areas_struct) structure is documented below.

<a name="areas_struct"></a>
The `areas` block supports:

* `id` - Indicates the geographic region ID.

* `name` - Indicates the geographic region name.

* `en_name` - Indicates the geographic region name in English.

* `es_name` - Indicates the geographic region name in Spanish.

* `pt_name` - Indicates the geographic region name in Portuguese.

* `station` - Indicates the site.
