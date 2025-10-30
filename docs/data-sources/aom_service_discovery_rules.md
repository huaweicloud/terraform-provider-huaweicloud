---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_service_discovery_rules"
description: |-
  Use this data source to get the list of AOM service discovery rules.
---

# huaweicloud_aom_service_discovery_rules

Use this data source to get the list of AOM service discovery rules.

## Example Usage

```hcl
data "huaweicloud_aom_service_discovery_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the service discovery rule ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - Indicates the service discovery rules list.
  The [rules](#attrblock--rules) structure is documented below.

<a name="attrblock--rules"></a>
The `rules` block supports:

* `id` - Indicates the rule ID.

* `name` - Indicates the rule name.

* `detect_log_enabled` - Indicates whether the log collection is enabled.

* `discovery_rule_enabled` - Indicates whether the rule is enabled.

* `discovery_rules` - Indicates the discovery rule.
  The [discovery_rules](#attrblock--rules--discovery_rules) structure is documented below.

* `is_default_rule` - Indicates whether the rule is the default one.

* `log_file_suffix` - Indicates the log file suffix.

* `log_path_rules` - Indicates the log path configuration rule.
  The [log_path_rules](#attrblock--rules--log_path_rules) structure is documented below.

* `name_rules` - Indicates the naming rules for discovered services and applications.
  The [name_rules](#attrblock--rules--name_rules) structure is documented below.

* `priority` - Indicates the rule priority.

* `service_type` - Indicates the service type.

* `created_at` - Indicates the create time of the rule.

* `description` - Indicates the description of the rule.

<a name="attrblock--rules--discovery_rules"></a>
The `discovery_rules` block supports:

* `check_content` - Indicates the matched value.

* `check_mode` - Indicates the match condition.

* `check_type` - Indicates the match type.

<a name="attrblock--rules--log_path_rules"></a>
The `log_path_rules` block supports:

* `args` - Indicates the command.

* `name_type` - Indicates the value type.

* `value` - Indicates the log path.

<a name="attrblock--rules--name_rules"></a>
The `name_rules` block supports:

* `application_name_rule` - Indicates the application name rule.
  The [basic_name_rule](#attrblock--rules--basic_name_rule) structure is documented below.

* `service_name_rule` - Indicates the service name rule.
  The [basic_name_rule](#attrblock--rules--basic_name_rule) structure is documented below.

<a name="attrblock--rules--basic_name_rule"></a>
The `basic_name_rule` block supports:

* `args` - Indicates the input value.

* `name_type` - Indicates the value type.

* `value` - Indicates the application name.
