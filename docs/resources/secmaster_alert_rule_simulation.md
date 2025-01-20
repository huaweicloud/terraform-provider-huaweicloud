---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rule_simulation"
description: |-
  Manages a SecMaster alert rule simulation resource within HuaweiCloud.
---

# huaweicloud_secmaster_alert_rule_simulation

Manages a SecMaster alert rule simulation resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "pipeline_id" {}
variable "query_rule" {}

resource "huaweicloud_secmaster_alert_rule_simulation" "test" {
  workspace_id   = var.workspace_id
  pipeline_id    = var.pipeline_id
  query_rule     = var.query_rule
  query_type     = "SQL"
  from_time      = "2025-01-16 19:04:05"
  to_time        = "2025-01-16 19:06:05"
  event_grouping = false

  triggers {
    mode              = "COUNT"
    operator          = "GT"
    expression        = 5
    severity          = "MEDIUM"
    accumulated_times = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the alert rule belongs.

* `pipeline_id` - (Required, String, NonUpdatable) Specifies the pipeline ID of the alert rule.

* `query_rule` - (Required, String, NonUpdatable) Specifies the query rule of the alert rule.

* `query_type` - (Required, String, NonUpdatable) Specifies the query type of the alert rule.
  The value can be: **SQL**.

* `from_time` - (Required, String, NonUpdatable) Specifies the start time of the alert rule simulation.

* `to_time` - (Required, String, NonUpdatable) Specifies the end time of the alert rule simulation.

* `event_grouping` - (Optional, Bool, NonUpdatable) Specifies whether to put events in a group.
  The value can be:
  + **true**: one alarm for all query results;
  + **false**: one alarm for each query result;

  Default to **true**.

* `triggers` - (Required, List, NonUpdatable) Specifies the triggers of the alert rule.
  The [triggers](#triggers) structure is documented below.

<a name="triggers"></a>
The `triggers` block supports:

* `expression` - (Required, String, NonUpdatable) Specifies the expression.

* `operator` - (Required, String, NonUpdatable) Specifies the operator.
  The value can be: **EQ**(equal), **NE**(not equal), **GT**(greater than) and **LT**(less than).

* `accumulated_times` - (Required, Int, NonUpdatable) Specifies the accumulated times.

* `mode` - (Required, String, NonUpdatable) Specifies the trigger mode.
  The value can be: **COUNT**.

* `severity` - (Required, String, NonUpdatable) Specifies the severity of the trigger.
  The value can be: **TIPS**, **LOW**, **MEDIUM**, **HIGH** and **FATAL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
