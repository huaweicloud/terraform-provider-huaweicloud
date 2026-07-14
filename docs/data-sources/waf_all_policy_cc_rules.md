---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_policy_cc_rules"
description: |-
  Use this data source to get the list of WAF all policy CC rules within HuaweiCloud.
---

# huaweicloud_waf_all_policy_cc_rules

Use this data source to get the list of WAF all policy CC rules within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_waf_all_policy_cc_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `policyids` - (Optional, String) Specifies the policy IDs.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**. If you want to query resources under all enterprise projects, set this parameter to
  **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `items` - The CC rule list.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `name` - The rule name.

* `id` - The rule ID.

* `policyid` - The policy ID.

* `url` - The URL to protect. When mode is set to `0`, this parameter has a return value.
  It specifies the URL link for the rule application, excluding the domain name.

* `prefix` - Whether the path is in prefix mode. When the protection URL ends with `*`, it is in prefix mode.
  When mode is set to `0`, this parameter has a return value.

* `mode` - The CC rule protection mode.  
  The valid values are as follows:
  + **0**: Standard protection, only supports limiting the protection path of the domain name.
  + **1**: Advanced protection, supports limiting the path, IP, Cookie, Header, and Params fields.

* `status` - The rule status.  
  The valid values are as follows:
  + **0**: Disabled.
  + **1**: Enabled.

* `conditions` - The CC rule protection rule rate limiting conditions.

  The [conditions](#items_conditions) structure is documented below.

* `action` - The protection action taken after the request limit is reached.

  The [action](#items_action) structure is documented below.

* `tag_type` - The rate limiting mode.  
  The valid values are as follows:
  + **ip**: IP rate limiting, distinguishes individual web visitors based on IP.
  + **cookie**: User rate limiting, distinguishes individual web visitors based on Cookie key value.
  + **header**: User rate limiting, distinguishes individual web visitors based on Header.
  + **other**: Distinguishes individual web visitors based on the Referer field (custom request source).
  + **policy**: Policy rate limiting.
  + **domain**: Domain name rate limiting.
  + **url**: URL rate limiting.

* `tag_index` - The user identifier. When the rate limiting mode is user rate limiting (cookie or header),
  this parameter is required.
  + When cookie is selected, set the cookie field name, which means the user needs to configure the variable name of a
    certain attribute in the cookie that can uniquely identify the web visitor according to the actual situation of the
    website.
  + When header is selected, set the custom HTTP header that needs protection, which means the user needs to configure
    the HTTP header that can identify the web visitor according to the actual situation of the website.

* `tag_condition` - The user identifier. When the rate limiting mode is other, this parameter is required.
  Distinguishes individual web visitors based on the Referer field (custom request source).

  The [tag_condition](#items_tag_condition) structure is documented below.

* `limit_num` - The rate limit. Range: `1` ~ `2,147,483,647`, unit: times.

* `limit_period` - The rate limiting period. Range: `1` ~ `3,600`, unit: seconds.

* `unlock_num` - The pass frequency. Range: `1` ~ `2,147,483,647`, unit: times.

* `lock_time` - The blocking time. When "Protection Action" is set to "Block", you can set the time to resume normal
  access to the page after blocking. Range: `0` ~ `65,535`, unit: seconds.

* `domain_aggregation` - Whether to enable domain name aggregation statistics.

* `region_aggregation` - Whether to enable global counting.

* `description` - The rule description.

* `timestamp` - The timestamp when the rule was created. `13`-bit millisecond timestamp.

<a name="items_conditions"></a>
The `conditions` block supports:

* `category` - The field type.  
  The valid values are as follows:
  + **url**
  + **ip**
  + **ipv6**
  + **params**
  + **cookie**
  + **header**

* `logic_operation` - The matching logic for the condition list.
  + If the field type `category` is URL, the matching logic can be: contain, not_contain, equal, not_equal, prefix,
    not_prefix, suffix, not_suffix, contain_any, not_contain_all, equal_any, not_equal_all, prefix_any, not_prefix_all,
    suffix_any, not_suffix_all, len_greater, len_less, len_equal, or len_not_equal.
  + If the field type `category` is IP or IPv6, the matching logic can be: equal, not_equal, equal_any, or not_equal_all.
  + If the field type `category` is params, cookie, or header, the matching logic can be: contain, not_contain, equal,
    not_equal, prefix, not_prefix, suffix, not_suffix, contain_any, not_contain_all, equal_any, not_equal_all,
    prefix_any, not_prefix_all, suffix_any, not_suffix_all, len_greater, len_less, len_equal, len_not_equal,
    num_greater, num_less, num_equal, num_not_equal, exist, or not_exist.

* `contents` - The logical matching content of the condition list. This parameter is required when the
  `logic_operation` parameter does not end with any or all.

* `value_list_id` - The reference table ID. This parameter is required when the `logic_operation` parameter
  ends with any or all. In addition, the reference table type must be consistent with the category type.

* `index` - The subfield.

<a name="items_action"></a>
The `action` block supports:

* `category` - The action type.  
  The valid values are as follows:
  + **captcha**: CAPTCHA verification. After blocking, users need to enter the correct verification code to resume
    access to the correct page.
  + **block**: Block.
  + **log**: Log only.
  + **dynamic_block**: In the previous rate limiting cycle, if the request frequency exceeds the rate limit,
    it will be blocked. Then in the next rate limiting cycle, if the request frequency exceeds the pass frequency,
    it will be blocked. This protection action is only supported when the CC protection rule mode is advanced mode.

* `detail` - The blocking page information. When the protection action (category) is set to block or dynamic_block,
  the blocking page to be returned needs to be set.

  The [detail](#items_action_detail) structure is documented below.

<a name="items_action_detail"></a>
The `detail` block supports:

* `response` - The response details.

  The [response](#items_action_detail_response) structure is documented below.

<a name="items_action_detail_response"></a>
The `response` block supports:

* `content_type` - The content type.  
  The valid values are as follows:
  + **application/json**
  + **text/html**
  + **text/xml**

* `content` - The block page content.

<a name="items_tag_condition"></a>
The `tag_condition` block supports:

* `category` - The user identification field. The valid value is **referer**.

* `contents` - The user identification field contents list.
