---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rule"
description: |-
  Manages a SecMaster alert rule resource within HuaweiCloud.
---

# huaweicloud_secmaster_alert_rule

Manages a SecMaster alert rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "pipeline_id" {}

resource "huaweicloud_secmaster_alert_rule" "test" {
  workspace_id = var.workspace_id
  pipeline_id  = var.pipeline_id
  name         = "test"
  description  = "this is a test rule created by terraform"
  status       = "ENABLED"
  severity     = "TIPS"
  type = {
    "name"     = "DNS protocol attacks"
    "category" = "DDoS attacks"
  }

  triggers {
    mode              = "COUNT"
    operator          = "GT"
    expression        = 5
    severity          = "MEDIUM"
    accumulated_times = 1
  }

  query_rule = "* | select status, count(*) as count group by status"
  query_type = "SQL"
  query_plan {
    query_interval      = 1
    query_interval_unit = "HOUR"
    time_window         = 1
    time_window_unit    = "HOUR"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the alert rule belongs.

  Changing this parameter will create a new resource.

* `pipeline_id` - (Required, String, ForceNew) Specifies the pipeline ID of the alert rule.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the alert rule name.

* `severity` - (Required, String) Specifies the severity of the alert rule.
  The value can be: **TIPS**, **LOW**, **MEDIUM**, **HIGH** and **FATAL**.

* `type` - (Required, Map) Specifies the type of the alert rule.

* `description` - (Required, String) Specifies the description of the alert rule.

* `status` - (Required, String) Specifies the status of the alert rule.
  The value can be **ENABLED** and **DISABLED**. Defaults to **ENABLED**.

* `query_rule` - (Required, String) Specifies the query rule of the alert rule.

* `query_type` - (Required, String) Specifies the query type of the alert rule.
  The value can be: **SQL**.

* `query_plan` - (Required, List) Specifies the query plan of the alert rule.
  The [query_plan](#query_plan) structure is documented below.

* `triggers` - (Required, List) Specifies the triggers of the alert rule.
  The [triggers](#triggers) structure is documented below.

* `custom_information` - (Optional, Map) Specifies the custom information of the alert rule.

* `event_grouping` - (Optional, Bool) Specifies whether to put events in a group.
  The value can be:
  + **true**: one alarm for all query results;
  + **false**: one alarm for each query result;

  Default to **true**.

* `debugging_alarm` - (Optional, Bool) Specifies whether to generate debugging alarms.
  Defaults to **true**.

* `suppression` - (Optional, Bool) Specifies whether to stop the query when an alarm is generated.

<a name="query_plan"></a>
The `query_plan` block supports:

* `query_interval` - (Required, Int) Specifies the query interval.
  + When `query_interval_unit` is **MINUTE**: the value range is `5` to `59`;
  + When `query_interval_unit` is **HOUR**: the value range is `1` to `23`;
  + When `query_interval_unit` is **DAY**: the value range is `1` to `14`;

* `query_interval_unit` - (Required, String) Specifies the query interval unit.
  The value can be: **MINUTE**, **HOUR** and **DAY**.

* `time_window` - (Required, Int) Specifies the time window.
  + When `time_window_unit` is **MINUTE**: the value range is `5` to `59`;
  + When `time_window_unit` is **HOUR**: the value range is `1` to `23`;
  + When `time_window_unit` is **DAY**: the value range is `1` to `14`;

* `time_window_unit` - (Required, String) Specifies the time window unit.
  The value can be: **MINUTE**, **HOUR** and **DAY**.

* `execution_delay` - (Optional, Int) Specifies the execution delay in minutes.

* `overtime_interval` - (Optional, Int) Specifies the overtime interval in minutes.

<a name="triggers"></a>
The `triggers` block supports:

* `expression` - (Required, String) Specifies the expression.

* `operator` - (Required, String) Specifies the operator.
  The value can be: **EQ**(equal), **NE**(not equal), **GT**(greater than) and **LT**(less than).

* `accumulated_times` - (Required, Int) Specifies the accumulated times.

* `mode` - (Required, String) Specifies the trigger mode.
  The value can be: **COUNT**.

* `severity` - (Required, String) Specifies the severity of the trigger.
  The value can be: **TIPS**, **LOW**, **MEDIUM**, **HIGH** and **FATAL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The created time.

* `updated_at` - The updated time.

## Import

The alert rule can be imported using theworkspace ID and the alert rule, e.g.

```bash
$ terraform import huaweicloud_secmaster_alert_rule.test <workspace_id>/<id>
```
