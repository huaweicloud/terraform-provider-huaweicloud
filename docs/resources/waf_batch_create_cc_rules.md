---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_cc_rules"
description: |-
  Manages a resource to batch create WAF CC protection rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_cc_rules

Manages a resource to batch create WAF CC protection rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating CC protection rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_cc_rules" "test" {
  name                  = "test_cc_rule"
  description           = "CC protection rule for admin paths"
  tag_type              = "ip"
  limit_num             = 100
  limit_period          = 60
  policy_ids            = var.policy_ids
  enterprise_project_id = var.enterprise_project_id
  
  conditions {
    category        = "url"
    logic_operation = "prefix"
    contents        = ["/admin"]
  }

  action {
    category = "block"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CC protection rule.

* `conditions` - (Required, List, NonUpdatable) Specifies the list of matching conditions for rate limiting.
  The [conditions](#Rule_conditions) structure is documented below.

* `action` - (Required, List, NonUpdatable) Specifies the action to take when the rate limit is reached.
  The [action](#Rule_action) structure is documented below.

* `tag_type` - (Required, String, NonUpdatable) Specifies how to identify unique web visitors for rate limiting.
  The value can be:
  + **ip**: Rate limit by IP address
  + **cookie**: Rate limit by cookie value (requires `tag_index`)
  + **header**: Rate limit by header value (requires `tag_index`)
  + **other**: Rate limit by Referer (requires `tag_condition`)
  + **policy**: Rate limit by policy
  + **domain**: Rate limit by domain
  + **url**: Rate limit by URL

* `limit_num` - (Required, Int, NonUpdatable) Specifies the maximum number of requests allowed in the rate limit period.
  The value ranges from `1` to `2,147,483,647`.

* `limit_period` - (Required, Int, NonUpdatable) Specifies the rate limit period in seconds.
  The value ranges from `1` to `3,600` (1 hour).

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the CC rule will be applied.

* `tag_index` - (Optional, String, NonUpdatable) Specifies the cookie or header name when `tag_type` is set to
  **cookie** or **header**.
  + When selecting a cookie, set the cookie field name. This is a specific attribute variable name in the cookie that
    the user needs to configure to uniquely identify the web visitor, based on the website's actual situation.
    User-identifying cookies do not support regular expressions; they must be exact matches. For example, if the website
    uses a field called "name" in the cookie to uniquely identify a user, then the "name" field can be selected to
    distinguish web visitors.
  + When selecting a header, set the custom HTTP headers that need to be protected. This means that users need to
    configure HTTP headers that can identify web visitors according to the actual situation of the website.

* `tag_condition` - (Optional, List, NonUpdatable) Specifies the condition when `tag_type` is set to **other**.
  The [tag_condition](#Tag_condition) structure is documented below.

* `unlock_num` - (Optional, Int, NonUpdatable) Specifies the number of requests to allow after rate limiting is triggered.
  The value ranges from `0` to `2,147,483,647`.
  This parameter is valid only when the action `category` is set to **dynamic_block**.

* `lock_time` - (Optional, Int, NonUpdatable) Specifies the duration (in seconds) to block requests when the rate limit
  is exceeded. The value ranges from `0` to `65,535` seconds.
  This parameter is valid only when the action `category` is set to **block**.

* `domain_aggregation` - (Optional, Bool, NonUpdatable) Specifies whether to enable domain aggregation statistics.
  Defaults to **false**.

* `region_aggregation` - (Optional, Bool, NonUpdatable) Specifies whether to enable global counting.
  Defaults to **false**.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the CC protection rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

<a name="Rule_conditions"></a>
The `conditions` block supports:

* `category` - (Required, String, NonUpdatable) Specifies the field type to match against.
  The value can be:
  + **url**
  + **ip**
  + **ipv6**
  + **asn**
  + **params**
  + **cookie**
  + **referer**
  + **user-agent**
  + **header**
  + **response_code**
  + **response_header**
  + **response_body**
  + **request_body**
  + **method**
  + **tls_fingerprint**
  + **tls_ja3_fingerprint**

* `logic_operation` - (Required, String, NonUpdatable) Specifies the matching logic operation.
  The available values depend on the `category`:
  + For **url**: **contain**, **not_contain**, **equal**, **not_equal**, **prefix**, **not_prefix**, **suffix**,
    **not_suffix**, **contain_any**, **not_contain_all**, **equal_any**, **not_equal_all**, **prefix_any**,
    **not_prefix_all**, **suffix_any**, **not_suffix_all**, **len_greater**, **len_less**, **len_equal**,
    **len_not_equal**
  + For **ip** or **ipv6**: **equal**, **not_equal**, **equal_any**, **not_equal_all**
  + For **params**, **cookie**, or **header**: **contain**, **not_contain**, **equal**, **not_equal**, **prefix**,
    **not_prefix**, **suffix**, **not_suffix**, **contain_any**, **not_contain_all**, **equal_any**, **not_equal_all**,
    **prefix_any**, **not_prefix_all**, **suffix_any**, **not_suffix_all**, **len_greater**, **len_less**, **len_equal**,
    **len_not_equal**, **num_greater**, **num_less**, **num_equal**, **num_not_equal**, **exist**, **not_exist**

* `contents` - (Optional, List, NonUpdatable) Specifies the content to match against.
  Required when `logic_operation` does not end with **any** or **all**.

* `value_list_id` - (Optional, String, NonUpdatable) Specifies the reference table ID.
  Required when `logic_operation` ends with **any** or **all**.

* `index` - (Optional, String, NonUpdatable) Specifies the subfield name when `category` is **params**, **cookie**,
  or **header**.
  For other categories, leave this empty.

<a name="Rule_action"></a>
The `action` block supports:

* `category` - (Required, String, NonUpdatable) Specifies the action to take when the rate limit is reached.
  The value can be:
  + **captcha**: After a human-machine verification is blocked, the user needs to enter the correct verification code
    to restore access to the correct page
  + **log**: Only records the request
  + **dynamic_block**: In the previous rate limit period, if the request frequency exceeds the "rate limit frequency",
    it will be blocked. In the next rate limit period, if the request frequency exceeds the "release frequency",
    it will be blocked. Note: Only when the CC protection rule mode is set to advanced mode can `dynamic_block`
    protection action be set
  + **block**: Blocks the request

* `detail` - (Optional, List, NonUpdatable) Specifies the detailed action configuration. The returned blocking page is
  only required when the protection action (`category`) is selected as either `block` or `dynamic_block`.
  + If you need the system's default blocking page to be returned, you do not need to pass this parameter.
  + If users want to protect against a custom blocking page, they can configure this setting.

  The [detail](#Action_detail) structure is documented below.

<a name="Action_detail"></a>
The `detail` block supports:

* `response` - (Optional, List, NonUpdatable) Specifies the custom response configuration.
  The [response](#Response) structure is documented below.

<a name="Response"></a>
The `response` block supports:

* `content_type` - (Optional, String, NonUpdatable) Specifies the content type of the response.
  The value can be **application/json**, **text/html**, or **text/xml**.

* `content` - (Optional, String, NonUpdatable) Specifies the content of the response.

<a name="Tag_condition"></a>
The `tag_condition` block supports:

* `category` - (Optional, String, NonUpdatable) Specifies the field type to match against. The value can be **referer**.

* `contents` - (Optional, List, NonUpdatable) Specifies the content to match against.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
