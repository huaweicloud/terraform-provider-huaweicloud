---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topics"
description: ""
---

# huaweicloud_smn_topics

Use this data source to get a list of SMN topics.

## Example Usage

```hcl
variable "topic_name" {}

data "huaweicloud_smn_topics" "tpoic_1" {
  name = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the SMN topics. If omitted, the
  provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the topic.

* `topic_urn` - (Optional, String) Specifies the topic URN.

* `display_name` - (Optional, String) Specifies the topic display name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the SMN topic.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID

* `topics` - An array of SMN topics found. Structure is documented below.

The `topics` block supports:

* `name` - The name of the topic.

* `id` - The topic ID. The value is the topic URN.

* `topic_urn` - The topic URN.

* `display_name` - The topic display name.

* `enterprise_project_id` - The enterprise project ID of the SMN topic.

* `push_policy` - Message pushing policy.
  + **0**: indicates that the message sending fails and the message is cached in the queue.
  + **1**: indicates that the failed message is discarded.

* `tags` - The tags of the SMN topic, key/value pair format.
