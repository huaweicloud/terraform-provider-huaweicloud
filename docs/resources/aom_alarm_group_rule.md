---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_group_rule"
description:  |-
  Manages an AOM alarm group rule resource within HuaweiCloud.
---

# huaweicloud_aom_alarm_group_rule

Manages an AOM alarm group rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "action_rule_name" {}
variable "enterprise_project_id" {}

resource "huaweicloud_aom_alarm_group_rule" "test" {
  name                  = var.name
  group_by              = ["resource_provider", "key-test", "resource_type"]
  group_interval        = 5
  group_repeat_waiting  = 0
  group_wait            = 0
  description           = "test"
  enterprise_project_id = var.enterprise_project_id

  detail {
    bind_notification_rule_ids = [var.action_rule_name]

    match {
      key     = "resource_type"
      operate = "EXIST"
    }

    match {
      key     = "resource_provider"
      operate = "EQUALS"
      value   = ["test"]
    }
  }

  detail {
    bind_notification_rule_ids = [var.action_rule_name]

    match {
      key     = "key-test"
      operate = "EQUALS"
      value   = ["value-test"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the alarm group rule name.
  Changing this creates a new resource.

* `detail` - (Required, List) Specifies the grouping conditions list.
  The [detail](#block--detail) structure is documented below.

* `group_by` - (Required, List) Specifies the combine notifications.

* `group_wait` - (Required, Int) Specifies the initial wait time.
  Value ranges from `0` to `600`. Unit is second.

* `group_interval` - (Required, Int) Specifies the batch processing interval.
  Value ranges from `5` to `1,800`. Unit is second.

* `group_repeat_waiting` - (Required, Int) Specifies the repeat interval.
  Value ranges from `0` to `1,296,000`. Unit is second.

* `description` - (Optional, String) Specifies the description of the rule.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the rule belongs.
  Changing this creates a new resource.

<a name="block--detail"></a>
The `detail` block supports:

* `bind_notification_rule_ids` - (Required, List) Specifies the action rule IDs.

* `match` - (Required, List) Specifies the matching conditions list.
  The [match](#block--detail--match) structure is documented below.

<a name="block--detail--match"></a>
The `match` block supports:

* `key` - (Required, String) Specifies the matching condition key.
  Valid value are as follows, or using specific key means taking tag as condition:
  + **event_severity**: event severity
  + **notification_scene**: notification scene
  + **resource_provider**: alarm source
  + **resource_type**: resource type

* `operate` - (Required, String) Specifies the matching condition operator. Valid values are **EQUALS**, **REGEX**, **EXIST**.

* `value` - (Optional, List) Specifies the matching condition value.
  + If `operate` is **EXIST**, it should be empty.
  + If `operate` is **REGEX**, it should be a regex expression.
  + If `operate` is **EQUALS**, it depends on `key`, can be as follows:
      - If `key` is **event_severity**, it can be **Critical**, **Major**, **Minor**, **Info**.
      - If `key` is **notification_scene**, it can be **notify_resolved**, **notify_triggered**.
      - If `key` is **resource_provider**, it can be specific alarm source.
      - If `key` is **resource_type**, it can be **cce-cluster**, **cluster**, **clusters-clustercert**, **clusters-nodepools**,
      **clusters-nodes**, **configmaps**, **deployments**, **ingresses**, **jobs**, **node**, **pods**,
      **podsecuritypolicies**, **releases**, **rolebindings**, **roles**, **routes**, **secrets**, **service**.
      - If `key` is specific key, it can be specific values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `name`.

* `created_at` - Indicates the rule create time.

* `updated_at` - Indicates the rule update time.

## Import

The rule can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_aom_alarm_group_rule.test <name>
```
