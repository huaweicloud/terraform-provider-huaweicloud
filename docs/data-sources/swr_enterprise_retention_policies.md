---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_retention_policies"
description: |-
  Use this data source to get the list of SWR enterprise instance retention policies.
---

# huaweicloud_swr_enterprise_retention_policies

Use this data source to get the list of SWR enterprise instance retention policies.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_retention_policies" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `name` - (Optional, String) Specifies the trigger name.

* `namespace_id` - (Optional, String) Specifies the namespace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - Indicates the policies.
  The [policies](#attrblock--policies) structure is documented below.

<a name="attrblock--policies"></a>
The `policies` block supports:

* `id` - Indicates the policy ID.

* `name` - Indicates the policy name.

* `namespace_id` - Indicates the namespace ID

* `namespace_name` - Indicates the namespace name.

* `algorithm` - Indicates the algorithm of policy.

* `enabled` - Indicates whether the policy is enabled.

* `rules` - Indicates the retention rules.
  The [rules](#attrblock--policies--rules) structure is documented below.

* `trigger` - Indicates the trigger config.
  The [trigger](#attrblock--policies--trigger) structure is documented below.

<a name="attrblock--policies--rules"></a>
The `rules` block supports:

* `id` - Indicates the retention policy rule ID.

* `action` - Indicates the policy action.

* `disabled` - Indicates whether the policy rule is disabled.

* `params` - Indicates the params.

* `priority` - Indicates the priority.

* `repo_scope_mode` - Indicates the repo scope mode.

* `scope_selectors` - Indicates the repository selectors.
  The [scope_selectors](#attrblock--policies--rules--scope_selectors) structure is documented below.

* `tag_selectors` - Indicates the repository version selector.
  The [tag_selectors](#attrblock--policies--rules--tag_selectors) structure is documented below.

* `template` - Indicates the template type.

<a name="attrblock--policies--rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - Indicates the repository selector key.

* `value` - Indicates the repository selector value.
  The [value](#attrblock--policies--rules--scope_selectors--value) structure is documented below.

<a name="attrblock--policies--rules--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - Indicates the selector matching type.

* `extras` - Indicates the extra infos.

* `kind` - Indicates the matching type.

* `pattern` - Indicates the pattern.

<a name="attrblock--policies--rules--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - Indicates the selector matching type.

* `extras` - Indicates the extra infos.

* `kind` - Indicates the matching type.

* `pattern` - Indicates the pattern.

<a name="attrblock--policies--trigger"></a>
The `trigger` block supports:

* `trigger_settings` - Indicates the trigger settings.
  The [trigger_settings](#attrblock--policies--trigger--trigger_settings) structure is documented below.

* `type` - Indicates the trigger type.

<a name="attrblock--policies--trigger--trigger_settings"></a>
The `trigger_settings` block supports:

* `cron` - Indicates the scheduled setting.
