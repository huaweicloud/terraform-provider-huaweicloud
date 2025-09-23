---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_alarm_notifications"
description: |-
  Use this data source to get the list of CNAD advanced alarm notifications.
---

# huaweicloud_cnad_advanced_alarm_notifications

Use this data source to get the list of CNAD advanced alarm notifications.

## Example Usage

```hcl
data "huaweicloud_cnad_advanced_alarm_notifications" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topic_urn` - The topic urn of SMN. Empty value means no alarm notifications is configured.

* `is_close_attack_source_flag` - Whether to enable the alarm content to shield the attack source information.
