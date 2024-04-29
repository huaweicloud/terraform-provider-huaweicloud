---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_user_group"
description: ""
---

# huaweicloud_workspace_user_group

Manages a Workspace user group resource within HuaweiCloud.

## Example Usage

### Create a local domain user group

```hcl
variable "group_name" {}

resource "huaweicloud_workspace_user_group" "test" {
  name        = var.group_name
  type        = "LOCAL"
  description = "Created by script"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Workspace user group resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the user group name.  
  -> AD domain user group do not support renaming.

* `type` - (Required, String) Specifies the type of user group.
  The valid values are as follows:
  + **LOCAL**: Local domain user group.
  + **AD**: AD domain user group.

* `description` - (Optional, String) Specifies the description of user group.

* `users` - (Optional, List) Specifies the user information under the user group.  
  The [users](#workspace_users) structure is documented below.

<a name="workspace_users"></a>
The `users` block supports:

* `id` - (Required, String) Specifies the user ID to be added to the user group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The user group ID in UUID format.

* `users` - (Optional, List) The user information under the user group.  
  The [users](#workspace_users_attr) structure is documented below.

* `created_at` - The creation time of the user group.

<a name="workspace_users_attr"></a>
The `users` block supports:

* `name` - The name of user.

* `email` - The email of user.

* `phone` - The phone of user.

* `description` - The description of user.

* `total_desktops` - The number of desktops the user has.

## Import

The user group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_user_group.test <id>
```
