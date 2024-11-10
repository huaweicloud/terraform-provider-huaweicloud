---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_global_protection_whitelist"
description: |-
  Manages a WAF global protection whitelist rule resource within HuaweiCloud.
---

# huaweicloud_waf_rule_global_protection_whitelist

Manages a WAF global protection whitelist rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The global protection whitelist rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

### WAF Global Protection Whitelist Rule with specified domain names

```hcl
variable "policy_id" {}
variable "enterprise_project_id" {}
variable "domains" {
  type = list(string)
}

resource "huaweicloud_waf_rule_global_protection_whitelist" "test" {
  policy_id             = var.policy_id
  domains               = var.domains
  enterprise_project_id = var.enterprise_project_id
  ignore_waf_protection = "xss;webshell"
  advanced_field        = "params"
  advanced_content      = "test_content"
  description           = "test description"

  conditions {
    field    = "ip"
    logic    = "equal"
    content  = "192.168.0.2"
    subfield = "x-forwarded-for"
  }
}
```

### WAF Global Protection Whitelist Rule with all domain names

```hcl
variable "policy_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_waf_rule_global_protection_whitelist" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  domains               = []
  ignore_waf_protection = "xss;webshell"
  advanced_field        = "params"
  advanced_content      = "test_content"
  description           = "test description"

  conditions {
    field    = "params"
    logic    = "contain"
    content  = "test content"
    subfield = "test_subfield"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID of WAF global protection whitelist rule.

  Changing this parameter will create a new resource.

* `domains` - (Required, List) Specifies the protected domain name bound with the policy or manually enter a single
  domain name corresponding to the wildcard domain name.
  If the array length is `0`, the rule takes effect for all domain names or websites.

* `ignore_waf_protection` - (Required, String) Specifies the rules that need to be ignored. You can provide multiple
  items and separate them with semicolons (;).

  + If you want to block a specific built-in rule, the value of this parameter is the rule ID.
  To query the rule ID, go to the WAF console, choose **Policies** and click the target policy name. On the displayed
  page, in the **Basic Web Protection** area, select the **Protection Rules** tab, and view the ID of the specific rule.
  You can also query the rule ID in the event details.

  + If you want to mask a type of basic web protection rules, set this parameter to the name of the type of basic web
  protection rules. Valid values are: **xss**(XSS attacks), **webshell**(Web shells), **vuln**(Other types of attacks),
  **sqli**(SQL injection attack), **robot**(Malicious crawlers), **rfi**(Remote file inclusion),
  **lfi**(Local file inclusion), **cmdi**(Command injection attack).

  + To bypass the basic web protection, set this parameter to **all**.

  + To bypass all WAF protection, set this parameter to **bypass**.

* `conditions` - (Required, List) Specifies the match condition list.
  The [conditions](#RuleGlobalProtectionWhitelist_conditions) structure is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF global protection
  whitelist rule.

  Changing this parameter will create a new resource.

* `advanced_field` - (Optional, String) Specifies the advanced field to ignore attacks of a specific field.
  After you add the rule, WAF will stop intercepting attack events of the specified field.
  The following fields are supported: **params**, **cookie**, **header**, **body** and **multipart**.

* `advanced_content` - (Optional, String) Specifies the advanced content value to ignore. This parameter is valid only
  when `advanced_field` is set to **params**, **cookie** or **header**.
  If not specified, WAF will ignore all attack events of the specific field.

* `description` - (Optional, String) Specifies the description of WAF global protection whitelist rule.

* `status` - (Optional, Int) Specifies the status of WAF global protection whitelist rule.
  Valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

  The default value is `1`.

<a name="RuleGlobalProtectionWhitelist_conditions"></a>
The `conditions` block supports:

* `field` - (Required, String) Specifies the field type. The value can be **ip**, **url**, **params**, **cookie**
  or **header**.

* `logic` - (Required, String) Specifies the condition matching logic.

  + If `field` is set to **ip**: Valid values are **equal** and **not_equal**.

  + If `field` is set to **url** or **header** or **params** or **cookie**: Valid values are **equal**, **not_equal**,
  **contain**, **not_contain**, **prefix**, **not_prefix**, **suffix** and **not_suffix**.

* `content` - (Required, String) Specifies the content of the match condition.

  + If `field` is set to **ip**, the value must be an IP address or IP address range.

  + If `field` is set to **url**, the value must be in the standard URL format.

  + If `field` is set to **params** or **cookie** or **header**, the content format is not limited.

* `subfield` - (Optional, String) Specifies the subfield of the condition.

  + If `field` is set to **ip** and the subfield is the client IP address, the parameter is not required.

  + If `field` is set to **ip** and the subfield is X-Forwarded-For, the parameter is required and the value should be
  **x-forwarded-for**.

  + If `field` is set to **params**, **header** or **cookie**, the parameter is required and the value is user-defined.

  + If `field` is set to **url**, the parameter cannot be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of a global protection whitelist rule.

## Import

There are two ways to import WAF rule global protection whitelist state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_global_protection_whitelist.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_global_protection_whitelist.test <policy_id>/<rule_id>/<enterprise_project_id>
```
