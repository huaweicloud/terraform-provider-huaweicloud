---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_rule_engine_rules"
description: |-
  Use this data source to get a list of rule engine rules within HuaweiCloud.
---

# huaweicloud_cdn_rule_engine_rules

Use this data source to get a list of rule engine rules within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_cdn_rule_engine_rules" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the accelerated domain name to which the rule engine rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the rule engine rules.  
  The [rules](#cdn_rule_engine_rules_rules) structure is documented below.

<a name="cdn_rule_engine_rules_rules"></a>
The `rules` block supports:

* `id` - The ID of the rule engine rule.

* `name` - The name of the rule engine rule.

* `status` - Whether the rule is enabled.

* `priority` - The priority of the rule engine rule.

* `conditions` - The trigger conditions of the current rule, in JSON format.

* `actions` - The list of actions to be performed when the rules are met.  
  The [actions](#cdn_rule_engine_rules_rules_actions) structure is documented below.

<a name="cdn_rule_engine_rules_rules_actions"></a>
The `actions` block supports:

* `flexible_origin` - The list of flexible origin configurations.  
  The [flexible_origin](#cdn_rule_engine_rules_rules_actions_flexible_origin) structure is documented below.

* `origin_request_header` - The list of origin request header configurations.  
  The [origin_request_header](#cdn_rule_engine_rules_rules_actions_origin_request_header) structure is documented below.

* `http_response_header` - The list of HTTP response header configurations.  
  The [http_response_header](#cdn_rule_engine_rules_rules_actions_http_response_header) structure is documented below.

* `access_control` - The access control configuration.  
  The [access_control](#cdn_rule_engine_rules_rules_actions_access_control) structure is documented below.

* `request_limit_rule` - The request rate limit configuration.  
  The [request_limit_rule](#cdn_rule_engine_rules_rules_actions_request_limit_rule) structure is documented below.

* `origin_request_url_rewrite` - The origin request URL rewrite configuration.  
  The [origin_request_url_rewrite](#cdn_rule_engine_rules_rules_actions_origin_request_url_rewrite) structure is
  documented below.

* `cache_rule` - The cache rule configuration.  
  The [cache_rule](#cdn_rule_engine_rules_rules_actions_cache_rule) structure is documented below.

* `request_url_rewrite` - The access URL rewrite configuration.  
  The [request_url_rewrite](#cdn_rule_engine_rules_rules_actions_request_url_rewrite) structure is documented below.

* `browser_cache_rule` - The browser cache rule configuration.  
  The [browser_cache_rule](#cdn_rule_engine_rules_rules_actions_browser_cache_rule) structure is documented below.

* `error_code_cache` - The list of error code cache configurations.  
  The [error_code_cache](#cdn_rule_engine_rules_rules_actions_error_code_cache) structure is documented below.

* `origin_range` - The origin range configuration.  
  The [origin_range](#cdn_rule_engine_rules_rules_actions_origin_range) structure is documented below.

<a name="cdn_rule_engine_rules_rules_actions_flexible_origin"></a>
The `flexible_origin` block supports:

* `sources_type` - The source type.

* `ip_or_domain` - The origin IP or domain name.

* `priority` - The origin priority.

* `weight` - The origin weight.

* `obs_bucket_type` - The OBS bucket type.

* `bucket_access_key` - The third-party object storage access key.

* `bucket_region` - The third-party object storage region.

* `bucket_name` - The third-party object storage name.

* `host_name` - The origin host name.

* `origin_protocol` - The origin protocol.

* `http_port` - The HTTP port number.

* `https_port` - The HTTPS port number.

<a name="cdn_rule_engine_rules_rules_actions_origin_request_header"></a>
The `origin_request_header` block supports:

* `name` - The back-to-origin request header parameter name.

* `action` - The back-to-origin request header setting type.

* `value` - The back-to-origin request header parameter value.

<a name="cdn_rule_engine_rules_rules_actions_http_response_header"></a>
The `http_response_header` block supports:

* `name` - The HTTP response header parameter name.

* `action` - The operation type of setting HTTP response header.

* `value` - The HTTP response header parameter value.

<a name="cdn_rule_engine_rules_rules_actions_access_control"></a>
The `access_control` block supports:

* `type` - The access control type.

<a name="cdn_rule_engine_rules_rules_actions_request_limit_rule"></a>
The `request_limit_rule` block supports:

* `limit_rate_after` - The rate limit condition.

* `limit_rate_value` - The rate limit value.

<a name="cdn_rule_engine_rules_rules_actions_origin_request_url_rewrite"></a>
The `origin_request_url_rewrite` block supports:

* `rewrite_type` - The rewrite type.

* `target_url` - The target URL.

* `source_url` - The source URL to be rewritten.

<a name="cdn_rule_engine_rules_rules_actions_cache_rule"></a>
The `cache_rule` block supports:

* `ttl` - The cache expiration time.

* `ttl_unit` - The cache expiration time unit.

* `follow_origin` - The cache expiration time source.

* `force_cache` - Whether to enable forced caching.

<a name="cdn_rule_engine_rules_rules_actions_request_url_rewrite"></a>
The `request_url_rewrite` block supports:

* `redirect_url` - The redirect URL.

* `execution_mode` - The execution mode.

* `redirect_status_code` - The redirect status code.

* `redirect_host` - The redirect host.

<a name="cdn_rule_engine_rules_rules_actions_browser_cache_rule"></a>
The `browser_cache_rule` block supports:

* `cache_type` - The cache effective type.

* `ttl` - The cache expiration time.

* `ttl_unit` - The cache expiration time unit.

<a name="cdn_rule_engine_rules_rules_actions_error_code_cache"></a>
The `error_code_cache` block supports:

* `code` - The error code to be cached.

* `ttl` - The error code cache time.

<a name="cdn_rule_engine_rules_rules_actions_origin_range"></a>
The `origin_range` block supports:

* `status` - The origin range status.
