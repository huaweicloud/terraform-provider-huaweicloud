---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_policies"
description: |-
  Use this data source to get a list of WAF policies.
---

# huaweicloud_waf_policies

Use this data source to get a list of WAF policies.

## Example Usage

```hcl
variable "policy_name" {}
variable "enterprise_project_id" {}

data "huaweicloud_waf_policies" "test" {
  name                  = var.policy_name
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the WAF policies. If omitted,
  the provider-level region will be used.

* `name` - (Optional, String) Specifies the policy name used for matching. The value is case-sensitive and supports
  fuzzy matching.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of WAF policies.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - A list of WAF policies.

The `policies` block supports:

* `id` - The WAF Policy ID.

* `name` - The WAF policy name.

* `full_detection` - The detection mode in precise protection.
  + **false**: Instant detection. When a request hits the blocking conditions in precise protection, WAF terminates
    checks and blocks the request immediately.
  + **true**: Full detection. If a request hits the blocking conditions in precise protection, WAF does not block the
    request immediately. Instead, it blocks the requests until other checks are finished.

* `protection_mode` - The protective action after a rule is matched. Valid values are:
  + **block**: WAF blocks and logs detected attacks.
  + **log**: WAF logs detected attacks only.

* `robot_action` - The protective actions for each rule in anti-crawler protection. Valid values are:
  + **block**: WAF blocks discovered attacks.
  + **log**: WAF only logs discovered attacks.

* `level` - The protection level. Valid values are:
  + **1**: Low. At this protection level, WAF blocks only requests with obvious attack features. If a large number of
    false alarms have been reported, this value is recommended.
  + **2**: Medium. This protection level meets web protection requirements in most scenarios.
  + **3**: High. At this protection level, WAF provides the finest granular protection and can intercept attacks with
    complex bypass features, such as Jolokia cyberattacks, common gateway interface (CGI) vulnerability detection,
    and Druid SQL injection attacks.

* `options` - The protection switches. The options object structure is documented below.

* `bind_hosts` - The protection switches. The object structure is documented below.

* `deep_inspection` - The deep inspection in basic web protection.

* `header_inspection` - The header inspection in basic web protection.

* `shiro_decryption_check` - The shiro decryption check in basic web protection.

The `options` block supports:

* `basic_web_protection` - Indicates whether Basic Web Protection is enabled.

* `general_check` - Indicates whether General Check in Basic Web Protection is enabled.

* `webshell` - Indicates whether the web shell detection in basic web protection is enabled.

* `crawler_engine` - Indicates whether the search engine is enabled.

* `crawler_scanner` - Indicates whether the anti-crawler detection is enabled.

* `crawler_script` - Indicates whether the script tool is enabled.

* `crawler_other` - Indicates whether other crawler check is enabled.

* `cc_attack_protection` - Indicates whether the cc attack protection rules are enabled.

* `precise_protection` - Indicates whether the precise protection is enabled.

* `blacklist` - Indicates whether the blacklist and whitelist protection is enabled.

* `data_masking` - Indicates whether data masking is enabled.

* `false_alarm_masking` - Indicates whether false alarm masking is enabled.

* `web_tamper_protection` - Indicates whether the web tamper protection is enabled.

* `geolocation_access_control` - Indicates whether the geolocation access control is enabled.

* `information_leakage_prevention` - Indicates whether the information leakage prevention is enabled.

* `bot_enable` - Indicates whether the anti-crawler protection is enabled.

* `known_attack_source` - Indicates whether the known attack source is enabled.

* `anti_crawler` - Indicates whether the javascript anti-crawler is enabled.

The `bind_hosts` block supports:

* `id` - The domain name ID.

* `hostname` - The domain name.

* `waf_type` - The deployment mode of WAF instance that is used for the domain name. The value can be **cloud** for
  cloud WAF or **premium** for dedicated WAF instances.

* `mode` - The special domain name mode. This attribute is only valid for dedicated mode.
