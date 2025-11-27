---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_triggers"
description: |-
  Use this data source to get the list of SWR enterprise instance triggers.
---

# huaweicloud_swr_enterprise_triggers

Use this data source to get the list of SWR enterprise instance triggers.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_triggers" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_id` - (Optional, String) Specifies the namespace ID.

* `name` - (Optional, String) Specifies the trigger name.

* `order_column` - (Optional, String) Specifies the order column.
  Value can be **desc** or **asc**, should use with `order_column`. Default to **desc**.

* `order_type` - (Optional, String) Specifies the order type.
  Value can be **created_at**, **updated_at** and **name**. Default to **created_at**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `triggers` - Indicates the triggers.
  The [triggers](#attrblock--triggers) structure is documented below.

<a name="attrblock--triggers"></a>
The `triggers` block supports:

* `id` - Indicates the trigger ID.

* `name` - Indicates the trigger name.

* `namespace_id` - Indicates the namespace ID

* `namespace_name` - Indicates the namespace name.

* `description` - Indicates the description of trigger.

* `enabled` - Indicates whether the trigger is enabled.

* `event_types` - Indicates the event types of trigger.

* `scope_rules` - Indicates the scope rules
  The [scope_rules](#attrblock--triggers--scope_rules) structure is documented below.

* `targets` - Indicates the target params.
  The [targets](#attrblock--triggers--targets) structure is documented below.

* `updated_at` - Indicates the last update time.

* `created_at` - Indicates the creation time.

* `creator` - Indicates the creator

<a name="attrblock--triggers--scope_rules"></a>
The `scope_rules` block supports:

* `repo_scope_mode` - Indicates the repository select method.

* `scope_selectors` - Indicates the repository selectors.
  The [scope_selectors](#attrblock--triggers--scope_rules--scope_selectors) structure is documented below.

* `tag_selectors` - Indicates the repository version selector.
  The [tag_selectors](#attrblock--triggers--scope_rules--tag_selectors) structure is documented below.

<a name="attrblock--triggers--scope_rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - Indicates the repository selector key.

* `value` - Indicates the repository selector value.
  The [value](#attrblock--triggers--scope_rules--scope_selectors--value) structure is documented below.

<a name="attrblock--triggers--scope_rules--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - Indicates the selector matching type.

* `kind` - Indicates the matching type.

* `pattern` - Indicates the pattern.

<a name="attrblock--triggers--scope_rules--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - Indicates the selector matching type.

* `kind` - Indicates the matching type.

* `pattern` - Indicates the pattern.

<a name="attrblock--triggers--targets"></a>
The `targets` block supports:

* `address` - Indicates the trigger address.

* `address_type` - Indicates the trigger address type.

* `auth_header` - Indicates the auth header.

* `skip_cert_verify` - Indicates whether to skip the verification of the certificate.

* `type` - Indicates the trigger type.
