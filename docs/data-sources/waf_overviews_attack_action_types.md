---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_attack_action_types"
description: |-
  Use this data source to query the events by protective actions.
---

# huaweicloud_waf_overviews_attack_action_types

Use this data source to query the events by protective actions.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_attack_action_types" "test" {
  from = var.start_time
  to   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `from` - (Required, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `to` - (Required, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The attcack action details.
  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `key` - The protective action.
  The valid values are as follows:
  + **block**: WAF blocks requests that trigger the rate limit set in the rule.
  + **log**: WAF logs requests that trigger the rate limit set in the rule but does not block them.
  + **captcha**: If the request exceeds the rate limit you configure, a verification code is displayed for
  human-machine verification. The request will not be allowed unless the verification is successful. Currently,
  verification code supports English.
  + **dynamic_block**: Requests that trigger the rule are blocked based on the allowable frequency you configure
  after the first rate limit period is over.
  + **advanced_captcha**: If your website visitor triggers the rate limit you set, CAPTCHA verification is required.

* `num` - The quantity.
