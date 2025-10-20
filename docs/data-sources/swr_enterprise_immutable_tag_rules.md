---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_immutable_tag_rules"
description: |-
  Use this data source to get the list of SWR enterprise instance immutable tag rules.
---

# huaweicloud_swr_enterprise_immutable_tag_rules

Use this data source to get the list of SWR enterprise instance immutable tag rules.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_immutable_tag_rules" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_id` - (Optional, String) Specifies the enterprise instance namespace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - Indicates the immutable tag rules.
  The [rules](#attrblock--rules) structure is documented below.

<a name="attrblock--rules"></a>
The `rules` block supports:

* `id` - Indicates the immutable tag rule ID

* `namespace_id` - Indicates the namespace ID

* `namespace_name` - Indicates the namespace name.

* `action` - Indicates the policy action.

* `disabled` - Indicates whether the policy rule is disabled.

* `priority` - Indicates the priority.

* `scope_selectors` - Indicates the repository selectors.
  The [scope_selectors](#attrblock--rules--scope_selectors) structure is documented below.

* `tag_selectors` - Indicates the repository version selector.
  The [tag_selectors](#attrblock--rules--tag_selectors) structure is documented below.

* `template` - Indicates the template type.

<a name="attrblock--rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - Indicates the repository selector key.

* `value` - Indicates the repository selector value.
  The [value](#attrblock--rules--scope_selectors--value) structure is documented below.

<a name="attrblock--rules--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - Indicates the selector matching type.

* `extras` - Indicates the extra infos.

* `kind` - Indicates the matching type.

* `pattern` - Indicates the pattern.

<a name="attrblock--rules--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - Indicates the selector matching type.

* `extras` - Indicates the extra infos.

* `kind` - Indicates the matching type.

* `pattern` - Indicates the pattern.
