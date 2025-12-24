---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_batch_authorize"
description: |-
  Use this resource to batch authorize applications within HuaweiCloud.
---

# huaweicloud_workspace_application_batch_authorize

Use this resource to batch authorize applications within HuaweiCloud.

-> This resource is a one-time action resource for batch authorizing applications. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Batch authorize applications with assigned users

```hcl
variable "application_ids" {
  type = list(string)
}
variable "user_names_to_be_authorized" {
  type = list(string)
}

resource "huaweicloud_workspace_application_batch_authorize" "test" {
  app_ids            = var.application_ids
  authorization_type = "ASSIGN_USER"

  dynamic "users" {
    for_each = var.user_names_to_be_authorized

    content {
      account      = users.value.name
      account_type = "SIMPLE"
    }
  }
}
```

### Batch authorize applications for all users

```hcl
resource "huaweicloud_workspace_application_batch_authorize" "test" {
  app_ids            = ["app_id_1", "app_id_2"]
  authorization_type = "ALL_USER"
}
```

### Batch authorize applications with user cancellation

```hcl
variable "application_ids" {
  type = list(string)
}
variable "authorized_user_names" {}

resource "huaweicloud_workspace_application_batch_authorize" "test" {
  app_ids            = var.application_ids
  authorization_type = "ASSIGN_USER"

  dynamic "del_users" {
    for_each = var.authorized_user_names

    content {
      account      = del_users.value.name
      account_type = "SIMPLE"
    }
  }

}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application batch authorization is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_ids` - (Required, List, NonUpdatable) Specifies the list of application IDs to be authorized or unauthorized.

* `authorization_type` - (Required, String, NonUpdatable) Specifies the authorization type.  
  Valid values are:
  + **ALL_USER** - All users can access the applications.
  + **ASSIGN_USER** - Only assigned users can access the applications.

* `users` - (Optional, List, NonUpdatable) Specifies the list of users to be authorized.  
  The maximum number is `100`.  
  The [users](#application_batch_authorize_user) structure is documented below.

* `del_users` - (Optional, List, NonUpdatable) Specifies the list of users to be unauthorized.  
  The maximum number is `100`.  
  The [del_users](#application_batch_authorize_user) structure is documented below.

<a name="application_batch_authorize_user"></a>
The `users` and `del_users` block supports:

* `account` - (Required, String, NonUpdatable) Specifies the account name.

* `account_type` - (Required, String, NonUpdatable) Specifies the account type.  
  The valid values are as follows:
  + **SIMPLE** - Simple user.
  + **USER_GROUP** - User group.

* `domain` - (Optional, String, NonUpdatable) Specifies the domain name. Required for user groups.

* `platform_type` - (Optional, String, NonUpdatable) Specifies the platform type.  
  The valid values are as follows:
  + **AD** - AD domain.
  + **LOCAL** - LiteAs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
