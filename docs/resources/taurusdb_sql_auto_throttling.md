---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_sql_auto_throttling"
description: |-
  Manages a TaurusDB SQL auto throttling resource within HuaweiCloud.
---

# huaweicloud_taurusdb_sql_auto_throttling

Manages a TaurusDB SQL auto throttling resource within HuaweiCloud.

-> **NOTE:** This feature is only available for primary nodes. Read-only nodes are not supported.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_taurusdb_sql_auto_throttling" "test" {
  instance_id     = var.instance_id
  node_id         = var.node_id
  start_time      = "00:00"
  end_time        = "01:00"
  condition       = "and"
  cpu_usage       = 70
  active_sessions = 3
  clear_time      = 3
  duration        = 2
  max_concurrency = 1000
  retain_sql_rule = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the ID of the TaurusDB instance.

* `node_id` - (Required, String, NoneUpdatable) Specifies the ID of the node. The node role must be the primary node.

* `start_time` - (Required, String) Specifies the start time of the throttling time window.
  The format is **hh:mm**, e.g. **01:00**.

* `end_time` - (Required, String) Specifies the end time of the throttling time window.
  The format is **hh:mm**, e.g. **23:59**.

* `condition` - (Required, String) Specifies the relationship between CPU usage and active sessions conditions.
  Valid values are:
  + **and**: Both conditions must be met.
  + **or**: Either condition is met.

* `cpu_usage` - (Required, Int) Specifies the CPU usage threshold used to determine the instance load.
  When the CPU usage exceeds the set value (>=70%), it will be combined with the active sessions condition
  to determine whether to trigger autonomous throttling.
  The value ranges from `70` to `100`.

* `active_sessions` - (Required, Int) Specifies the active sessions threshold used to determine the concurrent
  access volume of the instance. When the active sessions exceed the set value, it will be combined with the
  CPU usage condition to determine whether to trigger autonomous throttling.
  The value ranges from `1` to `5,000`.

* `clear_time` - (Required, Int) Specifies the maximum throttling duration each time, in minutes.
  When the maximum throttling duration within the throttling time window is exceeded, the system will
  automatically exit throttling. If the trigger conditions remain unchanged, throttling will be triggered
  again after a 1-minute interval.
  The value ranges from `1` to `1440`.

* `duration` - (Required, Int) Specifies the duration for which the limiting conditions are met, in minutes.
  This sets an observation period to ensure that autonomous throttling is triggered only after the CPU usage
  and active sessions conditions have been continuously met for the specified duration.
  The value ranges from `2` to `60`.

* `max_concurrency` - (Required, Int) Specifies the maximum concurrency used to limit the number of concurrent
  executions of SQL statements that match the keywords. When the concurrency exceeds the set value, excess
  SQL statements will be rejected.
  The value ranges from `0` to `1,000,000,000`.

* `retain_sql_rule` - (Required, String) Specifies whether to retain existing SQL limiting rules.
  Valid values are as follows:
  + **true**: Existing SQL throttling rules are retained.
  + **false**: Existing SQL throttling rules are cleared.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<node_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SQL auto throttling resource can be imported using the `instance_id` and `node_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_taurusdb_sql_auto_throttling.test <instance_id>/<node_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `retain_sql_rule`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the auto SQL limiting, or the resource definition should be updated
to align with the auto SQL limiting. Also you can ignore changes as below.

```hcl
resource "huaweicloud_taurusdb_sql_auto_throttling" "test" {
  ...

  lifecycle {
    ignore_changes = [
      retain_sql_rule,
    ]
  }
}
```
