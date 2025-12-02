---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_notification_rule"
description: |-
  Use this resource to manage notification rules for metrics within HuaweiCloud.
---

# huaweicloud_workspace_notification_rule

Use this resource to manage notification rules for metrics within HuaweiCloud.

-> You need to add delegate authorization for the cloud service first, otherwise notifications cannot be sent to SMN
   normally.

## Example Usage

```hcl
variable "metric_name" {}
variable "notified_topic_urn" {}

resource "huaweicloud_workspace_notification_rule" "test" {
  metric_name         = var.metric_name
  comparison_operator = ">="
  enable              = true
  notify_object       = var.notified_topic_urn
  threshold           = 30
  interval            = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the notification rule is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `metric_name` - (Required, String, ForceNew) Specifies the name of the metric.  
  The valid values are as follows:
  + **desktop_idle_duration**

* `comparison_operator` - (Required, String) Specifies the comparison operator for the metric value and threshold.  
  The valid values are as follows:
  + **>=** - Trigger when the metric value is greater than or equal to the threshold.

* `enable` - (Required, Bool) Specifies whether to enable the rule.

* `notify_object` - (Required, String) Specifies the notification object, which is the SMN topic URN.

* `threshold` - (Optional, Int) Specifies the threshold (in days) for the rule configuration.

* `interval` - (Optional, Int) Specifies the interval time (in days) for the next notification after triggering.  
  Default is once per day.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The notification rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_notification_rule.test <id>
```
