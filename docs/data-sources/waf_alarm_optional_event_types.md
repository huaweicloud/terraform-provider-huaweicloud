---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_alarm_optional_event_types"
description: |-
  Use this data source to query the alarm optional event types.
---

# huaweicloud_waf_alarm_optional_event_types

Use this data source to query the alarm optional event types.

## Example Usage

```hcl
data "huaweicloud_waf_alarm_optional_event_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `threats` - The list of optional event types in alarm notifications.

* `locale` - The description of each event type in alarm notifications.

  The [locale](#locale_struct) structure is documented below.

<a name="locale_struct"></a>
The `locale` block supports:

* `cmdi` - The command injection attack, in which an attacker injects malicious commands to perform unauthorized
  operations.

* `llm_prompt_injection` - The LLM prompt injection attack, in which an attacker constructs special inputs to tamper
  with the prompts of an AI model.

* `anticrawler` - The website anti-crawler policy, which is used to prevent automated programs from illegally obtaining
  website content.

* `custom_custom` - The precise protection, which is a custom security protection policy based on specific rules.

* `third_bot_river` - The third-party bot, which is an automated interaction program from a third-party service.

* `robot` - The malicious crawler, which is an automated program used to illegally obtain data or launch attacks.

* `custom_idc_ip` - The IDC intelligence, which is threat intelligence based on data center IP addresses.

* `antiscan_dir_traversal` - The directory traversal protection, which prevents attackers from accessing system files
  through directory traversal.

* `advanced_bot` - The advanced bot, which is an automated program with complex behavior patterns.

* `xss` - The XSS attack, in which an attacker obtains user information by injecting malicious scripts.

* `antiscan_high_freq_scan` - The scanning blocking, which identifies and blocks abnormal high-frequency requests.

* `webshell` - The web shells, which are malicious programs uploaded by attackers to remotely control websites.

* `cc` - The CC attack, a challenge attack, exhausting server resources by sending a large number of requests.

* `botm` - The BOT attacks, malicious attacks using automated programs.

* `illegal` - The invalid request, which violates security policies or service rules.

* `llm_prompt_sensitive` - The large model prompt word compliance detection, identifying sensitive information in
  prompts.

* `sqli` - The SQL injection: Attackers inject malicious SQL statements to obtain or tamper with data.

* `lfi` - The local file inclusion: Attackers exploit this vulnerability to include local files to obtain information.

* `antitamper` - The web Tamper Protection, Protecting Website Content from Unauthorized Modification.

* `custom_geoip` - The geolocation access control: geographical location-based access control policy.

* `rfi` - The remote file inclusion: Attackers exploit this vulnerability to execute malicious code using remote files.

* `vuln` - The other types of attacks, unclassified security vulnerabilities or attacks.

* `llm_response_sensitive` - The foundation model responds to compliance detection and identifies sensitive information
  in the AI model output.

* `custom_whiteblackip` - The IP address blacklist and whitelist, IP address-based access control policy.

* `leakage` - The website information leakage and accidental exposure of sensitive information.
