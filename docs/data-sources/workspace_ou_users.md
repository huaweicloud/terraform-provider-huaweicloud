---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_ou_users"
description: |-
  Use this data source to query the user list under the specified OU within HuaweiCloud.
---

# huaweicloud_workspace_ou_users

Use this data source to query the user list under the specified OU within HuaweiCloud.

## Example Usage

### Query all users under the specified OU

```hcl
variable "ou_dn" {}

data "huaweicloud_workspace_ou_users" "test" {
  ou_dn = var.ou_dn
}
```

### Query users by user name

```hcl
variable "ou_dn" {}
variable "user_name" {}

data "huaweicloud_workspace_ou_users" "test" {
  ou_dn     = var.ou_dn
  user_name = var.user_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the OU users are located.  
  If omitted, the provider-level region will be used.

* `ou_dn` - (Required, String) Specifies the distinguished name (DN) of the OU.

* `user_name` - (Optional, String) Specifies the name of the user to which the OU belongs.
  Fuzzy matching is supported.

* `has_existed` - (Optional, Bool) Specifies whether the user already exists in the user list.  
  If omitted, all users under the OU will be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of users that match the filter parameters.  
  The [users](#workspace_ou_users) structure is documented below.

* `enable_create_count` - The number of users that can be created.

<a name="workspace_ou_users"></a>
The `users` block supports:

* `name` - The name of the user.

* `expired_time` - The expiration time of the user.  
  `-1` indicates that the user never expires.

* `has_existed` - Whether the user already exists in the user list.
