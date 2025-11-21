---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_user_sessions"
description: |-
  Use this data source to get the Identity Center user sessions.
---

# huaweicloud_identitycenter_user_sessions

Use this data source to get the Identity Center user sessions.

## Example Usage

```hcl
variable "user_id" {}
variable "identity_store_id" {}

data "huaweicloud_identitycenter_user_sessions" "test"{
  user_id           = var.user_id
  identity_store_id = var.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `user_id` - (Required, String) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `session_list` - The list of the user session.
  The [session_list](#session_list_struct) structure is documented below.

<a name="session_list_struct"></a>
The `session_list` block supports:

* `creation_time` - The creation time of the user session.

* `ip_address` - The ip address of the user session.

* `session_id` - The ID of the user session.

* `session_not_valid_after` - The validity period of the user session.

* `user_agent` - The user agent of the user session.
