---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_inhibit_rule"
description:  |-
  Manages an AOM alarm inhibit rule resource within HuaweiCloud.
---

# huaweicloud_aom_alarm_inhibit_rule

Manages an AOM alarm inhibit rule resource within HuaweiCloud.

## Example Usage

### Create a basic alarm inhibit rule with source and target matches

```hcl
variable "inhibit_rule_name" {}
variable "source_matches" {
  type = list(object({
    conditions = list(object({
      key     = string
      operate = string
      values  = optional(list(string), null)
    }))
  }))
}
variable "target_matches" {
  type = list(object({
    conditions = list(object({
      key     = string
      operate = string
      values  = optional(list(string), null)
    }))
  }))
}

resource "huaweicloud_aom_alarm_inhibit_rule" "test" {
  name = var.inhibit_rule_name

  dynamic "source_matches" {
    for_each = var.source_matches

    content {
      dynamic "conditions" {
        for_each = source_matches.value.conditions

        content {
          key     = conditions.value.key
          operate = conditions.value.operate
          values  = conditions.value.values
        }
      }
    }
  }

  dynamic "target_matches" {
    for_each = var.target_matches

    content {
      dynamic "conditions" {
        for_each = target_matches.value.conditions

        content {
          key     = conditions.value.key
          operate = conditions.value.operate
          values  = conditions.value.values
        }
      }
    }
  }
}
```

### Create an orchestrated alarm inhibit rule with match_v3

```hcl
variable "inhibit_rule_name" {}
variable "alarm_group_rule_name" {}
variable "match_v3" {}

resource "huaweicloud_aom_alarm_inhibit_rule" "test" {
  name               = var.inhibit_rule_name
  binding_group_rule = var.alarm_group_rule_name
  match_v3           = var.match_v3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the alarm inhibit rule is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the alarm inhibit rule.
  The name can contain `1` to `100` characters, only letters, digits and underscores (_) are allowed.  
  The name cannot start or end with an underscore.

* `description` - (Optional, String) Specifies the description of the alarm inhibit rule.  
  The description contains a maximum of `1,024` characters.

* `binding_group_rule` - (Optional, String) Specifies the rule name associated with the alarm inhibit rule.  
  For details, please refer to the [documentation](https://support.huaweicloud.com/api-aom/AddInhibitRule.html#AddInhibitRule__request_InhibitMatchV3Tag).

  This parameter is available only when orchestrated alarm inhibit rule is supported.

* `match_v3` - (Optional, String) Specifies the orchestrated alarm inhibit rule definition, in JSON format.  
  This parameter is **required** when orchestrated alarm inhibit rule is supported.

* `source_matches` - (Optional, List) Specifies the parallel match conditions for root alerts that suppress
  other alerts.  
  The [source_matches](#alarm_inhibit_rule_matches) structure is documented below.  
  This parameter is **required** when orchestrated alarm inhibit rule is not supported.  
  The maximum number of elements is `10`.

* `target_matches` - (Optional, List) Specifies the parallel match conditions for target alerts that
  will be suppressed.  
  The [target_matches](#alarm_inhibit_rule_matches) structure is documented below.
  This parameter is **required** when orchestrated alarm inhibit rule is not supported.  
  The maximum number of elements is `10`.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project to which
  the alarm inhibit rule belongs.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

<a name="alarm_inhibit_rule_matches"></a>
The `source_matches` and `target_matches` block supports:

* `conditions` - (Required, List) Specifies the serial conditions within a parallel condition group.  
  The [conditions](#alarm_inhibit_rule_matches_conditions) structure is documented below.  
  The maximum number of elements is `10`.

<a name="alarm_inhibit_rule_matches_conditions"></a>
The `conditions` block supports:

* `key` - (Required, String) Specifies the key of the alarm.  
  The valid values are as follows:
  + Specified tag name: The tag name can only contain Chinese characters, letters, numbers and underscores (_).
  + **event_severity**: Event severity.
  + **resource_provider**: Alarm source.
  + **resource_type**: Resource type.

* `operate` - (Required, String) Specifies the match operator for the alarm key.  
  The valid values are as follows:
  + **EQUALS**
  + **EXIST**
  + **REGEX**

* `values` - (Optional, List) Specifies the value list corresponding to the key of the alarm.  
  Each value cannot exceed `256` characters when the `operate` is **REGEX**.  
  + If `key` is **event_severity** and the `operate` is **EQUALS**, the valid values are **Critical**, **Major**,
    **Minor** and **Info**.
  + If `key` is **resource_provider**, the value can be any of the resource types specified when creating an alarm rule
    or customizing an alarm report. Types can include **host**, **container**, **process**, etc.
  + If `key` is **resource_type**, the value can be any of the service names that triggered the alarm or event.
    The service names can include **AOM**, **LTS**, **CCE**, etc.
  + If `key` is a tag, the value can be any of the tag values ​​corresponding to the tag name. Tag values ​​can
    only contain Chinese characters, letters, numbers, and underscores (_).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the alarm inhibit rule name.

## Import

The alarm inhibit rule can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_aom_alarm_inhibit_rule.test <id>
```

For the alarm inhibit rule with the `enterprise_project_id`, its enterprise project ID need to be specified
additionanlly when importing, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_aom_alarm_inhibit_rule.test <id>/<enterprise_project_id>
```
