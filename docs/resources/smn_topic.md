---
subcategory: "Simple Message Notification (SMN)"
---

# huaweicloud\_smn\_topic

Manages a SMN Topic resource within HuaweiCloud.
This is an alternative to `huaweicloud_smn_topic_v2`

## Example Usage

```hcl
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the SMN topic resource. If omitted, the provider-level region will be used. Changing this creates a new SMN Topic resource.

* `name` - (Required) The name of the topic to be created.

* `display_name` - (Optional) Topic display name, which is presented as the
    name of the email sender in an email message.

* `topic_urn` - (Optional) Resource identifier of a topic, which is unique.

* `push_policy` - (Optional) Message pushing policy. 0 indicates that the message
    sending fails and the message is cached in the queue. 1 indicates that the
    failed message is discarded.

## Attributes Reference

The following attributes are exported:

* `create_time` - Time when the topic was created.

* `update_time` - Time when the topic was updated.

