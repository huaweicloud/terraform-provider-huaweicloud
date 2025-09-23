---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_precise_protection"
description: |-
  Manages a WAF precise protection rule resource within HuaweiCloud.
---

# huaweicloud_waf_rule_precise_protection

Manages a WAF precise protection rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The precise protection rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable "policy_id" {}
variable "enterprise_project_id" {}
variable "reference_table_id" {}

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  name                  = "rule_01"
  priority              = 10
  action                = "block"
  start_time            = "2023-05-01 12:00:00"
  end_time              = "2023-05-10 12:00:00"
  description           = "description information"
  status                = 1

  conditions {
    field   = "url"
    logic   = "contain"
    content = "login"
  }

  conditions {
    field    = "params"
    logic    = "contain"
    subfield = "param_info"
    content  = "register"
  }

  conditions {
    field              = "header"
    logic              = "prefix_any"
    subfield           = "test_sub"
    reference_table_id = var.reference_table_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID of WAF precise protection rule.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of WAF precise protection rule.

* `priority` - (Required, Int) Specifies the priority of a rule. Smaller values correspond to higher priorities.
  If two rules are assigned with the same priority, the rule added earlier has higher priority.
  The value ranges from 0 to 1000.

* `conditions` - (Required, List) Specifies the match condition list.
  The [conditions](#RulePreciseProtection_conditions) structure is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF precise protection
  rule. For enterprise users, if omitted, default enterprise project will be used.

  Changing this parameter will create a new resource.

* `action` - (Optional, String) Specifies the protective action of WAF precise protection rule.
  Valid values are **block**, **pass**, **log**. The default value is **block**.

* `known_attack_source_id` - (Optional, String) Specifies the known attack source ID.
  The requirements for using this parameter are as follows:
  + The field is valid only when the `action` is set to **block**.
  + The policy needs to be bound to a domain name.
  + Before enabling `Cookie` or `Params` known attack source rules, configure a session or user tag for the
  corresponding website domain name.
  Refer to [Configure Traffic Identifier](https://support.huaweicloud.com/intl/en-us/usermanual-waf/waf_01_0270.html)

* `status` - (Optional, Int) Specifies the status of WAF precise protection rule.
  Valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

  The default value is `1`.

* `start_time` - (Optional, String) Specifies the time when the precise protection rule takes effect.
  The time format is **yyyy-MM-dd HH:mm:ss**, e.g. **2023-05-01 15:04:05**.

* `end_time` - (Optional, String) Specifies the time when the precise protection rule expires.
  The time format is **yyyy-MM-dd HH:mm:ss**, e.g. **2023-05-02 15:04:05**.

-> The precise protection rule will take effect immediately when `start_time` and `end_time` are not specified.

* `description` - (Optional, String) Specifies the description of WAF precise protection rule.

<a name="RulePreciseProtection_conditions"></a>
The `conditions` block supports:

* `field` - (Required, String) Specifies the field of the condition. Valid values are **url**, **user-agent**,
  **referer**, **ip**, **method**, **request_line**, **request**, **params**, **cookie**, **header**.

* `logic` - (Required, String) Specifies the condition matching logic.

  + If `field` is set to **url** or **user-agent** or **referer**: Valid values are **contain**, **not_contain**,
    **equal**, **not_equal**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **contain_any**,
    **not_contain_all**, **equal_any**, **not_equal_all**, **equal_any**, **not_equal_all**, **prefix_any**,
    **not_prefix_all**, **suffix_any**, **not_suffix_all**, **len_greater**, **len_less**, **len_equal**,
    **len_not_equal**.

  + If `field` is set to **ip**: Valid values are **equal**, **not_equal**, **equal_any**, **not_equal_all**.

  + If `field` is set to **method**: Valid values are **equal**, **not_equal**.

  + If `field` is set to **request_line** or **request**: Valid values are **len_greater**, **len_less**, **len_equal**,
    **len_not_equal**.

  + If `field` is set to **params** or **cookie** or **header**: Valid values are **contain**, **not_contain**,
    **equal**, **not_equal**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **contain_any**,
    **not_contain_all**, **equal_any**, **not_equal_all**, **equal_any**, **not_equal_all**, **prefix_any**,
    **not_prefix_all**, **suffix_any**, **not_suffix_all**, **len_greater**, **len_less**, **len_equal**,
    **len_not_equal**, **num_greater**, **num_less**, **num_equal**, **num_not_equal**, **exist**, **not_exist**.

* `subfield` - (Optional, String) Specifies the subfield of the condition. The parameter is required when `field`
  is set to **params**, **header** and **cookie**.

  + If `field` is set to **cookie**: The parameter indicates cookie name.

  + If `field` is set to **params**: The parameter indicates param name.

  + If `field` is set to **header**: The parameter indicates an option in the header.

* `content` - (Optional, String) Specifies the content of the match condition. It is required when the `logic`
  does not end with **any** or **all**.

* `reference_table_id` - (Optional, String) Specifies the reference table id. It is required when the `logic`
  end with **any** or **all**. The type of reference table should be consistent with the type of `field`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of a precise protection rule.

## Import

There are two ways to import WAF rule precise protection state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_precise_protection.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_precise_protection.test <policy_id>/<rule_id>/<enterprise_project_id>
```
