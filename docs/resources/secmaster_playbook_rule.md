---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_rule"
description: |-
  Manages a SecMaster playbook rule resource within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_rule

Manages a SecMaster playbook rule resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
variable "workspace_id" {}
variable "version_id" {}

resource "huaweicloud_secmaster_playbook_rule" "test" {
  workspace_id    = var.workspace_id
  version_id      = var.version_id
  expression_type = "common"

  conditions {
    name   = "condition_0"
    detail = "123"
    data   = ["waf.alarm.level", ">", "3"]
  }

  logics = ["condition_0"]
}
```

### More Examples

For more detailed associated usage see [playbook instructions](/examples/secmaster/playbook/README.md)

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the playbook rule belongs.

  Changing this parameter will create a new resource.

* `version_id` - (Required, String, ForceNew) Specifies playbook version ID of the rule.

  Changing this parameter will create a new resource.

* `expression_type` - (Optional, String) Specifies the expression type of the rule.
  Required for event triggered playbooks.

* `conditions` - (Optional, List) Specifies the conditions of the rule.
  Required for event triggered playbooks.
The [conditions](#PlaybookRule_ConditionItem) structure is documented below.

* `logics` - (Optional, List) Specifies the logics of the rule.
  Required for event triggered playbooks.

* `cron` - (Optional, String) Specifies the cron expression.
  Required for timer triggered playbooks.

* `schedule_type` - (Optional, String) Specifies the schedule type.
  The value can be **second**, **hour** and **day** **week**. Required for timer triggered playbooks.

* `start_type` - (Optional, String) Specifies the playbook start type.
  The value can be: **IMMEDIATELY** and **CUSTOM**. Required for timer triggered playbooks.

* `end_type` - (Optional, String) Specifies the playbook end type.
  The value can be: **FOREVER** and **CUSTOM**. Required for timer triggered playbooks.

* `end_time` - (Optional, String) Specifies the playbook end time.
  Required for timer triggered playbooks.

* `repeat_range` - (Optional, String) Specifies the repeat range.
  Required for timer triggered playbooks.

* `only_once` - (Optional, Bool) Specifies the repeat range.
  Required for timer triggered playbooks.

* `execution_type` - (Optional, String) Specifies the execution type.
  The value can be **PARALLEL**. Required for timer triggered playbooks.

<a name="PlaybookRule_ConditionItem"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the condition name.

* `detail` - (Optional, String) Specifies the condition detail.

* `data` - (Optional, List) Specifies the condition data.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the created time of the playbook rule.

* `updated_at` - Indicates the updated time of the playbook rule.

## Import

The playbook rule can be imported using  workspace ID, the playbook version ID and the playbook rule ID,
separated by slashes, e.g.

```bash
$ terraform import huaweicloud_secmaster_playbook_rule.test <workspace_id>/<playbook_version_id>/<playbook_rule_id>
```
