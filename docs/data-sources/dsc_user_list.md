---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_user_list"
description: |-
  Use this data source to get the list of user accounts under the current project within HuaweiCloud.
---

# huaweicloud_dsc_user_list

Use this data source to get the list of user accounts under the current project within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_user_list" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `user_list` - The list of user information.

  The [user_list](#user_list_struct) structure is documented below.

<a name="user_list_struct"></a>
The `user_list` block supports:

* `user_id` - The unique ID of the user.

* `user_name` - The name of the user.
