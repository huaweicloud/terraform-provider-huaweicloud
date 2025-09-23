---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_message_detection"
description: |-
  Manages a SMN message detection resource within HuaweiCloud.
---

# huaweicloud_smn_message_detection

Manages a SMN message detection resource within HuaweiCloud.

## Example Usage

```hcl
variable "topic_urn" {}

resource "huaweicloud_smn_message_detection" "test" {
  topic_urn = var.topic_urn
  protocol  = "https"
  endpoint  = "https://example.com/notification/action"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the resource identifier of a topic.

* `protocol` - (Required, String, NonUpdatable) Specifies the protocol type. The value can be **http** or **https**.

* `endpoint` - (Required, String, NonUpdatable) Specifies the endpoint address to be detected.
  The address must start with **http://** or **https://** and cannot be left blank.

* `extension` - (Optional, Map, NonUpdatable) Specifies the extended key/value for subscriptions over HTTP or HTTPS.
  These key/value pairs will be carried as the request header when HTTP or HTTPS messages are sent.
  The key/value must meet the following requirements:
  + **key** can contain letters, digits, and hyphens (-). **key** cannot end with a hyphen (-) nor contain
    consecutive hyphens (-). **key** must start with **x-**but not **x-smn**. Examples: **x-abc-cba** or **x-abc**.
    **key** is case insensitive. **key** must be unique.
  + **value** must be an ASCII code. Unicode characters are not supported. Spaces are allowed.
  + You can specify up to 10 key/value pairs.
  + The total length of all key/value pairs cannot exceed **1,024** characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `result` - The message detection result.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
