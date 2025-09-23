---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_logtank"
description: ""
---

# huaweicloud_smn_logtank

Manages an SMN logtank resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
resource "huaweicloud_smn_topic" "topic_test" {
  name = "topic_test"
}

resource "huaweicloud_lts_group" "lts_group_test" {
  group_name  = "lts_group_test"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "lts_stream_test" {
  group_id    = huaweicloud_lts_group.lts_group_test.id
  stream_name = "lts_stream_test"
}

resource "huaweicloud_smn_logtank" "logtank_test" {
  topic_urn     = huaweicloud_smn_topic.topic_test.topic_urn
  log_group_id  = huaweicloud_lts_group.lts_group_test.id
  log_stream_id = huaweicloud_lts_stream.lts_stream_test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN logtank resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `topic_urn` - (Required, String, ForceNew) Resource identifier of a topic, which is unique.
  Changing this parameter will create a new resource.

* `log_group_id` - (Required, String) The lts log group ID.

* `log_stream_id` - (Required, String) The lts log stream ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the topic URN.

* `logtank_id` - The ID of the logtank.

* `created_at` - Time when the logtank was created.

* `updated_at` - Time when the logtank was updated.

## Import

SMN logtank can be imported using the `topic_urn` or using the `topic_urn` and `logtank_id` separated by a slash e.g.

```bash
$ terraform import huaweicloud_smn_logtank.logtank_test urn:smn:cn-south-1:09f960944c80f4802f85c003e0ed1d98:logtank_test
```

or

```bash
$ terraform import huaweicloud_smn_logtank.logtank_test urn:smn:cn-south-1:09f960944c80f4802f85c003e0ed1d98:logtank_test/d9dbc3baee5c43d18a79b3fe29292003
```
