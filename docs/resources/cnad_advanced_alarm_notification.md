---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_alarm_notification"
description: |-
  Manages a CNAD advanced alarm notification resource within HuaweiCloud.
---

# huaweicloud_cnad_advanced_alarm_notification

Manages a CNAD advanced alarm notification resource within HuaweiCloud.

## Example Usage

```hcl
variable "topic_urn" {}

resource "huaweicloud_cnad_advanced_alarm_notification" "test" {
  topic_urn = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `topic_urn` - (Required, String) Specifies the topic urn of SMN. It is required that the SMN has been subscribed successfully.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `topic_urn`).

* `is_close_attack_source_flag` - Whether to enable the alarm content to shield the attack source information.

## Import

The CNAD advanced alarm notification can be imported using the `topic_urn`, e.g.

```bash
$ terraform import huaweicloud_cnad_advanced_alarm_notification.test <topic_urn>
```
