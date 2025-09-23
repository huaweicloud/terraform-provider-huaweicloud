---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_users"
description: |-
  Use this data source to get the list of users who have performed operations in the last 30 days.
---

# huaweicloud_cts_users

Use this data source to get the list of users who have performed operations in the last 30 days.

## Example Usage

```hcl
data "huaweicloud_cts_users" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of users who have performed operations in the last 30 days.

  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `name` - The username.
