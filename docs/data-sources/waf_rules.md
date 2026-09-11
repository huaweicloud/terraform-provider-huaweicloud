---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules"
description: |-
  Use this data source to query protection rules of a specified type in all policies.
---

# huaweicloud_waf_rules

Use this data source to query protection rules of a specified type in all policies.

## Example Usage

```hcl
variable "rule_type" {}

data "huaweicloud_waf_rules" "test" {
  rule_type = var.rule_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the rules.
  If omitted, the provider-level region will be used.

* `rule_type` - (Required, String) Specifies the type of the rules to be queried.
  The valid values are as follows:
  + **cc**: Indicates CC attack protection.
  + **custom**: Indicates precise protection rules.
  + **whiteblackip**: Indicates blacklist and whitelist.
  + **geoip**: Indicates geolocation-based protection.
  + **ip-reputation**: Indicates threat intelligence.
  + **antitamper**: Indicates anti-tamper events.
  + **antileakage**: Indicates sensitive information leakage prevention.
  + **ignore**: Indicates global protection whitelist.
  + **privacy**: Indicates privacy masking.

* `policy_ids` - (Optional, String) Specifies the ID list of the policies to be queried.
  Multiple policy IDs are separated by commas (,).

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the rules belong.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - All rules that matched the filter parameters.
