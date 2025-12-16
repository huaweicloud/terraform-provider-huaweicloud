---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_user_batch_detach"
description: |-
  Use this resource to batch detach users from desktops within HuaweiCloud.
---

# huaweicloud_workspace_desktop_user_batch_detach

Use this resource to batch detach users from desktops within HuaweiCloud.

-> This resource is a one-time action resource for batch detaching users from desktops. Deleting this resource will
   not undo the detachment, but will only remove the resource information from the tfstate file.

## Example Usage

### Detach all users from a desktop

```hcl
variable "desktop_id" {}

resource "huaweicloud_workspace_desktop_user_batch_detach" "test" {
  desktops {
    desktop_id          = var.desktop_id
    is_detach_all_users = true
  }
}
```

### Detach specific users from a desktop

```hcl
variable "desktop_id" {}
variable "user_id" {}

resource "huaweicloud_workspace_desktop_user_batch_detach" "test" {
  desktops {
    desktop_id = var.desktop_id

    detach_user_infos {
      user_id = var.user_id
    }
  }
}
```

## Argument

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktops and users are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `desktops` - (Required, List, NonUpdatable) Specifies the list of desktop user detach information.  
  The [desktops](#desktop_user_batch_detach_desktops) structure is documented below.

<a name="desktop_user_batch_detach_desktops"></a>
The `desktops` block supports:

* `desktop_id` - (Optional, String, NonUpdatable) Specifies the ID of the desktop to be detached.

* `is_detach_all_users` - (Optional, Bool, NonUpdatable) Specifies whether to detach all users.  

  ->**Note** If `is_detach_all_users` set to **true**, `detach_user_infos` will be invalid.

* `detach_user_infos` - (Optional, List, NonUpdatable) Specifies the list of users to be detached.  
  The [detach_user_infos](#desktop_user_batch_detach_user_infos) structure is documented below.

<a name="desktop_user_batch_detach_user_infos"></a>
The `detach_user_infos` block supports:

* `user_id` - (Optional, String, NonUpdatable) Specifies the ID of the user.  
  This parameter is **Required** when `type` is **USER**.

* `user_name` - (Optional, String, NonUpdatable) Specifies the name of the user or user group.

* `user_group` - (Optional, String, NonUpdatable) Specifies the user group which the user belongs to.  
  The valid values are as follows:
  + **sudo**
  + **default**
  + **administrators**
  + **users**

* `type` - (Optional, String, NonUpdatable) Specifies the type of the object.  
  The valid values are as follows:
  + **USER**
  + **GROUP**

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
