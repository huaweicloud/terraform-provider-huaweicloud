---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_tag_threat_map"
description: |-
  Use this data source to query the event types supported by alarm notifications.
---

# huaweicloud_waf_tag_threat_map

Use this data source to query the event types supported by alarm notifications.

## Example Usage

```hcl
data "huaweicloud_waf_tag_threat_map" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `threats` - The event types avaiables in alarm notifications.

* `locale` - The description of each event type in alarm notifications.

  The [locale](locale_struct) structure is documented below.

<a name="locale_struct"></a>
The `locale` block supports:

* `cmdi` - Indicates command injection attack, in which an attacker injects malicious commands to perform unauthorized
  operations.

* `llm_prompt_injection` - Indicates LLM prompt injection attack, in which an attacker constructs special inputs to
  tamper with the prompts of an AI model.

* `anticrawler` - Indicates website anti-crawler policy, which is used to prevent automated programs from illegally
  obtaining website content.

* `custom_custom` - Indicates precise protection, which is a custom security protection policy based on specific rules.

* `third_bot_river` - Indicates Third-party BOT, which is an automated interaction program from a third-party service.

* `robot` - Indicates malicious crawler, which is an automated program used to illegally obtain data or launch attacks.

* `custom_idc_ip` - Indicates IDC intelligence, which is threat intelligence based on data center IP addresses.

* `antiscan_dir_traversal` - Indicates directory traversal protection, which prevents attackers from accessing system
  files through directory traversal.

* `advanced_bot` - Indicates advanced BOT, which is an automated program with complex behavior patterns.

* `xss` - Indicates XSS attack, in which an attacker obtains user information by injecting malicious scripts.

* `antiscan_high_freq_scan` - Indicates scanning blocking, which identifies and blocks abnormal high-frequency
  requests.

* `webshell` - Indicates web shells, which are malicious programs uploaded by attackers to remotely control websites.

* `cc` - Indicates CC attack, a challenge attack, exhausting server resources by sending a large number of requests.

* `botm` - Indicates BOT attacks, malicious attacks using automated programs.

* `illegal` - Indicates invalid request, which violates security policies or service rules.

* `llm_prompt_sensitive` - Indicates large model prompt word compliance detection, identifying sensitive information
  in prompts.

* `sqli` - Indicates SQL injection, attackers inject malicious SQL statements to obtain or tamper with data.

* `lfi` - Indicates local file inclusion, attackers exploit this vulnerability to include local files to
  obtain information.

* `antitamper` - Indicates web tamper protection, protecting website content from unauthorized modification.

* `custom_geoip` - Indicates geolocation access control, geographical location-based access control policy.

* `rfi` - Indicates remote file inclusion, attackers exploit this vulnerability to execute malicious code
  using remote files.

* `vuln` - Indicates other types of attacks, unclassified security vulnerabilities or attacks.

* `llm_response_sensitive` - Indicates foundation model responds to compliance detection and identifies sensitive
  information in the AI model output.

* `custom_whiteblackip` - Indicates IP address blacklist and whitelist, IP address-based access control policy.

* `leakage` - Indicates website information leakage and accidental exposure of sensitive information.
