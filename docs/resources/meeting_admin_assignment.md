---
subcategory: "Meeting"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_meeting_admin_assignment"
description: ""
---

# huaweicloud_meeting_admin_assignment

Using this resource to assign an administrator role to a user within HuaweiCloud.

## Example Usage

### Assign an administrator role to a user

```hcl
variable "app_id" {}
variable "app_key" {}
variable "user_account" {}

resource "huaweicloud_meeting_admin_assignment" "test" {
  app_id  = var.app_id
  app_key = var.app_key

  account = var.user_account
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Optional, String, ForceNew) Specifies the (HUAWEI Cloud meeting) user account name to which the
  default administrator belongs. Changing this parameter will create a new resource.

* `account_password` - (Optional, String, ForceNew) Specifies the user password.
  Required if `account_name` is set. Changing this parameter will create a new resource.

* `app_id` - (Optional, String, ForceNew) Specifies the ID of the Third-party application.
  Changing this parameter will create a new resource.

  -> You can apply for an application and obtain the App ID and App Key in the console.

* `app_key` - (Optional, String, ForceNew) Specifies the Key information of the Third-party APP.
  Required if `app_id` is set. Changing this parameter will create a new resource.

-> Exactly one of account authorization and application authorization you must select.

* `account` - (Required, String, ForceNew) Specifies the user account to be assigned the administrator role.
  The value can contain `1` to `64` characters.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (user account).

## Import

The assignment relationships can be imported using their `id` and authorization parameters, separated by slashes, e.g.

Import an administrator assignment and authenticated by account.

```bash
$ terraform import huaweicloud_meeting_admin_assignment.test <id>/<account_name>/<account_password>
```

Import an administrator assignment and authenticated by `APP ID`/`APP Key`.

```bash
$ terraform import huaweicloud_meeting_admin_assignment.test <id>/<app_id>/<app_key>/<corp_id>/<user_id>
```

For this resource, the `corp_id` and `user_id` are never used, you can omit them but the slashes cannot be missing.
