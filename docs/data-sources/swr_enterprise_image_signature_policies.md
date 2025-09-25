---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_image_signature_policies"
description: |-
  Use this data source to get the list of SWR enterprise image signature policies.
---

# huaweicloud_swr_enterprise_image_signature_policies

Use this data source to get the list of SWR enterprise image signature policies.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_image_signature_policies" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - Indicates the image signature policies.
  The [policies](#attrblock--policies) structure is documented below.

<a name="attrblock--policies"></a>
The `policies` block supports:

* `id` - Indicates the policy ID.

* `namespace_name` - Indicates the namespace name.

* `name` - Indicates the policy name.

* `signature_algorithm` - Indicates the signature algorithm.

* `signature_key` - Indicates the signature key.

* `signature_method` - Indicates the signature method.

* `namespace_id` - Indicates the namespace ID

* `description` - Indicates the description of policy.

* `enabled` - Indicates whether the policy is enabled.

* `scope_rules` - Indicates the scope rules
  The [scope_rules](#attrblock--policies--scope_rules) structure is documented below.

* `trigger` - Indicates the trigger config.
  The [trigger](#attrblock--policies--trigger) structure is documented below.

* `creator` - Indicates the creator

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

<a name="attrblock--policies--scope_rules"></a>
The `scope_rules` block supports:

* `repo_scope_mode` - Indicates the repository select method.

* `scope_selectors` - Indicates the repository selectors.
  The [scope_selectors](#attrblock--policies--scope_rules--scope_selectors) structure is documented below.

* `tag_selectors` - Indicates the repository version selector.
  The [rule_selector](#attrblock--policies--scope_rules--rule_selector) structure is documented below.

<a name="attrblock--policies--scope_rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - Indicates the repository selector key.

* `value` - Indicates the repository selector value.
  The [rule_selector](#attrblock--policies--scope_rules--rule_selector) structure is documented below.

<a name="attrblock--policies--scope_rules--rule_selector"></a>
The `rule_selector` block supports:

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
