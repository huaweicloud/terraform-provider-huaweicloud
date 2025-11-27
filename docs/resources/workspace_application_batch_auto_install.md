---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_batch_auto_install"
description: |-
  Use this resource to batch auto install applications within HuaweiCloud.
---

# huaweicloud_workspace_application_batch_auto_install

Use this resource to batch auto install applications within HuaweiCloud.

-> This resource is a one-time action resource for batch auto installing applications. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

-> Applications must support silent installation or unzip installation for batch auto installation.

## Example Usage

### Batch auto install applications for all users

```hcl
resource "huaweicloud_workspace_application_batch_auto_install" "test" {
  app_ids      = ["app_id_1", "app_id_2"]
  assign_scope = "ALL_USER"
}
```

### Batch auto install applications with assigned users

```hcl
variable "application_ids" {
  type = list(string)
}
variable "user_names" {
  type = list(string)
}

resource "huaweicloud_workspace_application_batch_auto_install" "test" {
  app_ids      = var.application_ids
  assign_scope = "ASSIGN_USER"

  dynamic "users" {
    for_each = var.user_names

    content {
      account      = users.value
      account_type = "SIMPLE"
    }
  }
}
```

### Batch auto install applications with user groups

```hcl
resource "huaweicloud_workspace_application_batch_auto_install" "test" {
  app_ids      = ["app_id_1", "app_id_2"]
  assign_scope = "ASSIGN_USER"

  users {
    account       = "group_name"
    account_type  = "USER_GROUP"
    domain        = "example.com"
    platform_type = "AD"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the applications are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_ids` - (Required, List, NonUpdatable) Specifies the list of application IDs to be automatically installed.  
  The number of applications in a single request ranges from `1` to `50`.

* `assign_scope` - (Required, String, NonUpdatable) Specifies the assignment scope.  
  The valid values are as follows:
  + **ALL_USER** - All users can access the applications.  
    When set to this value, all application authorizations will be modified to all users.
  + **ASSIGN_USER** - Only assigned users can access the applications.  
    When set to this value, the `users` parameter must be specified. If users do not have access permissions for the
    corresponding applications, the permissions will be automatically added.

* `users` - (Optional, List, NonUpdatable) Specifies the list of users.  
  Required when `assign_scope` is set to **ASSIGN_USER**.  
  The number of users in a single request ranges from `1` to `50`.  
  The [users](#application_batch_auto_install_user) structure is documented below.

<a name="application_batch_auto_install_user"></a>
The `users` block supports:

* `account` - (Required, String) Specifies the account name.  
  The account format must be: account(group).

* `account_type` - (Required, String) Specifies the account type.  
  The valid values are as follows:
  + **SIMPLE** - Simple user.
  + **USER_GROUP** - User group.

* `domain` - (Optional, String) Specifies the domain name.  
  Required for user groups, and defaults to local.com if not specified.

* `platform_type` - (Optional, String) Specifies the platform type.  
  The valid values are as follows:
  + **AD** - AD domain.
  + **LOCAL** - LiteAs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
