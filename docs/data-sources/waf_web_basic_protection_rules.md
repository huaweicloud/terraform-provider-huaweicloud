---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_web_basic_protection_rules"
description: |-
  Use this data source to get the list of web basic protection rules.
---

# huaweicloud_waf_web_basic_protection_rules

Use this data source to get the list of web basic protection rules.

## Example Usage

```hcl
data "huaweicloud_waf_web_basic_protection_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

* `from` - (Optional, Int) Specifies the start time (`13`-digit timestamp in ms).
  This parameter must be used together with parameter `to`.

* `to` - (Optional, Int) Specifies the end time (`13`-digit timestamp, in ms).
  This parameter must be used together with parameter `from`.

* `level` - (Optional, Int) Specifies the protection level of the rule set.
  If you select a loose rule set, there will be fewer false positives, but also more false negatives. If you select a
  tight rule set, the opposite is true. Valid values are:
  + `1`: loose
  + `2`: medium
  + `3`: tight

* `rule_id` - (Optional, String) Specifies the rule ID.

* `description` - (Optional, String) Specifies the description.

* `cve_number` - (Optional, String) Specifies the CVE ID.

* `risk_level` - (Optional, Int) Specifies the risk severity. Valid values are:
  + `1`: high risk
  + `2`: medium risk
  + `3`: low risk

* `protection_type_names` - (Optional, String) Specifies the protection type. Valid values are:
  + **vuln**: Others
  + **xss**: Cross-site scripting (XSS) attacks
  + **cmdi**: command injection attacks
  + **lfi**: local file inclusion attacks
  + **rfi**: remote file inclusion attacks
  + **webshell**: website Trojans
  + **robot**: malicious crawlers
  + **sqli**: SQL injections

* `application_type_names` - (Optional, String) Specifies the application type.
  For details, go to the Basic Web Protection page on the WAF console.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The array of built-in basic web protection rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The rule ID. Each rule has a unique ID.

* `cve_number` - The CVE ID.

* `risk_level` - The risk severity.

* `effective_time` - The effective time.

* `create_time` - The creation time.

* `update_time` - The update time.

* `description` - The rule description.

* `application_type` - The application type.

* `protection_type` - The protection type.
