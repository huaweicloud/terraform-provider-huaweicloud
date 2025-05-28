---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_config"
description: |-
  Use this data source to get the available features under a region.
---

# huaweicloud_waf_config

Use this data source to get the available features under a region.

## Example Usage

```hcl
data "huaweicloud_waf_config" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `eps` - Whether EPS is supported.
  The value can be **true** or **false**.
  + **true**: EPS is supported.
  + **false**: EPS is not supported.

* `tls` - Whether to support the TLS version.
  The value can be **true** or **false**.
  + **true**: TLS version is supported.
  + **false**: TLS version is not supported.

* `ipv6` - Whether IPv6 protection is supported.
  The value can be **true** or **false**.
  + **true**: IPv6 protection is supported.
  + **false**: IPv6 protection is not supported.

* `alert` - Whether alarm reporting is supported.
  The value can be **true** or **false**.
  + **true**: Alarm reporting is supported.
  + **false**: Alarm reporting is not supported.

* `custom` - Whether precise protection is supported.
  The value can be **true** or **false**.
  + **true**: Precise protection is supported.
  + **false**: Precise protection is not supported.

* `elb_mode` - Whether ELB mode is supported.
  The value can be **true** or **false**.
  + **true**: The ELB mode is supported.
  + **false**: The ELB mode is not supported.

* `event_lts` - Whether LTS logging is supported.
  The value can be **true** or **false**.
  + **true**: LTS logging is supported.
  + **false**: LTS logging is not supported.

* `multi_dns` - Whether multi-DNS resolution is supported.
  The value can be **true** or **false**.
  + **true**: The multi-DNS resolution is supported.
  + **false**: The multi-DNS resolution is not supported.

* `search_ip` - Whether IP address search is supported.
  The value can be **true** or **false**.
  + **true**: IP address search is supported.
  + **false**: IP address search is not supported.

* `cc_enhance` - Whether CC attack protection is supported.
  The value can be **true** or **false**.
  + **true**: CC attack protection is supported.
  + **false**: CC attack protection is not supported.

* `cname_switch` - Whether CNAME switchover is supported.
  The value can be **true** or **false**.
  + **true**: CNAME switchover is supported.
  + **false**: CNAME switchover is not supported.

* `custom_block` - Whether custom block page is supported.
  The value can be **true** or **false**.
  + **true**: The custom block page is supported.
  + **false**: The custom block page is not supported.

* `advanced_ignore` - Whether false alarm masking is supported.
  The value can be **true** or **false**.
  + **true**: The false alarm masking is supported.
  + **false**: The false alarm masking is not supported.

* `js_crawler_enable` - Whether JS anti-crawler is supported.
  The value can be **true** or **false**.
  + **true**: The JS anti-crawler is supported.
  + **false**: The JS anti-crawler is not supported.

* `deep_decode_enable` - Whether deep inspection in basic web protection is supported.
  The value can be **true** or **false**.
  + **true**: The deep inspection is supported.
  + **false**: The deep inspection is not supported.

* `overview_bandwidth` - Whether security overview bandwidth statistics is supported.
  The value can be **true** or **false**.
  + **true**: The security overview bandwidth statistics is supported.
  + **false**: The security overview bandwidth statistics is not supported.

* `proxy_use_oldcname` - Whether old cname resolution is supported.
  The value can be **true** or **false**.
  + **true**: The old cname resolution is supported.
  + **false**: The old cname resolution is not supported.

* `check_all_headers_enable` - Whether all header inspection is supported.
  The value can be **true** or **false**.
  + **true**: The all eader inspection is supported.
  + **false**: The all header inspection is not supported.

* `geoip_enable` - Whether to support geolocation access control.
  The value can be **true** or **false**.
  + **true**: The geolocation access control is supported.
  + **false**: The geolocation access control is not supported.

* `load_balance_enable` - Whether to support domain name access load balancing.
  The value can be **true** or **false**.
  + **true**: The domain name access load balancing is supported.
  + **false**: The domain name access load balancing is not supported.

* `ipv6_protection_enable` - Whether IPv6 protection is supported.
  The value can be **true** or **false**.
  + **true**: The IPv6 protection is supported.
  + **false**: The IPv6 protection is not supported.

* `policy_sharing_enable` - Whether to support policy sharing.
  The value can be **true** or **false**.
  + **true**: The policy sharing is supported.
  + **false**: The policy sharing is not supported.

* `ip_group` - Whether the IP address group is supported.
  The value can be **true** or **false**.
  + **true**: The IP address group is supported.
  + **false**: The IP address group is not supported.

* `robot_action_enable` - Whether to support website anti-crawler.
  The value can be **true** or **false**.
  + **true**: The website anti-crawler is supported.
  + **false**: The website anti-crawler is not supported.

* `http2_enable` - Whether to support HTTP2.
  The value can be **true** or **false**.
  + **true**: The HTTP2 is supported.
  + **false**: The HTTP2 is not supported.

* `timeout_config_enable` - Whether to support the timeout configuration.
  The value can be **true** or **false**.
  + **true**: The timeout configuration is supported.
  + **false**: The timeout configuration is not supported.
