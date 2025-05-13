---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_plugin"
description: ""
---

# huaweicloud_apig_plugin

Manages a plugin resource within HuaweiCloud.

## Example Usage

### Create a CORS plugin

```hcl
variable "instance_id" {}
variable "plugin_name" {}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = var.instance_id
  name        = var.plugin_name
  type        = "cors"
  content     = jsonencode(
    {
      allow_origin      = "*"
      allow_methods     = "GET,PUT,DELETE,HEAD,PATCH"
      allow_headers     = "Content-Type,Accept,Cache-Control"
      expose_headers    = "X-Request-Id,X-Apig-Latency"
      max_age           = 12700
      allow_credentials = true
    }
  )
}
```

### Create an HTTP Response Header Management plugin

```hcl
variable "instance_id" {}
variable "plugin_name" {}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = var.instance_id
  name        = var.plugin_name
  type        = "set_resp_headers"
  content     = jsonencode(
    {
      response_headers = [{
        name       = "X-Custom-Pwd"
        value      = "**********"
        value_type = "custom_value"
        action     = "override"
      },
      {
        name       = "X-Custom-Debug-Step"
        value      = "Beta"
        value_type = "custom_value"
        action     = "skip"
      },
      {
        name       = "X-Custom-Config"
        value      = "<HTTP response test>"
        value_type = "custom_value"
        action     = "append"
      },
      {
        name       = "X-Custom-Id"
        value      = ""
        value_type = "custom_value"
        action     = "delete"
      },
      {
        name       = "X-Custom-Log-Level"
        value      = "DEBUG"
        value_type = "custom_value"
        action     = "add"
      },
      {
        name       = "Sys-Param"
        value      = "$context.cacheStatus"
        value_type = "system_parameter"
        action     = "add"
      }]
    }
  )
}
```

### Create a Request Throttling 2.0 plugin

```hcl
variable "instance_id" {}
variable "plugin_name" {}
variable "application_id" {}
variable "user_id" {}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = var.instance_id
  name        = var.plugin_name
  type        = "rate_limit"
  content     = jsonencode(
    {
      "scope": "basic",
      "default_time_unit": "minute",
      "default_interval": 1,
      "api_limit": 25,
      "app_limit": 10,
      "user_limit": 15,
      "ip_limit": 25,
      "algorithm": "counter",
      "specials": [
        {
          "type": "app",
          "policies": [
            {
              "key": "${var.application_id}",
              "limit": 10
            }
          ]
        },
        {
          "type": "user",
          "policies": [
            {
              "key": "${var.user_id}",
              "limit": 10
            }
          ]
        }
      ],
      "parameters": [
        {
          "type": "path",
          "name": "reqPath",
          "value": "reqPath"
        },
        {
          "type": "method",
          "name": "method",
          "value": "method"
        },
        {
          "type": "system",
          "name": "serverName",
          "value": "serverName"
        }
      ],
      "rules": [
        {
          "rule_name": "rule-0001",
          "match_regex": "[\"AND\",[\"method\",\"~=\",\"POST\"],[\"method\",\"~=\",\"PATCH\"]]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 20
        },
        {
          "rule_name": "rule-0002",
          "match_regex": "[\"reqPath\",\"~~\",\"/terraform/test/*/\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 10
        },
        {
          "rule_name": "rule-0003",
          "match_regex": "[\"serverName\",\"==\",\"terraform\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 15
        },
        {
          "rule_name": "rule-0004",
          "match_regex": "[\"method\",\"in\",\"PATCH\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 5
        }
      ]
    }
  )
}
```

### Create a Kafka Log Push plugin

```hcl
variable "instance_id" {}
variable "plugin_name" {}
variable "connect_addresses" {
  type = list(string)
}
variable "topic_name" {}
variable "connect_port" {}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = var.instance_id
  name        = var.plugin_name
  type        = "kafka_log"
  content     = jsonencode(
    {
      "broker_list": [for v in var.connect_addresses: format("%s:%d", v, var.connect_port)],
      "topic": "${var.topic_name}",
      "key": "",
      "max_retry_count": 3,
      "retry_backoff": 10,
      "sasl_config": {
        "security_protocol": "PLAINTEXT",
        "sasl_mechanisms": "PLAIN",
        "sasl_username": "",
        "sasl_password": "",
        "ssl_ca_content": ""
      },
      "meta_config": {
        "system": {
          "start_time": true,
          "request_id": true,
          "client_ip": true,
          "api_id": false,
          "user_name": false,
          "app_id": false,
          "access_model1": false,
          "request_time": true,
          "http_status": true,
          "server_protocol": false,
          "scheme": true,
          "request_method": true,
          "host": false,
          "api_uri_mode": false,
          "uri": false,
          "request_size": false,
          "response_size": false,
          "upstream_uri": false,
          "upstream_addr": true,
          "upstream_status": true,
          "upstream_connect_time": false,
          "upstream_header_time": false,
          "upstream_response_time": true,
          "all_upstream_response_time": false,
          "region_id": false,
          "auth_type": false,
          "http_x_forwarded_for": true,
          "http_user_agent": true,
          "error_type": true,
          "access_model2": false,
          "inner_time": false,
          "proxy_protocol_vni": false,
          "proxy_protocol_vpce_id": false,
          "proxy_protocol_addr": false,
          "body_bytes_sent": false,
          "api_name": false,
          "app_name": false,
          "provider_app_id": false,
          "provider_app_name": false,
          "custom_data_log01": false,
          "custom_data_log02": false,
          "custom_data_log03": false,
          "custom_data_log04": false,
          "custom_data_log05": false,
          "custom_data_log06": false,
          "custom_data_log07": false,
          "custom_data_log08": false,
          "custom_data_log09": false,
          "custom_data_log10": false,
          "response_source": false
        },
        "call_data": {
          "log_request_header": true,
          "log_request_query_string": true,
          "log_request_body": true,
          "log_response_header": true,
          "log_response_body": true,
          "request_header_filter": "X-Custom-Auth-Type",
          "request_query_string_filter": "authId",
          "response_header_filter": "X-Trace-Id",
          "custom_authorizer": {
            "frontend": [
              "user_name",
              "user_age"
            ],
            "backend": [
              "userName",
              "userAge"
            ]
          }
        }
      }
    }
  )
}
```

### Create a Circuit Breaker plugin

```hcl
variable "instance_id" {}
variable "plugin_name" {}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = var.instance_id
  name        = var.plugin_name
  type        = "breaker"
  content     = jsonencode(
    {
      "breaker_condition": {
        "breaker_type": "condition",
        "breaker_mode": "counter",
        "unhealthy_condition": "[\"OR\",[\"$context.statusCode\",\"in\",\"500,501,504\"],[\"$context.backendResponseTime\",\">\",6000]]",
        "unhealthy_threshold": 30,
        "min_call_threshold": 20,
        "unhealthy_percentage": 51,
        "time_window": 15,
        "open_breaker_time": 15
      },
      "downgrade_default": null,
      "downgrade_parameters": [
        {
          "type": "path",
          "name": "reqPath",
          "value": "reqPath"
        },
        {
          "type": "method",
          "name": "method",
          "value": "method"
        },
        {
          "type": "query",
          "name": "authType",
          "value": "authType"
        }
      ],
      "downgrade_rules": [
        {
          "breaker_condition": {
            "breaker_type": "timeout",
            "breaker_mode": "percentage",
            "unhealthy_condition": "",
            "unhealthy_threshold": 30,
            "min_call_threshold": 20,
            "unhealthy_percentage": 51,
            "time_window": 15,
            "open_breaker_time": 15
          },
          "downgrade_backend": null,
          "rule_name": "rule-qkqe",
          "match_regex": "[\"authType\",\"~=\",\"Token\"]",
          "parameters": [
            "reqPath",
            "method",
            "authType"
          ]
        }
      ],
      "scope": "basic"
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the plugin is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the plugin
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the plugin name.  
  The valid length is limited from `3` to `255`, only Chinese characters, English letters, digits, hyphens (-) and
  underscores (_) are allowed. The name must start with an English letter or Chinese character.

* `type` - (Required, String, ForceNew) Specifies the plugin type.  
  The valid values are as follows:
  + **cors**: CORS, specify preflight request headers and response headers and automatically create preflight request
    APIs for cross-origin API access.
  + **set_resp_headers**: HTTP Response Header Management, customize HTTP headers that will be contained in an API
    response.
  + **rate_limit**: Request Throttling 2.0, limits the number of times an API can be called within a specific time
    period. It supports parameter-based, basic, and excluded throttling.
  + **kafka_log**: Kafka Log Push, Push detailed API calling logs to kafka for you to easily obtain logs.
  + **breaker**: Circuit Breaker, circuit breaker protect the system when performance issues occur on backend service.
  + **third_auth**: Third-Party Authorizer.
  + **proxy_cache**: Proxy Cache.
  + **proxy_mirror**: Proxy Mirror.

  Changing this will create a new resource.

* `content` - (Required, String) Specifies the configuration details for plugin.
  + For `CORS` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0021.html).
  + For `HTTP Response Header Management` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0022.html).
  + For `Request Throttling 2.0` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0054.html).
  + For `Kafka Log Push` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0061.html).
  + For `Circuit Breaker` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0023.html).
  + For `Third-Party Authorizer` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0077.html).
  + For `Proxy Cache` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0111.html).
  + For `Proxy Mirror` plugins, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0112.html).

  -> All default values in content need to be filled in, otherwise `terraform plan` will prompt the script to change.
     You can confirm the content of a policy through the console (**APIG** -> **API Policies** -> **Your Policy**
     -> **Policy Content**).

  ~> For sensitive values filled in `content` (such as **SASL password**), the API response will not be returned, so
     `terraform plan` will prompt script changes, and you can ignore these changes through `ignorechanges` in the
     `lifecycle`.

* `description` - (Optional, String) Specifies the plugin description.  
  The valid length is limited from `3` to `255` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the plugin.

* `created_at` - The creation time of the plugin.

* `updated_at` - The latest update time of the plugin.

## Import

Plugins can be imported using their related dedicated instance ID (`instance_id`) and their ID (`id`), separated by a
slash, e.g.

```bash
$ terraform import huaweicloud_apig_plugin.test <instance_id>/<id>
```
