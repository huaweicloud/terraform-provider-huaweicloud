---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_users"
description: |-
  Use this data source to query the Workspace users under a specified region within HuaweiCloud.
---

# huaweicloud_workspace_users

Use this data source to query the Workspace users under a specified region within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_users" "test" {}
```

### Filter by user name

```hcl
variable "user_name" {}

data "huaweicloud_workspace_users" "test" {
  user_name = var.user_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the users are located.  
  If omitted, the provider-level region will be used.

* `user_name` - (Optional, String) Specifies the user name to be queried.

* `description` - (Optional, String) Specifies the user description for fuzzy matching.

* `active_type` - (Optional, String) Specifies the activation type of the user.  
  The valid values are as follows:
  + **USER_ACTIVATE**
  + **ADMIN_ACTIVATE**

* `group_name` - (Optional, String) Specifies the user group name for exact matching.

* `is_query_total_desktops` - (Optional, Bool) Specifies whether to query the number of desktops bound to the user.  
  The default value is **true**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of users that matched filter parameters.  
  The [users](#workspace_users) structure is documented below.

<a name="workspace_users"></a>
The `users` block supports:

* `id` - The ID of the user.

* `sid` - The SID of the user.

* `user_name` - The name of user.

* `user_email` - The email address of the user.

* `total_desktops` - The total number of desktops bound to the user.

* `user_phone` - The phone number of the user.

* `active_type` - The activation type of the user.

* `is_pre_user` - Whether the user is a pre-created user.

* `account_expires` - The account expired time.

* `password_never_expired` - Whether the password never expires.

* `account_expired` - Whether the account has expired.

* `enable_change_password` - Whether the user is allowed to change password.

* `next_login_change_password` - Whether the password needs to be reset on next login.

* `description` - The description of the user.

* `locked` - Whether the account is locked.

* `disabled` - Whether the account is disabled.

* `share_space_subscription` - Whether the user has subscribed to collaboration.

* `share_space_desktops` - The number of collaboration desktops bound to the user.

* `group_names` - The list of group name that the user has joined.

* `enterprise_project_id` - The ID of the enterprise project.

* `user_info_map` - The user information mapping, including user service level, operation mode and type.
