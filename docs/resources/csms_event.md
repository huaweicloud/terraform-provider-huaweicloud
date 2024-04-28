---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_event"
description: ""
---

# huaweicloud_csms_event

Manages a CSMS (Cloud Secret Management Service) event resource within HuaweiCloud.

## Example Usage

```hcl
variable "notification_target_type" {}
variable "notification_target_id" {}
variable "notification_target_name" {}

resource "huaweicloud_csms_event" "test" {
  name                     = "test_name"
  event_types              = ["SECRET_VERSION_CREATED", "SECRET_ROTATED"]
  status                   = "ENABLED"
  notification_target_type = var.notification_target_type
  notification_target_id   = var.notification_target_id
  notification_target_name = var.notification_target_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the event name. The valid length is limited from `1` to `64`.
  Only letters, digits, underscores (_) and hyphens (-) are allowed.

  Changing this parameter will create a new resource.

* `event_types` - (Required, List) Specifies the event type list. Valid values are:
  + **SECRET_VERSION_CREATED**: Triggered when a version of a secret is created.
  + **SECRET_VERSION_EXPIRED**: Triggered when a secret version expires, and only once per expiration.
  + **SECRET_ROTATED**: Triggered when a secret is rotated. Currently, only RDS secrets can be automatically rotated.
  + **SECRET_DELETED**: Triggered when a secret is deleted.

* `status` - (Required, String) Specifies the event status. Valid values are **ENABLED** and **DISABLED**.
  Only the event in **ENABLED** status can be triggered.

* `notification_target_type` - (Required, String) Specifies the notification target type.
  Currently, only **SMN** is supported.

* `notification_target_id` - (Required, String) Specifies the notification target ID.
  Currently, only SMN topic URN is supported.

* `notification_target_name` - (Required, String) Specifies the notification target name.
  Currently, only SMN topic name is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `name`.

* `event_id` - The event ID.

## Import

The CSMS event can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_csms_event.test <id>
```
