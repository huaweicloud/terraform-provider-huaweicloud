---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_policy_v2"
description: |-
  Manages a WAF policy V2 resource within HuaweiCloud.
---

# huaweicloud_waf_policy_v2

Manages a WAF policy V2 resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The policy resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable "name" {}
variable "enterprise_project_id" {}

resource "huaweicloud_waf_policy_v2" "test" {
  name                  = var.name
  enterprise_project_id = var.enterprise_project_id
  log_action_replaced   = "true"
  full_detection        = "false"
  level                 = 2

  action {
    category = "log"
  }

  robot_action {
    category = "log"
  }

  options {
    antileakage     = "false"
    antitamper      = "true"
    bot_enable      = "true"
    cc              = "true"
    crawler_other   = "false"
    crawler_scanner = "true"
    crawler_script  = "false"
    custom          = "true"
    geoip           = "true"
    ignore          = "true"
    modulex_enabled = "false"
    privacy         = "true"
    webattack       = "true"
    webshell        = "false"
    whiteblackip    = "true"
  }

  extend = {
    "extend" = jsonencode(
      {
        deep_decode             = true
        log_action_replaced     = false
        shiro_rememberMe_enable = true
      }
    )
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF policy V2 resource.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `name` - (Required, String) Specifies the policy name.

* `log_action_replaced` - (Optional, String, NonUpdatable) Specifies the log action replacement flag.
  This indicates whether Web Basic Protection will hit the policy rules and block when the "Protection Action" of the
  CC rule and the Precision Protection rule is set to "Log Only". The default value is **true**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

* `level` - (Optional, Int) Specifies the web basic protection level. Valid values are:
  + `1`: Low. At this protection level, WAF blocks only requests with obvious attack features. If a large number of
    false alarms have been reported, this value is recommended.
  + `2`: Medium. This protection level meets web protection requirements in most scenarios.
  + `3`: High. At this protection level, WAF provides the finest granular protection and can intercept attacks with
    complex bypass features, such as Jolokia cyberattacks, common gateway interface (CGI) vulnerability detection,
    and Druid SQL injection attacks.

* `full_detection` - (Optional, String) Specifies the detection mode in precise protection. Valid values are:
  + **false**: Instant detection. When a request hits the blocking conditions in precise protection, WAF terminates
    checks and blocks the request immediately.
  + **true**: Full detection. If a request hits the blocking conditions in precise protection, WAF does not block the
    request immediately. Instead, it blocks the requests until other checks are finished.

* `robot_action` - (Optional, List) Specifies the anti-crawler feature action information. Maximum items: `1`.

  The [robot_action](#Robot_Action) structure is documented below.

* `action` - (Optional, List) Specifies the protection action information. Maximum items: `1`.

  The [action](#Action) structure is documented below.

* `options` - (Optional, List) Specifies the switch information of protection items in the protection policy.
  Maximum items: `1`.

  The [options](#Options) structure is documented below.

* `hosts` - (Optional, List) Specifies the array of domain IDs bound to the protection policy.

* `bind_host` - (Optional, List) Specifies the array of domain information bound to the protection policy,
  which contains more detailed domain information than the hosts field.

  The [bind_host](#Bind_Host) structure is documented below.

* `extend` - (Optional, Map) Specifies the extension field for storing some switch configuration information
  in web basic protection. When modifying fields such as `shiro_rememberMe_enable`, `deep_decode`, `check_all_headers`,
  an additional layer of extend field nesting is required.
  Example: key is **extend**, value is **{"shiro_rememberMe_enable":true}**.

<a name="Bind_Host"></a>
The `bind_host` block supports:

* `id` - (Optional, String) Specifies the domain name ID.

* `hostname` - (Optional, String) Specifies the domain name.

* `waf_type` - (Optional, String) Specifies the deployment mode corresponding to the domain name. Valid values are:
  + **cloud**: Cloud mode.
  + **premium**: Dedicated mode.

* `mode` - (Optional, String) Specifies the special domain name mode. This attribute is only valid for dedicated mode.

<a name="Robot_Action"></a>
The `robot_action` block supports:

* `category` - (Optional, String) Specifies the protective action for feature anti-crawler. Valid values are:
  + **log**: Record only.
  + **block**: Block.

<a name="Action"></a>
The `action` block supports:

* `category` - (Optional, String) Specifies the web basic protection action. Valid values are:
  + **log**: Record only.
  + **block**: Block.

* `followed_action_id` - (Optional, String) Specifies the attack punishment rule ID.

<a name="Options"></a>
The `options` block supports:

* `webattack` - (Optional, String) Specifies whether basic protection is enabled.

* `common` - (Optional, String) Specifies whether conventional detection is enabled.

* `crawler_engine` - (Optional, String) Specifies whether the search engine is enabled.

* `crawler_scanner` - (Optional, String) Specifies whether anti-crawler detection is enabled.

* `crawler_script` - (Optional, String) Specifies whether script anti-crawler is enabled.

* `crawler_other` - (Optional, String) Specifies whether other crawlers are enabled.

* `webshell` - (Optional, String) Specifies whether Webshell detection is enabled.

* `cc` - (Optional, String) Specifies whether CC rules are enabled.

* `custom` - (Optional, String) Specifies whether precise protection is enabled.

* `whiteblackip` - (Optional, String) Specifies whether blacklist and whitelist protection is enabled.

* `geoip` - (Optional, String) Specifies whether geolocation access control rules are enabled.

* `ignore` - (Optional, String) Specifies whether global whitelist is enabled.

* `privacy` - (Optional, String) Specifies whether privacy masking is enabled.

* `antitamper` - (Optional, String) Specifies whether web tamper protection rules are enabled.

* `antileakage` - (Optional, String) Specifies whether sensitive information leakage prevention rules are enabled.

* `bot_enable` - (Optional, String) Specifies whether the website anti-crawler master switch is enabled.

* `modulex_enabled` - (Optional, String) Specifies whether modulex intelligent CC protection is enabled.
  This is a public beta feature and only supports record-only mode during the beta period.

  -> Fields `webattack`, `common`, `crawler_engine`, `crawler_scanner`, `crawler_script`, `crawler_other`, `webshell`,
  `cc`, `custom`, `whiteblackip`, `geoip`, `ignore`, `privacy`, `antitamper`, `antileakage`, `bot_enable`,
  and `modulex_enabled` are boolean values ​​of type string. Valid values are **true** and **false**.
  If not configured, these fields will be set to default values ​​by the service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The policy ID in UUID format.

* `extend_attribute` - The extended fields are used to store information such as switch configurations in basic web
  protection.

* `timestamp` - The time to create protection policy.

## Import

There are two ways to import WAF policy V2 resources.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_policy_v2.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_policy_v2.test <id>/<enterprise_project_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `log_action_replaced`, `extend`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_waf_policy_v2" "test" {
  ...

  lifecycle {
    ignore_changes = [
      log_action_replaced,
      extend,
    ]
  }
}
```
