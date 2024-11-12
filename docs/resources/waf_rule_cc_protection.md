---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_cc_protection"
description: |-
  Manages a WAF cc protection rule resource within HuaweiCloud.
---

# huaweicloud_waf_rule_cc_protection

Manages a WAF cc protection rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The cc protection rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable "policy_id" {}
variable "enterprise_project_id" {}
variable "reference_table_id" {}

resource "huaweicloud_waf_rule_cc_protection" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  name                  = "test_rule"
  protective_action     = "block"
  rate_limit_mode       = "cookie"
  block_page_type       = "application/json"
  page_content          = "test page content"
  user_identifier       = "test_identifier"
  limit_num             = 10
  limit_period          = 60
  lock_time             = 5
  request_aggregation   = true
  all_waf_instances     = true
  description           = "test description"

  conditions {
    field    = "params"
    logic    = "contain"
    content  = "test content"
    subfield = "test_subfield"
  }

  conditions {
    field              = "header"
    logic              = "prefix_any"
    subfield           = "test_subfield"
    reference_table_id = var.reference_table_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID of WAF cc protection rule.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the rule name of WAF cc protection rule.

* `conditions` - (Required, List) Specifies the match condition list.
The [conditions](#RuleCCProtection_conditions) structure is documented below.

* `protective_action` - (Required, String) Specifies the protective action taken when the number of requests reaches
  the upper limit. Valid values are as follows:
  + **captcha**: Verification code. The user needs to enter the correct verification code after blocking to restore the
  correct access page.
  + **block**: Block the requests.
  + **log**: Record only.
  + **dynamic_block**: Dynamic block the requests. If the request frequency exceeds the "speed limit frequency" during
  the previous speed limit cycle. In the next speed limit cycle, if the request frequency exceeds the
  "release frequency", it will be blocked.

* `rate_limit_mode` - (Required, String) Specifies the rate limit mode. Valid values are as follows:
  + **ip**: A web visitor is identified by the IP address.
  + **cookie**: A web visitor is identified by the cookie key value.
  + **header**: A web visitor is identified by the header key value.
  + **other**: A web visitor is identified by the Referer field (user-defined request source).
  + **policy**: A web visitor is identified by rule.
  + **domain**: A web visitor is identified by domain name.
  + **url**: A web visitor is identified by url.

* `limit_num` - (Required, Int) Specifies the number of requests allowed from a web visitor in a rate limiting period.
  The value ranges from `1` to `2,147,483,647`.

* `limit_period` - (Required, Int) Specifies the rate limiting period. The value ranges from `1` to `3,600` in seconds.

* `block_page_type` - (Optional, String) Specifies the type of the returned page. The options are **application/json**,
  **text/html** and **text/xml**. This parameter is valid when `protective_action` is set to **block** or **dynamic_block**.
  If not specified the system default block page will be used.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF cc protection rule.
  For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

* `page_content` - (Optional, String) Specifies the content of the returned page.
  This parameter is required when `block_page_type` has value.

* `user_identifier` - (Optional, String) Specifies the user identifier.
  This parameter is required when `rate_limit_mode` is set to **cookie** or **header**.
  + If `rate_limit_mode` is set to **cookie**, this parameter indicates cookie name.
  + If `rate_limit_mode` is set to **header**, this parameter indicates header name.

* `other_user_identifier` - (Optional, String) Specifies the other user identifier.
  This parameter is required when `rate_limit_mode` is set to **other**, indicates the user-defined request field.

* `unlock_num` - (Optional, Int) Specifies the allowable frequency. The value ranges from `0` to `2,147,483,647`.
  This parameter is valid when `protective_action` is set to **dynamic_block**.

* `lock_time` - (Optional, Int) Specifies the lock time for resuming normal page access after blocking can be set.
  The value ranges from `0` to `65,535` in seconds. This parameter is valid when `protective_action` is set to **block**.

* `request_aggregation` - (Optional, Bool) Specifies whether to enable domain aggregation statistics.
  This parameter is valid when `rate_limit_mode` is not set to **policy**. Default to **false**.

* `all_waf_instances` - (Optional, Bool) Specifies whether to enable global counting. Default to **false**.

* `description` - (Optional, String) Specifies the description of WAF cc protection rule.

* `status` - (Optional, Int) Specifies the status of WAF cc protection rule.
  Valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

  The default value is `1`.

<a name="RuleCCProtection_conditions"></a>
The `conditions` block supports:

* `field` - (Required, String) Specifies the field type.
  The value can be **url**, **ip**, **ipv6**, **params**, **cookie**, **header** or **response_code**.

* `logic` - (Required, String) Specifies the condition matching logic.

  + If `field` is set to **url**: Valid values are **contain**, **not_contain**, **equal**, **not_equal**, **prefix**,
  **not_prefix**, **suffix**, **not_suffix**, **contain_any**, **not_contain_all**, **equal_any**, **not_equal_all**,
  **equal_any**, **not_equal_all**, **prefix_any**, **not_prefix_all**, **suffix_any**, **not_suffix_all**,
  **len_greater**, **len_less**, **len_equal** and **len_not_equal**.

  + If `field` is set to **ip** or **ipv6**: Valid values are **equal**, **not_equal**, **equal_any** and
  **not_equal_all**.

  + If `field` is set to **response_code**: Valid values are **equal** and **not_equal**.

  + If `field` is set to **params**, **cookie** or **header**: Valid values are **contain**, **not_contain**,
  **equal**, **not_equal**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **contain_any**,
  **not_contain_all**, **equal_any**, **not_equal_all**, **equal_any**, **not_equal_all**, **prefix_any**,
  **not_prefix_all**, **suffix_any**, **not_suffix_all**, **len_greater**, **len_less**, **len_equal**,
  **len_not_equal**, **num_greater**, **num_less**, **num_equal**, **num_not_equal**, **exist** and **not_exist**.

* `subfield` - (Optional, String) Specifies the subfield of the condition.
  It is required when `field` is set to **params**, **header** or **cookie**.

* `content` - (Optional, String) Specifies the content of the match condition.
  It is required when the `logic` does not end with **any** or **all**.

* `reference_table_id` - (Optional, String) Specifies the reference table ID.
  It is required when the `logic` end with **any** or **all**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of a WAF cc protection rule.

## Import

There are two ways to import WAF rule cc protection state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_cc_protection.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_cc_protection.test <policy_id>/<rule_id>/<enterprise_project_id>
```
