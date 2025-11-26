---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_group_authorization_notification_resend"
description: |-
  Use this resource to resend application group authorization notifications of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_group_authorization_notification_resend

Use this resource to resend application group authorization notifications of the Workspace APP within HuaweiCloud.

-> This resource is a one-time action resource for resending application group authorization notifications. Deleting
  this resource will not clear the corresponding request record, but will only remove the resource information from
  the tfstate file.

## Example Usage

### Resend by authorization notification records

```hcl
variable "notification_records" {
  type = list(object({
    id                = string
    account           = optional(string)
    account_auth_type = optional(string)
    account_auth_name = optional(string)
    app_group_id      = optional(string)
    app_group_name    = optional(string)
    mail_send_type    = optional(string)
  }))
}

resource "huaweicloud_workspace_app_group_authorization_notification_resend" "test" {
  is_notification_record = true

  dynamic "records" {
    for_each = var.notification_records

    content {
      id                = records.value["id"]
      account           = records.value["account"]
      account_auth_type = records.value["account_auth_type"]
      account_auth_name = records.value["account_auth_name"]
      app_group_id      = records.value["app_group_id"]
      app_group_name    = records.value["app_group_name"]
      mail_send_type    = records.value["mail_send_type"]
    }
  }
}
```

### Resend by authorization records

```hcl
variable "authorization_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_group_authorization_notification_resend" "test" {
  dynamic "records" {
    for_each = var.authorization_ids

    content {
      id = records.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application group authorization notifications
  to be resent are located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `records` - (Required, List, NonUpdatable) Specifies the list of record IDs to resend authorization notification.
  The [records](#app_group_authorization_notification_resend_records) structure is documented below.

* `is_notification_record` - (Optional, Bool, NonUpdatable) Specifies whether to resend according to the authorization
  notification records.  
  Defaults to **false**.
  + **true** - Resend according to the authorization notification records.
  + **false** - Resend according to the authorization records.

<a name="app_group_authorization_notification_resend_records"></a>
The `records` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the ID of the authorization notification record or authorization
  record.

* `account` - (Optional, String, NonUpdatable) Specifies the name of the authorized account.

* `account_auth_type` - (Optional, String, NonUpdatable) Specifies the type of the authorized object.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**

* `account_auth_name` - (Optional, String, NonUpdatable) Specifies the name of the authorized object.

* `app_group_id` - (Optional, String, NonUpdatable) Specifies the ID of the application group.

* `app_group_name` - (Optional, String, NonUpdatable) Specifies the name of the application group.

* `mail_send_type` - (Optional, String, NonUpdatable) Specifies the type of authorization notification.  
  The valid values are as follows:
  + **ADD_GROUP_AUTHORIZATION** - Add group authorization email.
  + **DEL_GROUP_AUTHORIZATION** - Delete group authorization email.
  + **ADD_GROUP_AUTHORIZATION_SMS** - Add group authorization SMS.
  + **DEL_GROUP_AUTHORIZATION_SMS** - Delete group authorization SMS.

-> The parameters `mail_send_type`, `account`, `account_auth_type`, `account_auth_name`, `app_group_id`, and
   `app_group_name` are valid only if `is_notification_record` is set to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
