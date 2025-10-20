---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_immutable_tag_rule"
description: |-
  Manages a SWR enterprise instance immutable tag rule resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_immutable_tag_rule

Manages a SWR enterprise instance immutable tag rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_swr_enterprise_immutable_tag_rule" "test" {
  instance_id    = var.instance_id
  namespace_name = "library"

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

  priority = 0
  action   = "immutable"
  disabled = false
  template = "immutable_template"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String, NonUpdatable) Specifies the namespace name.

* `scope_selectors` - (Required, List) Specifies the repository selectors.
  The [scope_selectors](#block--scope_selectors) structure is documented below.

* `tag_selectors` - (Required, List) Specifies the repository version selector.
  The [tag_selectors](#block--tag_selectors) structure is documented below.

* `action` - (Optional, String) Specifies the policy action. Valid value is **immutable**.

* `disabled` - (Optional, Bool) Specifies whether the policy rule is disabled. Default to **false**.

* `priority` - (Optional, Int) Specifies the priority. Default to **0**.

* `template` - (Optional, String) Specifies the template type. Valid value is **immutable_template**.

<a name="block--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - (Required, String) Specifies the repository selector key. Valid value is **repository**.

* `value` - (Required, List) Specifies the repository selector value.
  The [value](#block--scope_selectors--value) structure is documented below.

<a name="block--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - (Required, String) Specifies the selector matching type. Valid value is **repoMatches**.

* `kind` - (Required, String) Specifies the matching type. Valid value is **doublestar**.

* `pattern` - (Required, String) Specifies the pattern.

* `extras` - (Optional, String) Specifies the extra infos.

<a name="block--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - (Required, String) Specifies the selector matching type. Valid value is **repoMatches**.

* `kind` - (Required, String) Specifies the matching type. Valid value is **doublestar**.

* `pattern` - (Required, String) Specifies the pattern.

* `extras` - (Optional, String) Specifies the extra infos.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `immutable_rule_id` - Indicates the policy ID

* `namespace_id` - Indicates the namespace ID

## Import

The immutable tag rule can be imported using `instance_id`, `namespace_name` and `immutable_rule_id` separated by
slashes, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_immutable_tag_rule.test <instance_id>/<namespace_name>/<immutable_rule_id>
```
