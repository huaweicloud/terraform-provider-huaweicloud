---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_retention_policy"
description: |-
  Manages a SWR enterprise instance retention policy resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_retention_policy

Manages a SWR enterprise instance retention policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "name" {}

resource "huaweicloud_swr_enterprise_retention_policy" "test"{
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  name           = var.name
  algorithm      = "or"
  enabled        = true
  
  rules {
    priority        = 0
    action          = "retain"
    repo_scope_mode = "regular"
    disabled        = false
    template        = "latestPushedK"

    params = {
      latestPushedK = jsonencode(1)
    }

    scope_selectors {
      key = "repository"

      value {
        kind       = "doublestar"
        decoration = "repoMatches"
        pattern    = "**"
      }
    }
    
    tag_selectors {
      kind       = "doublestar"
      decoration = "matches"
      pattern    = "**"
    }
  }

  trigger {
    type = "manual"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String, NonUpdatable) Specifies the namespace name.

* `name` - (Required, String) Specifies the policy name.

* `algorithm` - (Required, String) Specifies the algorithm of policy. Valid value is **or**.

* `enabled` - (Required, Bool) Specifies whether the policy is enabled.

* `rules` - (Required, List) Specifies the retention rules.
  The [rules](#block--rules) structure is documented below.

* `trigger` - (Required, List) Specifies the trigger config.
  The [trigger](#block--trigger) structure is documented below.

<a name="block--rules"></a>
The `rules` block supports:

* `action` - (Required, String) Specifies the policy action. Valid value is **retain**.

* `priority` - (Required, Int) Specifies the priority. Valid value is **0**.

* `repo_scope_mode` - (Required, String) Specifies the repo scope mode. Valid values are **regular** and **selection**.

* `template` - (Required, String) Specifies the template type.
  Valid values are **latestPulledN**, **latestPushedK**, **nDaysSinceLastPush**, **nDaysSinceLastPull**.

* `params` - (Required, Map) Specifies the params. Value should be in JSON format.

* `scope_selectors` - (Required, List) Specifies the repository selectors.
  The [scope_selectors](#block--rules--scope_selectors) structure is documented below.

* `tag_selectors` - (Required, List) Specifies the repository version selector.
  The [tag_selectors](#block--rules--tag_selectors) structure is documented below.

* `disabled` - (Optional, Bool) Specifies whether the policy rule is disabled. Default to **false**.

<a name="block--rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - (Required, String) Specifies the repository selector key. Valid value is **repository**.

* `value` - (Required, List) Specifies the repository selector value.
  The [value](#block--rules--scope_selectors--value) structure is documented below.

<a name="block--rules--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - (Required, String) Specifies the selector matching type.

* `kind` - (Required, String) Specifies the matching type.

* `pattern` - (Required, String) Specifies the pattern.

* `extras` - (Optional, String) Specifies the extra infos.

<a name="block--rules--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - (Required, String) Specifies the selector matching type.

* `kind` - (Required, String) Specifies the matching type.

* `pattern` - (Required, String) Specifies the pattern.

* `extras` - (Optional, String) Specifies the extra infos.

<a name="block--trigger"></a>
The `trigger` block supports:

* `type` - (Required, String) Specifies the trigger type. Valid values are **manual** and **scheduled**.

* `trigger_settings` - (Optional, List) Specifies the trigger settings.
  The [trigger_settings](#block--trigger--trigger_settings) structure is documented below.

<a name="block--trigger--trigger_settings"></a>
The `trigger_settings` block supports:

* `cron` - (Optional, String) Specifies the scheduled setting.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `namespace_id` - Indicates the namespace ID

* `policy_id` - Indicates the policy ID

* `rules` - Indicates the retention rules.
  The [rules](#attrblock--rules) structure is documented below.

<a name="attrblock--rules"></a>
The `rules` block supports:

* `id` - Indicates the retention policy rule ID.

## Import

The retention policy can be imported using `instance_id`, `namespace_name` and `policy_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_retention_policy.test <instance_id>/<namespace_name>/<policy_id>
```
