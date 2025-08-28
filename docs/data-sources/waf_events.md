---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_events"
description: |-
  Use this data source to get the list of WAF attack events.
---

# huaweicloud_waf_events

Use this data source to get the list of WAF attack events.

## Example Usage

```hcl
variable "recent" {}

data "huaweicloud_waf_events" "test" {
  recent = var.recent
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `recent` - (Optional, String) Specifies the time range for querying logs.
  The value can be **yesterday**, **today**, **3days**, **1week** or **1month**.

* `from` - (Optional, Int) Specifies the start time.
  The format is 13-digit timestamp in millisecond.

  -> The field `from` must be used together with field `to`.

* `to` - (Optional, Int) Specifies the end time.
  The format is 13-digit timestamp in millisecond.

  -> The field `to` must be used together with field `from`.

-> Exactly one of `recent` or `from`, `to` must be set. If both `recent` and `from`, `to` are set, `recent`
  takes effect.

* `attacks` - (Optional, List) Specifies the attack type.
  The valid values are as follows:
  + **vuln**: Other attack types.
  + **sqli**: SQL injections.
  + **lfi**: Local file inclusion attacks.
  + **cmdi**: Command injections.
  + **xss**: XSS attacks.
  + **robot**: Malicious crawlers.
  + **rfi**: Remote file inclusion attacks.
  + **custom_custom**: Precise protection.
  + **cc**: CC attacks.
  + **webshell**: Website Trojans.
  + **custom_whiteblackip**: Attacks blocked based on blacklist and whitelist settings.
  + **custom_geoip**: Attacks blocked based on geolocations.
  + **antitamper**: Anti-tamper events.
  + **anticrawler**: Anti-crawler events.
  + **leakage**: Website data leakage prevention.
  + **illegal**: Unauthorized request.
  + **antiscan_high_freq_scan**: Blocking of high-frequency scanning.
  + **antiscan_dir_traversal**: Directory traversal protection.

* `hosts` - (Optional, List) Specifies the domain ID list.

* `sips` - (Optional, List) Specifies the source IP addresses.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The attack event list.
  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The event ID.

* `time` - The timestamp when the attack occurred, in milliseconds.

* `policyid` - The policy ID.

* `sip` - The source IP address.

* `host` - The domain name.

* `url` - The attacked URL link.

* `attack` - The attack type.

* `rule` - The hit rule ID.

* `payload` - The hit payload.

* `payload_location` - The hit payload location.

* `action` - The protection action.

* `request_line` - The request method and path.

* `headers` - The HTTP request headers.

* `cookie` - The request cookie.

* `status` - The response code status.

* `process_time` - The processing time.

* `region` - The geographic location.

* `host_id` - The domain ID.

* `response_time` - The response time.

* `response_size` - The response body size.

* `response_body` - The response body.

* `request_body` - The request body.
