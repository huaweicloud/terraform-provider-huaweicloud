---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_anti_crawler"
description: |-
  Manages a WAF rule anti crawler resource within HuaweiCloud.
---

# huaweicloud_waf_rule_anti_crawler

Manages a WAF rule anti crawler resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The anti crawler rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable policy_id {}
variable enterprise_project_id {}

resource "huaweicloud_waf_rule_anti_crawler" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  name                  = "test_name"
  protection_mode       = "anticrawler_specific_url"
  priority              = 100
  description           = "test description"

  conditions {
    field   = "user-agent"
    logic   = "contain"
    content = "TR"
  }

  conditions {
    field   = "url"
    logic   = "equal"
    content = "/test/path"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID.

  Changing this parameter will create a new resource.

* `protection_mode` - (Required, String, ForceNew) Specifies the protection mode of WAF anti crawler rule.
  Changing this parameter will create a new resource. Valid values are as follows:
    + **anticrawler_specific_url**: Used to protect a specific path specified by the rule.
    + **anticrawler_except_url**: Used to protect all paths except the one specified by the rule.

  -> All rules in the current mode will take effect, while rules in another mode will become invalid.

* `name` - (Required, String) Specifies the rule name. The value should be a maximum of `128` characters. Only letters,
  digits, hyphens (-), underscores (_), colons (:) and periods (.) are allowed.

* `priority` - (Required, Int) Specifies the priority. A smaller value indicates a higher priority. If the value is
  the same, the rule is created earlier and the priority is higher. Value ranges from `0` to `65,535`.

* `conditions` - (Required, List) Specifies the match condition list.
  The [conditions](#RuleAntiCrawler_conditions) structure is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF anti crawler rule.
  For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the rule description.

<a name="RuleAntiCrawler_conditions"></a>
The `conditions` block supports:

* `field` - (Required, String) Specifies the field type. The valid values are **url** and **user-agent**.

* `logic` - (Required, String) Specifies the logic for matching the condition. The valid values are **contain**,
  **not_contain**, **equal**, **not_equal**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **contain_any**,
  **not_contain_all**, **equal_any**, **not_equal_all**, **prefix_any**, **not_prefix_all**, **suffix_any** and
  **not_suffix_all**.

* `content` - (Optional, String) Specifies the content of the condition.
  It is valid and required when the `logic` does not end with `any` or `all`.

* `reference_table_id` - (Optional, String) Specifies the reference table ID.
  It is valid and required when the `logic` end with `any` or `all`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The rule status.

## Import

There are two ways to import WAF rule anti crawler state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_anti_crawler.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_anti_crawler.test <policy_id>/<rule_id>/<enterprise_project_id>
```
