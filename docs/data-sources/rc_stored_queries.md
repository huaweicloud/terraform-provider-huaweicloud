---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_stored_queries"
description: |-
  Use this data source to get the stored ResourceQL queries within HuaweiCloud Resource Center.
---

# huaweicloud_rc_stored_queries

Use this data source to get the stored ResourceQL queries within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_stored_queries" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the ResourceQL name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `value` - The stored query list.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The value block supports:

* `id` - The unique ID of the operation.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `deleted_at` - The deletion time.

* `expired_at` - The expiration time.

* `name` - The ResourceQL name.

* `description` - The ResourceQL description.

* `expression` - The ResourceQL expression.

* `created` - The ResourceQL created time.

* `updated` - The ResourceQL updated time.
