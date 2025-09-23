---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_notification_mask"
description: |-
  Manages a CES notification mask resource within HuaweiCloud.
---

# huaweicloud_ces_notification_mask

Manages a CES notification mask resource within HuaweiCloud.

## Example Usage

### Masked By Resource

```hcl
variable "mask_name" {}

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "RESOURCE"
  mask_name     = var.mask_name
  
  resources {
    namespace = "SYS.OBS"

    dimensions {
      name  = "bucket_name"
      value = "bucket-one"
    }
  }

  resources {
    namespace = "SYS.OBS"
    
    dimensions {
      name  = "bucket_name"
      value = "bucket-two"
    }
  }

  mask_type  = "START_END_TIME"
  start_date = "2025-03-27"
  start_time = "09:12:09"
  end_date   = "2025-03-27"
  end_time   = "20:12:09"
}
```

### Masked By Policy

```hcl
variable "mask_name" {}
variable "alarm_policy_id" {}

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "RESOURCE_POLICY_NOTIFICATION"
  mask_name     = var.mask_name
  relation_ids  = [var.alarm_policy_id]
  
  resources {
    namespace = "SYS.OBS"

    dimensions {
      name = "bucket_name"
      value = "*"
    }
  }

  mask_type  = "FOREVER_TIME"
}
```

### Masked By Rule

```hcl
variable "alarm_rule" {}

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "ALARM_RULE"
  relation_id   = var.alarm_rule

  mask_type  = "START_END_TIME"
  start_date = "2025-03-27"
  start_time = "09:12:09"
  end_date   = "2025-03-27"
  end_time   = "20:12:10"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `relation_type` - (Required, String, NonUpdatable) Specifies the type of a resource that is associated with an alarm notification
  masking rule.
  The valid values are as follows:
  + **ALARM_RULE**: alarm rules;
  + **RESOURCE**: resources;
  + **RESOURCE_POLICY_NOTIFICATION**: alarm policies for the resource;

* `mask_type` - (Required, String) Specifies the alarm notification masking type.

* `relation_ids` - (Optional, List) Specifies the alarm policy IDs.

* `relation_id` - (Optional, String) Specifies the alarm rule ID.

* `mask_name` - (Optional, String) Specifies the masking rule name.

* `start_date` - (Optional, String) Specifies the masking start date, in **yyyy-MM-dd** format.

* `start_time` - (Optional, String) Specifies the masking start time, in **HH:mm:ss** format.

* `end_date` - (Optional, String) Specifies the masking end date, in **yyyy-MM-dd** format.

* `end_time` - (Optional, String) Specifies the masking end time, in **HH:mm:ss** format.

* `resources` - (Optional, List) Specifies the resource for which alarm notifications will be masked.
  The [resources](#Resources) structure is documented below.

<a name="Resources"></a>
The `resources` block supports:

* `namespace` - (Required, String) Specifies the resource namespace in **service.item** format.

* `dimensions` - (Required, List) Specifies the resource dimension information.
  The [dimensions](#ResourcesDimensions) structure is documented below.

<a name="ResourcesDimensions"></a>
The `dimensions` block supports:

* `name` - (Required, String) Specifies the dimension of a resource.

* `value` - (Required, String) Specifies the value of a resource dimension.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `mask_status` - Specifies the alarm notification masking status.

* `policies` - The alarm policy list.
  The [policies](#Policies) structure is documented below.

<a name="Policies"></a>
The `policies` block supports:

* `alarm_policy_id` - The alarm policy ID.

* `metric_name` - The metric name of a resource.

* `extra_info` - The extended metric information.
  The [extra_info](#PoliciesExtraInfo) structure is documented below.

* `period` - The period for determining whether to generate an alarm, in seconds.

* `filter` - The data rollup method.

* `comparison_operator` - The operator.

* `value` - The alarm threshold.

* `unit` - The data unit.

* `count` - The number of consecutive times that alarm conditions are met.

* `type` - The alarm policy type.

* `suppress_duration` - The interval for triggering alarms.

* `alarm_level` - The alarm severity.

* `selected_unit` - The unit you selected, which is used for subsequent metric data display and calculation.

<a name="PoliciesExtraInfo"></a>
The `extra_info` block supports:

* `origin_metric_name` - The original metric name.

* `metric_prefix` - The metric name prefix.

* `custom_proc_name` - The name of a user process.

* `metric_type` - The metric type.

## Import

When `relation_type` value is `RESOURCE` or `RESOURCE_POLICY_NOTIFICATION`, the notification mask can be imported
using `relation_type` and `id`, e.g.

```bash
$ terraform import huaweicloud_ces_notification_mask.test <relation_type>/<id>
```

When `relation_type` value is `ALARM_RULE`, the notification mask can be imported using `relation_type` and
`relation_id`, e.g.

```bash
$ terraform import huaweicloud_ces_notification_mask.test <relation_type>/<relation_id>
```
