---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_schemas"
description: |-
  Use this data source to get the query schemas within HuaweiCloud Resource Center.
---

# huaweicloud_rc_schemas

Use this data source to get the query schemas within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_schemas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `value` - The schema list.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `type` - The resource type.

* `schema` - The schema content in JSON format.
