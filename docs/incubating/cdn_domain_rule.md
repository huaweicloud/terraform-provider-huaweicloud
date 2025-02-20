---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_rule"
description: |-
  Manages a CDN domain rule resource within HuaweiCloud.
---

# huaweicloud_cdn_domain_rule

Manages a CDN domain rule resource within HuaweiCloud.

## Example Usage

### Create a CDN domain rule

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_domain_rule" "test" {
  name = var.domain_name

  rules {
    name     = "test-rule"
    priority = 1
    status   = "on"

    actions {
      cache_rule {
        follow_origin = "on"
        force_cache   = "on"
        ttl           = 30
        ttl_unit      = "d"
      }
    }

    conditions {
      match {
        criteria = jsonencode(
          [
            {
              case_sensitive = false
              match_pattern = [
                "HTTP",
              ]
              match_target_name = ""
              match_target_type = "scheme"
              match_type        = "contains"
              negate            = false
            },
            {
              case_sensitive = false
              match_pattern = [
                "GET",
              ]
              match_target_name = ""
              match_target_type = "method"
              match_type        = "contains"
              negate            = false
            },
          ]
        )
        logic = "and"
      }
    }
  }
}
```

### Create several CDN domain rules

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_domain_rule" "test" {
  name = var.domain_name

  rules {
    name     = "test-rule2"
    priority = 2
    status   = "on"

    actions {
      access_control {
        type = "block"
      }
    }
    actions {
      cache_rule {
        follow_origin = "on"
        force_cache   = "off"
        ttl           = 30
        ttl_unit      = "d"
      }
    }
    actions {
      flexible_origin {
        host_name       = "test.name.com.cn"
        http_port       = 80
        https_port      = 443
        ip_or_domain    = "110.110.110.11"
        obs_bucket_type = null
        origin_protocol = "follow"
        priority        = 3
        sources_type    = "ipaddr"
        weight          = 4
      }
    }
    actions {
      http_response_header {
        action = "set"
        name   = "Content-Language"
        value  = "en-Us"
      }
    }
    actions {
      origin_request_header {
        action = "set"
        name   = "X-Token"
        value  = "abc"
      }
    }
    actions {
      origin_request_url_rewrite {
        rewrite_type = "simple"
        source_url   = null
        target_url   = "/test/*.jpg"
      }
    }
    actions {
      request_url_rewrite {
        execution_mode       = "break"
        redirect_host        = null
        redirect_status_code = 0
        redirect_url         = "/index/test.html"
      }
    }

    conditions {
      match {
        criteria = jsonencode(
          [
            {
              case_sensitive = false
              criteria = [
                {
                  case_sensitive = false
                  match_pattern = [
                    "HTTP",
                  ]
                  match_target_name = ""
                  match_target_type = "scheme"
                  match_type        = "contains"
                  negate            = false
                },
              ]
              logic  = "and"
              negate = false
            },
            {
              case_sensitive = false
              criteria = [
                {
                  case_sensitive = true
                  match_pattern = [
                    "HTTP",
                  ]
                  match_target_name = "aa"
                  match_target_type = "arg"
                  match_type        = "contains"
                  negate            = false
                },
                {
                  case_sensitive = false
                  criteria = [
                    {
                      case_sensitive = true
                      match_pattern = [
                        "HTTP",
                      ]
                      match_target_name = ""
                      match_target_type = "filename"
                      match_type        = "contains"
                      negate            = false
                    },
                    {
                      case_sensitive = true
                      match_pattern = [
                        "HTTP",
                      ]
                      match_target_name = ""
                      match_target_type = "ua"
                      match_type        = "contains"
                      negate            = true
                    },
                  ]
                  logic  = "or"
                  negate = false
                },
              ]
              logic  = "and"
              negate = false
            },
          ]
        )
        logic = "or"
      }
    }
  }
  rules {
    name     = "test-rule1"
    priority = 1
    status   = "on"

    actions {
      access_control {
        type = "trust"
      }
    }
    actions {
      cache_rule {
        follow_origin = "on"
        force_cache   = "on"
        ttl           = 30
        ttl_unit      = "d"
      }
    }

    conditions {
      match {
        criteria = jsonencode(
          [
            {
              case_sensitive = false
              match_pattern = [
                "HTTP",
              ]
              match_target_name = ""
              match_target_type = "scheme"
              match_type        = "contains"
              negate            = false
            },
            {
              case_sensitive = false
              match_pattern = [
                "GET",
              ]
              match_target_name = ""
              match_target_type = "method"
              match_type        = "contains"
              negate            = false
            },
          ]
        )
        logic = "and"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies acceleration domain name.

* `rules` - (Required, List) Specifies a list of domain rules.
  The [rules](#domain_rule_rules) structure is documented below.

<a name="domain_rule_rules"></a>
The `rules` block supports:

* `name` - (Required, String) Specifies the rule name. The valid length is limit from `1` to `50`.

* `status` - (Required, String) Specifies the rule status. Valid values are **on** and **off**.

* `priority` - (Required, Int) Specifies the rule priority. The rule with larger value has higher priority.
  The value ranges from `1` to `100`, and the priority cannot be the same.

* `conditions` - (Required, List) Specifies the trigger conditions. The length of this field is `1`.
  The [conditions](#rules_conditions) structure is documented below.

* `actions` - (Required, List) Specifies a list of actions to be performed when the rules are met.
  The [actions](#rules_actions) structure is documented below.

<a name="rules_conditions"></a>
The `conditions` block supports:

* `match` - (Required, List) Specifies the value of the matching condition. The length of this field is `1`.
  The [match](#conditions_match) structure is documented below.

<a name="conditions_match"></a>
The `match` block supports:

* `logic` - (Required, String) Specifies the logical operator. Valid values are **and** with **or**.

* `criteria` - (Required, String) Specifies the match criteria list in JSON format. This field is a nested structure.
  Please refer to the usage of this field in the document example. Example without nested structure: `[{"case_sensitive"
  :false,"match_pattern":["HTTP"],"match_target_name":"","match_target_type":"scheme","match_type":"contains",
  "negate":false}]`. Example with nested structure: `[{"case_sensitive":false,"match_pattern":["HTTP"],
  "match_target_name":"","match_target_type":"scheme","match_type":"contains","negate":false},{"case_sensitive":false,
  "criteria":[{"case_sensitive":false,"match_pattern":["GET"],"match_target_name":"","match_target_type":"method",
  "match_type":"contains","negate":true},{"case_sensitive":false,"criteria":[{"case_sensitive":true,"match_pattern":
  ["HTTP"],"match_target_name":"","match_target_type":"filename","match_type":"contains","negate":false}],"logic":"and",
  "negate":false}],"logic":"and","negate":false}]`. The fields and descriptions included are as follows:
  + `match_target_type`: Specifies the match target type. Valid values are:
      - **schema**: The protocol type requested to be used.
      - **method**: The request method used by the request.
      - **path**: The request URL path.
      - **arg**: The query parameters in request URL.
      - **extension**: The file suffix of requested content.
      - **filename**: The file name of the requested content.
      - **header**: The HTTP request headers.
      - **clientip**: The requested client IP.
      - **clientip_version**: The requested client IP version.
      - **ua**: The User-Agent in the request header.
      - **ngx_variable**: The nginx variables.
  + `match_target_name`: Specifies the match target name.
      - When the matching target type is **schema**, **method**, **path**, **extension**, **filename**, or **ua**,
        configure this field to be empty.
      - When the matching target type is **arg**, this field represents the query parameter name.
        The length is limit from `1` to `100`. It consists of numbers, uppercase and lowercase letters, underscores and
        underlines. It can only start with a letter.
      - When the matching target type is **header**, this field indicates the name of the request header.
        The length is limit from `1` to `100`. It consists of numbers, uppercase and lowercase letters, underscores and
        underlines. It can only start with a letter.
      - When the matching target type is **clientip**, it indicates the IP source.
      - When the matching target type is **clientip_version**, it indicates the IP version source.
      - When the matching target type is **ngx_variable**, it represents the nginx variable name.
        The length is limit from `1` to `100`. It consists of numbers, uppercase and lowercase letters, underscores and
        underlines. It can only start with a letter.
  + `match_type`: The matching algorithm. Valid values are:
      - **contains**: Contains matches. When there are multiple strings in `match_pattern`, matching one of them means
        the match is successful.
      - **regex**: Regular match. Regularly matches the string in `match_pattern`. When `match_type` is regex,
        only one `match_pattern` can be filled in. This function requires the whitelist to be turned on.
  + `match_pattern`: The match content.
      - When the matching target type is **schema**, the value can be **HTTP** or **HTTPS**.
      - When the matching target type is **method**, the value can be **GET**, **PUT**, **POST**, **DELETE**, **HEAD**,
        **OPTIONS**, **PATCH**, **TRACE**, or **CONNECT**.
      - When the matching target type is **clientip_version**, the values are **IPv4** and **IPv6**.
      - When the matching target type is **path** or **ua**, this field supports wildcard *****.
  + `negate`: Whether to negate.
  + `case_sensitive`: Whether case-sensitive.
  + `logic`: Nested conditional logical operators. Valid values are **and** with **or**.
  + `criteria`: Nested condition list. The structure of this field is the same as the outermost field `criteria`.

<a name="rules_actions"></a>
The `actions` block supports:

* `http_response_header` - (Optional, List) Specifies the http response header configurations.
  The [http_response_header](#actions_http_response_header) structure is documented below.

* `access_control` - (Optional, List) Specifies the access control configurations. The maximum length of this field is `1`.
  The [access_control](#actions_access_control) structure is documented below.

* `request_url_rewrite` - (Optional, List) Specifies the request url rewrite configurations. The maximum length of this
  field is `1`.
  The [request_url_rewrite](#actions_request_url_rewrite) structure is documented below.

* `cache_rule` - (Optional, List) Specifies the cache rule configurations. The maximum length of this field is `1`.
  The [cache_rule](#actions_cache_rule) structure is documented below.

* `origin_request_url_rewrite` - (Optional, List) Specifies the origin request url rewrite configurations. The maximum
  length of this field is `1`.
  The [origin_request_url_rewrite](#actions_origin_request_url_rewrite) structure is documented below.

* `flexible_origin` - (Optional, List) Specifies the flexible origin configurations.
  The [flexible_origin](#actions_flexible_origin) structure is documented below.

* `origin_request_header` - (Optional, List) Specifies the origin request header configurations.
  The [origin_request_header](#actions_origin_request_header) structure is documented below.

<a name="actions_http_response_header"></a>
The `http_response_header` block supports:

* `action` - (Required, String) Specifies the operation type of setting HTTP response header. Valid values are **set**
  and **delete**.

* `name` - (Required, String) Specifies the HTTP response header parameter.

* `value` - (Optional, String) Specifies the value of HTTP response header parameters.

<a name="actions_access_control"></a>
The `access_control` block supports:

* `type` - (Required, String) Specifies the access control type. Valid values are **block** and **trust**.

<a name="actions_request_url_rewrite"></a>
The `request_url_rewrite` block supports:

* `redirect_url` - (Required, String) Specifies the redirect URL. The redirected URL starts with a forward slash (`/`) and
  does not contain the `http://` header and domain name, for example: /test/index.html

* `execution_mode` - (Required, String) Specifies the execution mode. Valid values are
  + **redirect**: If the requested URL matches the current rule, the request will be redirected to the target path.
    After the current rule is executed, the remaining rules will continue to be matched when other configured rules exist.
  + **break**:  If the requested URL matches the current rule, the request will be rewritten to the target path.
    After the current rule is executed, the remaining rules will no longer be matched when other configuration rules exist.
    At this time, configuring the redirection Host and redirection status code is not supported, and status code `200`
    will be returned.

* `redirect_status_code` - (Optional, Int) Specifies the redirect status code. Valid values are `301`, `302`, `303`,
  and `307`.

* `redirect_host` - (Optional, String) Specifies the redirect host. Defaults to the current domain name.
  The host supported character length limit from `1` to `255`. This field must start with `http://` or `https://`,
  for example: `http://www.example.com`.

<a name="actions_cache_rule"></a>
The `cache_rule` block supports:

* `ttl` - (Required, Int) Specifies the cache expiration time. The expiration time supports up to 365 days.

* `ttl_unit` - (Required, String) Specifies the cache expiration time unit. Valid values: **s**, **m**, **h**, and **d**.

* `follow_origin` - (Required, String) Specifies the cache expiration time source. Valid values:
  + **off**: The cache expiration time of the CDN node follows the cache expiration time in the cache rules.
  + **on**: The cache expiration time of the CDN node follows the settings of the origin site.
  + **min_ttl**: The cache expiration time of the CDN node is the minimum value between the cache rule and the origin site.

* `force_cache` - (Optional, String) Specifies whether to enable forced caching. Valid values are **on** and **off**.
  Defaults to **off**.

<a name="actions_origin_request_url_rewrite"></a>
The `origin_request_url_rewrite` block supports:

* `rewrite_type` - (Required, String) Specifies the rewrite type. Valid values:
  + **simple**: Write all URLs to the target back to the origin.
  + **wildcard**: Capture rewriting, use wildcards to capture the content in the URL to be rewritten back to the source,
    and fill it into the target return to source URL to rewrite the back to the source URL.

* `target_url` - (Required, String) Specifies the target URL. The URL starts with a forward slash (`/`), does not
  contain the `http(s)://` header and domain name, and the length does not exceed `256` characters.
  The wildcard character ***** can be captured by $n (n=`1`,`2`,`3`..., for example: `/newtest/$1/$2.jpg`)

* `source_url` - (Optional, String) Specifies the URL to be rewritten back to the source.
  When the matching mode is all files, this parameter is not supported. It must be passed when the rewriting method
  is capture rewriting, and wildcard ***** matching is supported.

<a name="actions_flexible_origin"></a>
The `flexible_origin` block supports:

* `priority` - (Required, Int) Specifies the origin priority. Valid value ranges from `1` to `100`.

* `weight` - (Required, Int) Specifies the weight. Valid value ranges from `1` to `100`.

* `sources_type` - (Required, String) Specifies the source type. Valid values are: **ipaddr**, **domain**, and **obs_bucket**.

* `ip_or_domain` - (Required, String) Specifies the origin IP or domain name.

* `origin_protocol` - (Required, String) Specifies the origin protocol. Valid values are:
  + **http**: CDN uses HTTP protocol to return to the origin.
  + **https**: CDN uses HTTPS protocol to return to the origin.
  + **follow**: The back-to-origin protocol is consistent with the client access protocol.
    For example: the client uses the HTTP protocol to access the CDN, and the CDN will also use the HTTP protocol to
    return to the origin.

* `obs_bucket_type` - (Optional, String) Specifies the OBS bucket type. Valid values are **private** and **public**.

* `http_port` - (Optional, Int) Specifies the HTTP port. Ranges from `1` to `65,535`. Defaults to `80`.

* `https_port` - (Optional, Int) Specifies the HTTPS port. Ranges from `1` to `65,535`. Defaults to `443`.

* `host_name` - (Optional, String) Specifies the host name.
  + When the origin site type is the origin site IP or the origin site domain name, defaults to the accelerated domain
    name when the source HOST is empty.
  + When the source site type is the OBS bucket domain name, defaults to the OBS bucket domain name when the source
    HOST is empty.

<a name="actions_origin_request_header"></a>
The `origin_request_header` block supports:

* `action` - (Required, String) Specifies the back-to-origin request header setting type. Valid values are **set**
  and **delete**.

* `name` - (Required, String) Specifies the back-to-origin request header parameters.

* `value` - (Optional, String) Specifies the value of the return-to-origin request header parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The acceleration domain name ID.

* `rules/rule_id` - The domain rule ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The CDN domain rule resource can be imported using the domain `name`, e.g.

```bash
$ terraform import huaweicloud_cdn_domain_rule.test <name>
```
