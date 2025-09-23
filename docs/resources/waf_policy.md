---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_policy"
description: |-
  Manages a WAF policy resource within HuaweiCloud.
---

# huaweicloud_waf_policy

Manages a WAF policy resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The policy resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_waf_policy" "test" {
  name                   = "test_policy"
  protection_mode        = "log"
  robot_action           = "block"
  level                  = 2
  deep_inspection        = true
  header_inspection      = true
  shiro_decryption_check = true
  enterprise_project_id  = var.enterprise_project_id

  options {
    crawler_scanner                = true
    crawler_script                 = true
    false_alarm_masking            = true
    general_check                  = true
    geolocation_access_control     = true
    information_leakage_prevention = true
    known_attack_source            = true
    precise_protection             = true
    web_tamper_protection          = true
    webshell                       = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF policy resource. If omitted, the
  provider-level region will be used. Changing this setting will push a new certificate.

* `name` - (Required, String) Specifies the policy name. The maximum length is `256` characters. Only digits, letters,
  underscores (_), and hyphens (-) are allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF policy.
  For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

* `full_detection` - (Optional, Bool) Specifies the detection mode in precise protection. Defaults to **false**.
  + **false**: Instant detection. When a request hits the blocking conditions in precise protection, WAF terminates
    checks and blocks the request immediately.
  + **true**: Full detection. If a request hits the blocking conditions in precise protection, WAF does not block the
    request immediately. Instead, it blocks the requests until other checks are finished.

* `protection_mode` - (Optional, String) Specifies the protective action after a rule is matched. Defaults to **log**.
  Valid values are:
  + **block**: WAF blocks and logs detected attacks.
  + **log**: WAF logs detected attacks only.

* `robot_action` - (Optional, String) Specifies the protective actions for each rule in anti-crawler protection.
  Defaults to **log**. Valid values are:
  + **block**: WAF blocks discovered attacks.
  + **log**: WAF only logs discovered attacks.

* `level` - (Optional, Int) Specifies the protection level. Defaults to `2`. Valid values are:
  + `1`: Low. At this protection level, WAF blocks only requests with obvious attack features. If a large number of
    false alarms have been reported, this value is recommended.
  + `2`: Medium. This protection level meets web protection requirements in most scenarios.
  + `3`: High. At this protection level, WAF provides the finest granular protection and can intercept attacks with
    complex bypass features, such as Jolokia cyberattacks, common gateway interface (CGI) vulnerability detection,
    and Druid SQL injection attacks.

* `deep_inspection` - (Optional, Bool) Specifies the deep inspection in basic web protection. Defaults to **false**.

* `header_inspection` - (Optional, Bool) Specifies the header inspection in basic web protection. Defaults to **false**.

* `shiro_decryption_check` - (Optional, Bool) Specifies the shiro decryption check in basic web protection.
  Defaults to **false**.

* `options` - (Optional, List) Specifies the switch options of the protection item in the policy.
  The [options](#Policy_Options) structure is documented below.

<a name="Policy_Options"></a>
The `options` block supports:

* `basic_web_protection` - (Optional, Bool) Specifies whether basic web protection is enabled. Defaults to **false**.

* `general_check` - (Optional, Bool) Specifies whether the general check in basic web protection is enabled.
  Defaults to **false**.

* `webshell` - (Optional, Bool) Specifies whether the web shell detection in basic web protection is enabled.
  Defaults to **false**.

* `crawler_engine` - (Optional, Bool) Specifies whether the search engine is enabled. Defaults to **false**.

* `crawler_scanner` - (Optional, Bool) Specifies whether the anti-crawler detection is enabled. Defaults to **false**.

* `crawler_script` - (Optional, Bool) Specifies whether the script tool is enabled. Defaults to **false**.

* `crawler_other` - (Optional, Bool) Specifies whether other crawler check is enabled. Defaults to **false**.

* `cc_attack_protection` - (Optional, Bool) Specifies whether the cc attack protection rules are enabled.
  Defaults to **false**.

* `precise_protection` - (Optional, Bool) Specifies whether the precise protection is enabled. Defaults to **false**.

* `blacklist` - (Optional, Bool) Specifies whether the blacklist and whitelist protection is enabled.
  Defaults to **false**.

* `data_masking` - (Optional, Bool) Specifies whether data masking is enabled. Defaults to **false**.

* `false_alarm_masking` - (Optional, Bool) Specifies whether false alarm masking is enabled. Defaults to **false**.

* `web_tamper_protection` - (Optional, Bool) Specifies whether the web tamper protection is enabled.
  Defaults to **false**.

* `geolocation_access_control` - (Optional, Bool) Specifies whether the geolocation access control is enabled.
  Defaults to **false**.

* `information_leakage_prevention` - (Optional, Bool) Specifies whether the information leakage prevention is enabled.
  Defaults to **false**.

* `bot_enable` - (Optional, Bool) Specifies whether the anti-crawler protection is enabled. Defaults to **false**.

* `known_attack_source` - (Optional, Bool) Specifies whether the known attack source is enabled. Defaults to **false**.

* `anti_crawler` - (Optional, Bool) Specifies whether the javascript anti-crawler is enabled. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The policy ID in UUID format.

* `bind_hosts` - The protection switches. The options object structure is documented below.

The `bind_hosts` block supports:

* `id` - The domain name ID.

* `hostname` - The domain name.

* `waf_type` - The deployment mode of WAF instance that is used for the domain name. The value can be **cloud** for
  cloud WAF or **premium** for dedicated WAF instances.

* `mode` - The special domain name mode. This attribute is only valid for dedicated mode.

## Import

There are two ways to import WAF policy state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_policy.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_policy.test <id>/<enterprise_project_id>
```
