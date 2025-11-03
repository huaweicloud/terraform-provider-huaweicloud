---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_notification_records"
description: |-
  Use this data source to get the list of notification records.
---

# huaweicloud_csms_notification_records

Use this data source to get the list of notification records.

-> Only event notification records generated within three months are stored.

## Example Usage

```hcl
data "huaweicloud_csms_notification_records" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of the versions.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `event_name` - The event name.

* `trigger_event_type` - The event type.
  The valid values are as follows:
  + **SECRET_VERSION_CREATED**
  + **SECRET_VERSION_EXPIRED**
  + **SECRET_ROTATED**
  + **SECRET_DELETED**
  + **SECRET_ROTATED_FAILED**

* `secret_name` - The secret name.

* `secret_type` - The secret type.
  The valid values are as follows:
  + **COMMON**
  + **RDS-FG**
  + **GaussDB-FG**

* `notification_target_name` - The name of the object to which the event notification is sent.

* `notification_target_id` - The ID of the object to which the event notification is sent.

* `notification_content` - The event notification content.

* `notification_status` - The event notification status.
  The valid values are as follows:
  + **SUCCESS**
  + **FAIL**
  + **INVALID**

* `create_time` - The creation time of the event notification record.
