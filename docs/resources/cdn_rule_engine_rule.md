---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_rule_engine_rule"
description: |-
  Manages a rule engine rule resource within HuaweiCloud.
---

# huaweicloud_cdn_rule_engine_rule

Manages a rule engine rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "rule_name" {}

resource "huaweicloud_cdn_rule_engine_rule" "test" {
  domain_name = var.domain_name
  name        = var.rule_name
  status      = "on"
  priority    = 2

  conditions = jsonencode({
    "match": {
      "logic": "and",
      "criteria": [
        {
          "match_target_type": "extension",
          "match_type": "contains",
          "match_pattern": [".txt", ".png"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "scheme",
          "match_type": "contains",
          "match_pattern": ["HTTPS"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "method",
          "match_type": "contains",
          "match_pattern": ["PATCH"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "path",
          "match_type": "contains",
          "match_pattern": ["/test"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "arg",
          "match_target_name": "test",
          "match_type": "contains",
          "match_pattern": ["123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "filename",
          "match_type": "contains",
          "match_pattern": ["test", "123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "header",
          "match_target_name": "test",
          "match_type": "contains",
          "match_pattern": ["123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "clientip",
          "match_target_name": "connect",
          "match_type": "contains",
          "match_pattern": ["1.1.1.1"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "clientip_version",
          "match_target_name": "connect",
          "match_type": "contains",
          "match_pattern": ["IPv4"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "ua",
          "match_type": "contains",
          "match_pattern": ["test"],
          "negate": false,
          "case_sensitive": true
        }
      ]
    }
  })

  actions {
    cache_rule {
      ttl           = 10
      ttl_unit      = "m"
      follow_origin = "min_ttl"
      force_cache   = "off"
    }
  }

  actions {
    access_control {
      type = "block"
    }
  }

  actions {
    http_response_header {
      name   = "Access-Control-Deny-Origin"
      value  = "*"
      action = "delete"
    }
  }

  actions {
    browser_cache_rule {
      cache_type = "follow_origin"
    }
  }

  actions {
    request_url_rewrite {
      execution_mode = "break"
      redirect_url   = "/path/$1"
    }
  }

  actions {
    flexible_origin {
      sources_type      = "ipaddr"
      ip_or_domain      = "1.1.1.1"
      http_port         = 80
      https_port        = 443
      origin_protocol  = "follow"
      host_name         = "1.1.1.1"
      priority          = 2
      weight            = 10
    }
    flexible_origin {
      sources_type      = "third_bucket"
      ip_or_domain      = "test.third-bucket.com"
      bucket_access_key = "test-ak"
      bucket_secret_key = "test-sk"
      bucket_region     = "cn-north-4"
      bucket_name       = "test-third-bucket-name"
      http_port         = 80
      https_port        = 443
      origin_protocol   = "follow"
      host_name         = "1.1.1.1"
      priority          = 1
      weight            = 2
    }
  }

  actions {
    origin_request_header {
      action = "delete"
      name   = "test"
      value  = "123"
    }
  }

  actions {
    origin_request_url_rewrite {
      rewrite_type = "simple"
      target_url   = "/test"
    }
  }

  actions {
    origin_range {
      status = "on"
    }
  }

  actions {
    request_limit_rule {
      limit_rate_after = 2
      limit_rate_value = 3
    }
  }

  actions {
    error_code_cache {
      code = 403
      ttl  = 123
    }

    error_code_cache {
      code = 404
      ttl  = 123
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the rule engine rule is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, NonUpdatable) Specifies the accelerated domain name to which the rule engine rule
  belongs.  
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the rule engine rule.  
  The valid length is limit from `1` to `50`.

* `status` - (Required, String) Whether to enable the rule engine rule.  
  The valid values are as follows:
  + **on**
  + **off**

* `priority` - (Required, Int) Specifies the priority of the rule engine rule.
  The valid value is range from `1` to `100`.

* `conditions` - (Optional, String) Specifies the trigger conditions of the rule engine rule, in JSON format.

* `actions` - (Required, List) Specifies the list of actions to be performed when the rule engine rule is met.
  The [actions](#cdn_rule_engine_rule_actions) structure is documented below.

<a name="cdn_rule_engine_rule_actions"></a>
The `actions` block supports:

* `flexible_origin` - (Optional, List) Specifies the list of flexible origin configurations.  
  The [flexible_origin](#cdn_rule_engine_rule_actions_flexible_origin) structure is documented below.

* `origin_request_header` - (Optional, List) Specifies the list of origin request header configurations.  
  The [origin_request_header](#cdn_rule_engine_rule_actions_origin_request_header) structure is documented below.

* `http_response_header` - (Optional, List) Specifies the list of HTTP response header configurations.  
  The [http_response_header](#cdn_rule_engine_rule_actions_http_response_header) structure is documented below.

* `access_control` - (Optional, List) Specifies the access control configuration.  
  The [access_control](#cdn_rule_engine_rule_actions_access_control) structure is documented below.

* `request_limit_rule` - (Optional, List) Specifies the request rate limit configuration.  
  The [request_limit_rule](#cdn_rule_engine_rule_actions_request_limit_rule) structure is documented below.

* `origin_request_url_rewrite` - (Optional, List) Specifies the origin request URL rewrite configuration.  
  The [origin_request_url_rewrite](#cdn_rule_engine_rule_actions_origin_request_url_rewrite) structure is documented
  below.

* `cache_rule` - (Optional, List) Specifies the cache rule configuration.  
  The [cache_rule](#cdn_rule_engine_rule_actions_cache_rule) structure is documented below.

* `request_url_rewrite` - (Optional, List) Specifies the access URL rewrite configuration.  
  The [request_url_rewrite](#cdn_rule_engine_rule_actions_request_url_rewrite) structure is documented below.

* `browser_cache_rule` - (Optional, List) Specifies the browser cache rule configuration.  
  The [browser_cache_rule](#cdn_rule_engine_rule_actions_browser_cache_rule) structure is documented below.

* `error_code_cache` - (Optional, List) Specifies the list of error code cache configurations.  
  The [error_code_cache](#cdn_rule_engine_rule_actions_error_code_cache) structure is documented below.

* `origin_range` - (Optional, List) Specifies the origin range configuration.  
  The [error_code_cache](#cdn_rule_engine_rule_actions_origin_range) structure is documented below.

-> Each of the above configuration items must be declared in a separate actions structure.

<a name="cdn_rule_engine_rule_actions_flexible_origin"></a>
The `flexible_origin` block supports:

* `sources_type` - (Required, String) Specifies the source type.  
  The valid values are as follows:
  + **ipaddr**
  + **domain**
  + **obs_bucket**
  + **third_bucket**

* `ip_or_domain` - (Required, String) Specifies the origin IP or domain name.

* `priority` - (Required, Int) Specifies the origin priority.  
  The valid value is range from `1` to `100`.

* `weight` - (Required, Int) Specifies the origin weight.  
  The valid value is range from `1` to `100`.

* `obs_bucket_type` - (Optional, String) Specifies the OBS bucket type.  
  The valid values are as follows:
  + **private**
  + **public**

* `bucket_access_key` - (Optional, String) Specifies the third-party object storage access key.

* `bucket_secret_key` - (Optional, String) Specifies the third-party object storage secret key.

* `bucket_region` - (Optional, String) Specifies the third-party object storage region.

* `bucket_name` - (Optional, String) Specifies the third-party object storage name.

* `host_name` - (Optional, String) Specifies the origin host name.

* `origin_protocol` - (Optional, String) Specifies the origin protocol.  
  The valid values are as follows:
  + **follow**
  + **http**
  + **https**

* `http_port` - (Optional, Int) Specifies the HTTP port number.  
  The valid value is range from `1` to `65,535`.

* `https_port` - (Optional, Int) Specifies the HTTPS port number.  
  The valid value is range from `1` to `65,535`.

<a name="cdn_rule_engine_rule_actions_origin_request_header"></a>
The `origin_request_header` block supports:

* `name` - (Required, String) Specifies the back-to-origin request header parameter name.

* `action` - (Required, String) Specifies the back-to-origin request header setting type.  
  The valid values are as follows:
  + **delete**
  + **set**

* `value` - (Optional, String) Specifies the back-to-origin request header parameter value.

<a name="cdn_rule_engine_rule_actions_http_response_header"></a>
The `http_response_header` block supports:

* `name` - (Required, String) Specifies the HTTP response header parameter name.

* `action` - (Required, String) Specifies the operation type of setting HTTP response header.  
  The valid values are as follows:
  + **set**
  + **delete**

* `value` - (Optional, String) Specifies the HTTP response header parameter value.

<a name="cdn_rule_engine_rule_actions_access_control"></a>
The `access_control` block supports:

* `type` - (Required, String) Specifies the access control type.  
  The valid values are as follows:
  + **block**
  + **trust**

<a name="cdn_rule_engine_rule_actions_request_limit_rule"></a>
The `request_limit_rule` block supports:

* `limit_rate_after` - (Required, Int) Specifies the rate limit condition.

* `limit_rate_value` - (Required, Int) Specifies the rate limit value.

<a name="cdn_rule_engine_rule_actions_origin_request_url_rewrite"></a>
The `origin_request_url_rewrite` block supports:

* `rewrite_type` - (Required, String) Specifies the rewrite type.  
  The valid values are as follows:
  + **simple**
  + **wildcard**
  + **regex**

* `target_url` - (Required, String) Specifies the target URL.

* `source_url` - (Optional, String) Specifies the source URL to be rewritten.

  -> This parameter is required if the value of parameter `rewrite_type` is **wildcard**.

<a name="cdn_rule_engine_rule_actions_cache_rule"></a>
The `cache_rule` block supports:

* `ttl` - (Required, Int) Specifies the cache expiration time.

* `ttl_unit` - (Required, String) Specifies the cache expiration time unit.  
  The valid values are as follows:
  + **s**
  + **m**
  + **h**
  + **d**

* `follow_origin` - (Required, String) Specifies the cache expiration time source.  
  The valid values are as follows:
  + **off**
  + **on**
  + **min_ttl**

* `force_cache` - (Optional, String) Whether to enable forced caching.  
  The valid values are as follows:
  + **on**
  + **off**

<a name="cdn_rule_engine_rule_actions_request_url_rewrite"></a>
The `request_url_rewrite` block supports:

* `execution_mode` - (Required, String) Specifies the execution mode.  
  The valid values are as follows:
  + **redirect**
  + **break**

* `redirect_url` - (Required, String) Specifies the redirect URL.

* `redirect_status_code` - (Optional, Int) Specifies the redirect status code.  
  The valid values are as follows:
  + `301`
  + `302`
  + `303`
  + `307`

* `redirect_host` - (Optional, String) Specifies the redirect host.

-> Parameter `redirect_status_code` and `redirect_host` can only be configured when the value of `cache_type` is
   **redirect**.

<a name="cdn_rule_engine_rule_actions_browser_cache_rule"></a>
The `browser_cache_rule` block supports:

* `cache_type` - (Required, String) Specifies the cache effective type.  
  The valid values are as follows:
  + **follow_origin**
  + **ttl**
  + **never**

* `ttl` - (Optional, Int) Specifies the cache expiration time.

* `ttl_unit` - (Optional, String) Specifies the cache expiration time unit.  
  The valid values are as follows:
  + **s**
  + **m**
  + **h**
  + **d**

-> Parameter `ttl` and `ttl_unit` can only be configured when the value of `cache_type` is **ttl**.

<a name="cdn_rule_engine_rule_actions_error_code_cache"></a>
The `error_code_cache` block supports:

* `code` - (Required, Int) Specifies the error code to be cached.  
  The valid values are as follows:
  + `301`
  + `302`
  + `400`
  + `401`
  + `403`
  + `404`
  + `405`
  + `407`
  + `414`
  + `451`
  + `500`
  + `501`
  + `502`
  + `503`
  + `504`
  + `509`
  + `514`

* `ttl` - (Required, Int) Specifies the error code cache time.

<a name="cdn_rule_engine_rule_actions_origin_range"></a>
The `origin_range` block supports:

* `status` - (Required, String) Specifies the origin range status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The rule engine rule can be imported using the format `<domain_name>/<id>` or `<domain_name>/<name>`, e.g.

```bash
$ terraform import huaweicloud_cdn_rule_engine_rule.test <domain_name>/<id>
```

or

```bash
$ terraform import huaweicloud_cdn_rule_engine_rule.test <domain_name>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `actions.*.flexible_origin.*.bucket_secret_key`.
It is generally recommended running `terraform plan` after importing rule engine rule.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cdn_rule_engine_rule" "test" {
  ...

  lifecycle {
    ignore_changes = [
      actions.5.flexible_origin.1.bucket_secret_key,
    ]
  }
}
```
