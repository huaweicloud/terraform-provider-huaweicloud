---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_group_authorization"
description: |-
  Manages an APP group authorization resource within HuaweiCloud.
---

# huaweicloud_workspace_app_group_authorization

Manages an APP group authorization resource within HuaweiCloud.

-> Deleting this resource will revoke authorization for the users or user groups.

## Example Usage

```hcl
variable "app_group_id" {}
variable "user_groups" {
  type = list(object({
    id   = string
    name = string
  }))
}

resource "huaweicloud_workspace_app_group_authorization" "test" {
  app_group_id = var.app_group_id

  dynamic "accounts" {
    for_each = var.user_groups

    content {
      id      = accounts.value.id
      account = accounts.value.name
      type    = "USER_GROUP"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `app_group_id` - (Required, String, ForceNew) Specifies the ID of the application group.
  Changing this creates a new resource.

* `accounts` - (Required, List, ForceNew) Specifies the list of the accounts to be authorized. The maximum length is `50`.
  Changing this creates a new resource.  
  The [accounts](#app_group_auth_accounts) structure is documented below.

  -> If the parameter contains non-existent objects, the resource creation will fail, but the remaining existing objects
     will be authorized successfully.

<a name="app_group_auth_accounts"></a>
The `accounts` block supports:

* `id` - (Optional, String, ForceNew) Specifies the ID of the user (group).
  Changing this creates a new resource.  
  This parameter is required when `type` is set to **USER_GROUP**.

* `account` - (Required, String, ForceNew) Specifies the name of the user (group).
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the object to be authorized.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `app_group_id`).
