---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_login_white_list"
description: |-
  Manages an HSS event login white list resource within HuaweiCloud.
---

# huaweicloud_hss_event_login_white_list

Manages an HSS event login white list resource within HuaweiCloud.

## Example Usage

```hcl
variable "private_ip" {}
variable "login_ip" {}
variable "login_user_name" {}

resource "huaweicloud_hss_event_login_white_list" "test" {
  private_ip      = var.private_ip
  login_ip        = var.login_ip
  login_user_name = var.login_user_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `private_ip` - (Required, String, NonUpdatable) Specifies the private IP address of the host.

* `login_ip` - (Required, String, NonUpdatable) Specifies the login IP address.

* `login_user_name` - (Required, String, NonUpdatable) Specifies the login username.

* `remarks` - (Optional, String, NonUpdatable) Specifies the remarks of the white list.

* `handle_event` - (Optional, Bool, NonUpdatable) Specifies whether to handle the related alarm events simultaneously.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `delete_all` - (Optional, Bool) Specifies whether to delete all login white lists. When set to `true`, all
  login white lists under HSS will be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `enterprise_project_name` - The enterprise project name.

* `update_time` - The update time in milliseconds.
