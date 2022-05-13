---
subcategory: "Simple Message Notification (SMN)"
---

# huaweicloud_smn_topic

Manages an SMN Topic resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN topic resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic to be created. The name can contains of 1 to 255
  characters and must start with a letter or digit, and can only contain letters, digits, underscores (_), and hyphens (-).
  Changing this parameter will create a new resource.

* `display_name` - (Optional, String) Specifies the topic display name, which is presented as the name of the email
  sender in an email message. The name can contains of 0 to 192 characters.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the SMN Topic, Value 0
  indicates the default enterprise project. Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the tags of the SMN topic, key/value pair format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the topic urn.

* `topic_urn` - Resource identifier of a topic, which is unique.

* `push_policy` - Message pushing policy. 0 indicates that the message sending fails and the message is cached in the
  queue. 1 indicates that the failed message is discarded.

* `create_time` - Time when the topic was created.

* `update_time` - Time when the topic was updated.

## Import

SMN topic can be imported using the `id` (topic urn), e.g.

```
$ terraform import huaweicloud_smn_topic.topic_1 urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:topic_1
```
