---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_group_authorization_notification_records"
description: |-
  Use this data source to query the application group authorization notification record list of
  the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_group_authorization_notification_records

Use this data source to query the application group authorization notification record list of
the Workspace APP within HuaweiCloud.

## Example Usage

### Query all authorization notification records under the specified application group

```hcl
variable "application_group_id" {}

data "huaweicloud_workspace_app_group_authorization_notification_records" "test" {
  app_group_id = var.application_group_id
}
```

### Query the authorization notification records by account name

```hcl
variable "application_group_id" {}
variable "account_name" {}

data "huaweicloud_workspace_app_group_authorization_notification_records" "test" {
  app_group_id = var.application_group_id
  account      = var.account_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the authorization notification records are located.  
  If omitted, the provider-level region will be used.

* `app_group_id` - (Required, String) Specifies the ID of the application group.

* `account` - (Optional, String) Specifies the name of the authorized user (group).  
  Fuzzy search is supported.

* `mail_send_type` - (Optional, String) Specifies the type of authorization operation.  
  The valid values are as follows:
  + **ADD_GROUP_AUTHORIZATION**
  + **DEL_GROUP_AUTHORIZATION**

* `mail_send_result` - (Optional, String) Specifies the result of the notification sending.
  The valid values are as follows:
  + **SUCCESS**
  + **FAIL**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The authorization notification record list that match filter parameters.  
  The [records](#app_group_authorization_notification_records) structure is documented below.

<a name="app_group_authorization_notification_records"></a>
The `records` block supports:

* `id` - The ID of the record.

* `account` - The name of the user (group).

* `account_auth_type` - The type of the account.
  + **USER**
  + **USER_GROUP**

* `account_auth_name` - The name of the authorized object.

* `app_group_id` - The ID of the application group.

* `app_group_name` - The name of the application group.

* `mail_send_type` - The type of authorization operation.
  + **ADD_GROUP_AUTHORIZATION**
  + **DEL_GROUP_AUTHORIZATION**
  + **ADD_GROUP_AUTHORIZATION_SMS**
  + **DEL_GROUP_AUTHORIZATION_SMS**

* `mail_send_result` - The result of the notification sending.

* `error_msg` - The error message when the notification failed to be sent.

* `send_at` - The time when the authorization notification was sent.
