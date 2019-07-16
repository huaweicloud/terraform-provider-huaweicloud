---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topic_v2"
sidebar_current: "docs-huaweicloud-resource-smn-topic-v2"
description: |-
  Manages a V2 topic resource within HuaweiCloud.
---

# huaweicloud\_smn\_topic\_v2

Manages a V2 topic resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the topic to be created.

* `display_name` - (Optional) Topic display name, which is presented as the
    name of the email sender in an email message.

* `topic_urn` - (Optional) Resource identifier of a topic, which is unique.

* `push_policy` - (Optional) Message pushing policy. 0 indicates that the message
    sending fails and the message is cached in the queue. 1 indicates that the
    failed message is discarded.

* `create_time` - (Optional) Time when the topic was created.

* `update_time` - (Optional) Time when the topic was updated.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `topic_urn` - See Argument Reference above.
* `push_policy` - See Argument Reference above.
* `create_time` - See Argument Reference above.
* `update_time` - See Argument Reference above.
