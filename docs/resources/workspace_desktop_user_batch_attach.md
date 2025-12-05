---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_user_batch_attach"
description: |-
  Use this resource to batch attach users to desktops within HuaweiCloud.
---

# huaweicloud_workspace_desktop_user_batch_attach

Use this resource to batch attach users to desktops within HuaweiCloud.

-> This resource is a one-time action resource for batch attaching users to a desktop. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Attach one user to desktop

```hcl
variable "desktop_id" {}
variable "user_name" {}

resource "huaweicloud_workspace_desktop_user_batch_attach" "test" {
  desktops {
    desktop_id = var.desktop_id
    user_name  = var.user_name
  }
}
```

### Batch attach users to desktop

```hcl
variable "desktop_id" {}
variable "attach_user_information" {
  type = list(object({
    user_name  = string
    user_group = string
  }))
}

resource "huaweicloud_workspace_desktop_user_batch_attach" "test" {
  desktops {
    desktop_id = var.desktop_id
    
    dynamic "attach_user_infos" {
      for_each = var.attach_user_information
      
      content {
        user_name  = attach_user_infos.value.user_name
        user_group = attach_user_infos.value.user_group
      }
    }
  }
}
```

## Argument

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktops are located.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `desktops` - (Required, List, NonUpdatable) Specifies the list of desktop information to be assigned.  
  The [desktops](#desktop_user_batch_attach_desktops) structure is documented below.

* `image_id` - (Optional, String, NonUpdatable) Specifies the image ID used to change the desktop image.

* `image_type` - (Optional, String, NonUpdatable) Specifies the image type used to change the desktop image.

* `desktop_name_policy_id` - (Optional, String, NonUpdatable) Specifies the policy ID used to specify the desktop name
  generation policy.

  ->**Note** If `desktop_name` is specified, `desktop_name_policy_id` will be invalid.

<a name="desktop_user_batch_attach_desktops"></a>
The `desktops` block supports:

* `user_group` - (Required, String, NonUpdatable) Specifies the user group to which the desktop user belongs.  
  This parameter is required when attach_user_infos is empty, and attach_user_infos takes higher priority.  
  The valid values are as follows:
  + **sudo**: Linux administrator group.
  + **default**: Linux default user group.
  + **administrators**: Windows administrator group. Administrators have full access to the desktop and can make any
    changes needed (except disable operations).
  + **users**: Windows standard user group. Standard users can use most software and can change system settings that
    do not affect other users.

* `desktop_id` - (Optional, String, NonUpdatable) Specifies the ID of the desktop to be assigned.

* `computer_name` - (Optional, String, NonUpdatable) Specifies the desktop name.  
  The desktop name must be unique. Only uppercase letters, lowercase letters, numbers, hyphens (-) and underscores (_)
  are allowed. It must start with a letter and cannot end with a hyphen (-).
  The length ranges from `1` to `15` characters.

* `user_name` - (Optional, String, NonUpdatable) Specifies the user to whom the desktop belongs.  
  After the desktop is successfully assigned, this user can log in to the desktop. Only uppercase letters, lowercase
  letters, numbers, hyphens (-) and underscores (_) are allowed. For **LITE_AD** domain type, use lowercase or uppercase
  letters at the beginning, with a length range of `1` to `20`. For **LOCAL_AD** domain type, usernames can begin with
  lowercase letters, uppercase letters or numbers, with a length range of `1` to `32`. Windows desktop users support
  up to `20` characters, and Linux desktop users support up to `32` characters. The username cannot be the same as the
  assigned machine name.

* `user_email` - (Optional, String, NonUpdatable) Specifies valid user email.  
  After the desktop is successfully assigned, the system will notify the user via email.

* `is_clear_data` - (Optional, Bool, NonUpdatable) Specifies whether to clean up desktop data when binding.  
  This field only takes effect when the unbinding and binding are for the same user. The default value is **true**.

* `attach_user_infos` - (Optional, List, NonUpdatable) Specifies the list of user information to be assigned.  
  This is only valid when assigning a multi-user desktop to multiple users.  
  The [attach_user_infos](#desktop_user_batch_attach_user_infos) structure is documented below.

<a name="desktop_user_batch_attach_user_infos"></a>
The `attach_user_infos` block supports:

* `type` - (Optional, String, NonUpdatable) Specifies the object type.  
  The valid values are as follows:
  + **USER**
  + **GROUP**

* `user_id` - (Optional, String, NonUpdatable) Specifies the user ID. The backend service will verify whether the
  group ID exists.
  + When the type field is **USER**, fill in the user ID.
  + When the type field is **GROUP**, fill in the user group ID.

* `user_name` - (Optional, String, NonUpdatable) Specifies the name of the desktop assignment object.  
  When the type is **USER**, fill in the user name. When the type is **GROUP**, fill in the user group name.
  + When the type is **USER**: The user to whom the desktop belongs. After the desktop is successfully assigned, this
    user can log in to the desktop. Only uppercase letters, lowercase letters, numbers, hyphens (-) and underscores (_)
    are allowed.
    - For **LITE_AD** domain type, use lowercase or uppercase letters at the beginning, with a length range
      of `1` to `20`.
    - For **LOCAL_AD** domain type, usernames can begin with lowercase letters, uppercase letters or numbers, with a
      length range of `1` to `64`. Windows desktop users support up to `20` characters, and Linux desktop users support
      up to `64` characters. The backend service will verify whether the username exists, and the username cannot
      be the same as the machine name.
  + When the type is **GROUP**: Can only be Chinese, letters, numbers and special symbols `-`, `_`.

* `user_group` - (Optional, String, NonUpdatable) Specifies the user group to which the desktop user belongs.  
  The valid values are as follows:
  + **sudo**: Linux administrator group.
  + **default**: Linux default user group.
  + **administrators**: Windows administrator group. Administrators have full access to the desktop and can make any
    changes needed (except disable operations).
  + **users**: Windows standard user group. Standard users can use most software and can change system settings that
    do not affect other users.

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
