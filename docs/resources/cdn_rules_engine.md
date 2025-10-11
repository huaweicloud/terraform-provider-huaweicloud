# huaweicloud_cdn_rules_engine

The Huawei Cloud Content Delivery Network (CDN) Rules Engine resource is used to create and manage CDN rules engine rules.

## Overview

The rules engine feature enables flexible and fine-grained rule configuration through a graphical interface. By restricting trigger conditions, you can control the scope of resources where the configuration takes effect, meeting the needs of various scenarios.

**Note:** Please submit a service ticket to enable the rules engine feature before using this resource.

## Supported API Endpoints

- **Create Rule**: `POST /v1.0/cdn/configuration/domains/{domain_name}/rules`
- **List Rules**: `GET /v1.0/cdn/configuration/domains/{domain_name}/rules`
- **Full Update Rule**: `PUT /v1.0/cdn/configuration/domains/{domain_name}/rules/full-update`
- **Batch Update Rule Status and Priority**: `POST /v1.0/cdn/configuration/domains/{domain_name}/rules/batch-update`
- **Delete Rule**: `DELETE /v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}`

## Example Usage

### Basic Usage

```hcl
resource "huaweicloud_cdn_rules_engine" "example" {
  domain_name = "example.com"
  name        = "example_rule"
  status      = "on"
  priority    = 1

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "path"
        match_type        = "contains"
        match_pattern     = ["/admin"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    access_control {
      type = "block"
    }
  }
}
```

### Access Control Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "access_control" {
  domain_name = "example.com"
  name        = "access_control_rule"
  status      = "on"
  priority    = 1

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "path"
        match_type        = "contains"
        match_pattern     = ["/admin", "/private"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    access_control {
      type = "block"
    }
  }
}
```

### Cache Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "cache_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "cache_rule"
  status      = "on"
  priority    = 2

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "extension"
        match_type        = "contains"
        match_pattern     = [".jpg", ".png", ".gif"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    cache_rule {
      ttl          = 3600
      ttl_unit     = "s"
      follow_origin = "off"
      force_cache   = "on"
    }
  }
}
```

### URL Rewrite Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "url_rewrite_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "url_rewrite_rule"
  status      = "on"
  priority    = 3

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "path"
        match_type        = "contains"
        match_pattern     = ["/old-path/*"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    request_url_rewrite {
      redirect_status_code = 301
      redirect_url         = "/new-path/$1"
      execution_mode       = "redirect"
    }
  }
}
```

### Browser Cache Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "browser_cache_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "browser_cache_rule"
  status      = "on"
  priority    = 4

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "extension"
        match_type        = "contains"
        match_pattern     = ["css", "js"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    browser_cache_rule {
      cache_type = "ttl"
      ttl        = 86400
      ttl_unit   = "s"
    }
  }
}
```

### HTTP Response Header Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "http_response_header_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "http_response_header_rule"
  status      = "on"
  priority    = 5

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "path"
        match_type        = "contains"
        match_pattern     = ["/api/*"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    http_response_header {
      name   = "Access-Control-Allow-Origin"
      value  = "*"
      action = "set"
    }
    
    http_response_header {
      name   = "Cache-Control"
      value  = "no-cache"
      action = "set"
    }
  }
}
```

### Request Limit Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "request_limit_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "request_limit_rule"
  status      = "on"
  priority    = 6

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "clientip"
        match_target_name = "xff"
        match_type        = "contains"
        match_pattern     = ["192.168.1.100"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    request_limit_rules {
      limit_rate_after  = 1024
      limit_rate_value  = 1024
    }
  }
}
```

### Error Code Cache Rule

```hcl
resource "huaweicloud_cdn_rules_engine" "error_code_cache_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "error_code_cache_rule"
  status      = "on"
  priority    = 7

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "path"
        match_type        = "contains"
        match_pattern     = ["/static/*"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    error_code_cache {
      code = 404
      ttl  = 60
    }

    error_code_cache {
      code = 403
      ttl  = 60
    }
  }
}
```

## Argument Reference

### Required Arguments

| Name | Type | Description |
|------|------|-------------|
| `domain_name` | string | The acceleration domain name. |
| `name` | string | The rule name, 1-50 characters. |
| `status` | string | Rule status, valid values: `on`, `off`. |
| `priority` | int | Rule priority, 1-100, higher value means higher priority. |
| `conditions` | list | The trigger conditions for the rule. |
| `actions` | list | The actions to execute when the rule is triggered. |

### Condition Block (conditions)

| Name | Type | Description |
|------|------|-------------|
| `match` | list | Rule match conditions. |

#### match Block

| Name | Type | Description |
|------|------|-------------|
| `logic` | string | Logical operator, valid values: `and`, `or`. |
| `criteria` | list | List of match criteria. |

#### criteria Block

| Name | Type | Description |
|------|------|-------------|
| `match_target_type` | string | Match target type, supports: `schema`, `method`, `path`, `arg`, `extension`, `filename`, `header`, `clientip`, `clientip_version`, `ua`, `ngx_variable`. |
| `match_target_name` | string | Match target name (optional). |
| `match_type` | string | Match algorithm, currently only supports `contains`. |
| `match_pattern` | list | List of match patterns. |
| `negate` | bool | Whether to negate, default is `false`. |
| `case_sensitive` | bool | Whether case sensitive, default is `false`. |

### Action Block (actions)

#### access_control

| Name | Type | Description |
|------|------|-------------|
| `type` | string | Access control type, valid values: `block`, `trust`. |

#### cache_rule

| Name | Type | Description |
|------|------|-------------|
| `ttl` | int | Cache expiration time on CDN nodes. |
| `ttl_unit` | string | Cache expiration time unit, valid values: `s`, `m`, `h`, `d`. |
| `follow_origin` | string | Source of cache expiration time, valid values: `on`, `off`, `min_ttl`. |
| `force_cache` | string | Whether to enable forced cache, valid values: `on`, `off`. |

#### request_url_rewrite

| Name | Type | Description |
|------|------|-------------|
| `redirect_status_code` | int | Redirect status code, valid values: `301`, `302`, `303`, `307`. |
| `redirect_url` | string | Redirect URL. |
| `redirect_host` | string | Redirect host (optional). |
| `execution_mode` | string | Execution mode, valid values: `redirect`, `break`. |

#### browser_cache_rule

| Name | Type | Description |
|------|------|-------------|
| `cache_type` | string | Cache effective type, valid values: `follow_origin`, `ttl`, `never`. |
| `ttl` | int | Cache expiration time (required if cache_type is `ttl`). |
| `ttl_unit` | string | Cache expiration time unit, valid values: `s`, `m`, `h`, `d`. |

#### http_response_header

| Name | Type | Description |
|------|------|-------------|
| `name` | string | HTTP response header name. |
| `value` | string | HTTP response header value (optional). |
| `action` | string | HTTP response header operation, valid values: `set`, `delete`. |

#### origin_request_header

| Name | Type | Description |
|------|------|-------------|
| `name` | string | Origin request header name. |
| `value` | string | Origin request header value (optional). |
| `action` | string | Origin request header operation, valid values: `set`, `delete`. |

#### request_limit_rules

| Name | Type | Description |
|------|------|-------------|
| `limit_rate_after` | int | Rate limit condition, in bytes. |
| `limit_rate_value` | int | Rate limit value, in Bps. |

#### error_code_cache

| Name | Type | Description |
|------|------|-------------|
| `code` | int | Error code to cache. |
| `ttl` | int | Error code cache time, in seconds. |

### Computed Attributes

| Name | Type | Description |
|------|------|-------------|
| `rule_id` | string | The rule ID. |

## Import

You can import an existing resource using the following format:

```bash
terraform import huaweicloud_cdn_rules_engine.example <domain_name>/<rule_id>
```

## Notes

1. **Feature Enablement**: You must submit a service ticket to enable the rules engine feature before use.
2. **Priority**: Priorities must not be duplicated. Higher values indicate higher priority.
3. **Rule Name**: Rule name must be 1-50 characters.
4. **Match Conditions**: Supports various match target types, configure as needed.
5. **Action Configuration**: Each rule can configure multiple actions to meet complex business needs.
6. **Update Operation**: The update operation uses the full update API and will overwrite all configuration.

## Supported Match Target Types

- `schema`: Protocol type used by the client request (HTTP, HTTPS)
- `method`: HTTP method used by the client request (GET, PUT, POST, DELETE, HEAD, OPTIONS, PATCH, TRACE, CONNECT)
- `path`: URL path of the client request
- `arg`: Query parameter in the client request URL
- `extension`: File extension of the requested content
- `filename`: File name of the requested content
- `header`: HTTP request header
- `clientip`: Client IP address
- `clientip_version`: Client IP version (IPv4, IPv6)
- `ua`: User-Agent in the client request header
- `ngx_variable`: Nginx variable 