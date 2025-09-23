---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_l7policy"
description: ""
---

# huaweicloud_elb_l7policy

Manages an ELB L7 Policy resource within HuaweiCloud.

## Example Usage

### ELB L7 Policy redirect to pool

```hcl
variable listener_id {}
variable pool_id {}

resource "huaweicloud_elb_l7policy" "policy_1" {
  name             = "policy_1"
  action           = "REDIRECT_TO_POOL"
  priority         = 20
  description      = "test description"
  listener_id      = var.listener_id
  redirect_pool_id = var.pool_id

  redirect_pools_extend_config {
    rewrite_url_enabled = true
    
    rewrite_url_config {
      host  = "test.com"
      path  = "/path"
      query = "abc"
    }
  }
}
```

### ELB L7 Policy redirect to listener

```hcl
variable listener_id {}
variable redirect_listener_id {}

resource "huaweicloud_elb_l7policy" "policy_1" {
  name                 = "policy_1"
  action               = "REDIRECT_TO_LISTENER"
  description          = "test description"
  listener_id          = var.listener_id
  redirect_listener_id = var.redirect_listener_id
}
```

### ELB L7 Policy redirect to URL

```hcl
variable listener_id {}

resource "huaweicloud_elb_l7policy" "policy_1" {
  name        = "policy_1"
  action      = "REDIRECT_TO_URL"
  priority    = 20
  description = "test description"
  listener_id = var.listener_id
  
  redirect_url_config {
    protocol    = "HTTP"
    host        = "test.com"
    port        = "6666"
    path        = "/test_policy"
    query       = "test_query"
    status_code = "301"
  }
}
```

### ELB L7 Policy redirect to fixed response

```hcl
variable listener_id {}

resource "huaweicloud_elb_l7policy" "policy_1" {
  name        = "policy_1"
  action      = "FIXED_RESPONSE"
  priority    = 20
  description = "test description"
  listener_id = var.listener_id
  
  fixed_response_config {
    status_code  = "200"
    content_type = "application/json"
    message_body = "it is a test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the L7 Policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new L7 Policy.

* `name` - (Optional, String) Human-readable name for the L7 Policy. Does not have to be unique.

* `priority` - (Optional, Int) The forwarding policy priority. A smaller value indicates a higher priority. The value
  must be unique for forwarding policies of the same listener. This parameter will take effect only when
  `advanced_forwarding_enabled` of the listener is set to **true**. If `action` is set to **REDIRECT_TO_LISTENER**,
  the value can only be 0.

* `description` - (Optional, String) Human-readable description for the L7 Policy.

* `listener_id` - (Required, String, ForceNew) The Listener on which the L7 Policy will be associated with. Changing
  this creates a new L7 Policy.

* `action` - (Optional, String, ForceNew) Whether requests are forwarded to another backend server group
  or redirected to an HTTPS listener. Changing this creates a new L7 Policy. The value ranges:
  + **REDIRECT_TO_POOL**: Requests are forwarded to the backend server group specified by `redirect_pool_id` or
    `redirect_pools_config`, the `protocol` of the listener must be **HTTP** or **HTTPS**.
  + **REDIRECT_TO_LISTENER**: Requests are redirected from the HTTP listener specified by `listener_id` to the
    HTTPS listener specified by `redirect_listener_id`, the `protocol` of the listener must be **HTTP**.
  + **REDIRECT_TO_URL**: Requests are forwarded to another URL whose config specified by `redirect_url_config`.
  + **FIXED_RESPONSE**: Requests are forwarded to a fixed response body specified by `fixed_response_config`.
  Defaults to **REDIRECT_TO_POOL**.

* `redirect_pool_id` - (Optional, String) The ID of the backend server group to which traffic is forwarded.
  This parameter will take effect when `action` is set to **REDIRECT_TO_POOL**. The backend server group must meet the
  following requirements:
  + Cannot be the default backend server group of the listener.
  + Cannot be the backend server group used by forwarding policies of other listeners.

* `redirect_pools_config` - (Optional, List) The list of the backend server groups to which traffic is forwarded.
  traffic is redirected. This parameter will take effect when `action` is set to **REDIRECT_TO_POOL**.
  The [redirect_pools_config](#redirect_pools_config_object) structure is documented below.

  -> **NOTE:** Exactly one of `redirect_pool_id` or `redirect_pools_config` should be specified when `action` is set to
  **REDIRECT_TO_POOL**.

* `redirect_pools_sticky_session_config` - (Optional, List) The session persistence between backend server groups which
  associated with the policy. This parameter will take effect when `action` is set to **REDIRECT_TO_POOL**.
  The [redirect_pools_sticky_session_config](#redirect_pools_sticky_session_config_object) structure is documented below.

* `redirect_pools_extend_config` - (Optional, List) The config of the backend server group to which the
  traffic is redirected. This parameter will take effect when `action` is set to **REDIRECT_TO_POOL**.
  The [redirect_pools_extend_config](#redirect_pools_extend_config_object) structure is documented below.

* `redirect_listener_id` - (Optional, String) The ID of the listener to which the traffic is redirected.
  This parameter is mandatory when `action` is set to **REDIRECT_TO_LISTENER**. The listener must meet the
  following requirements:
  + Can only be an HTTPS listener.
  + Can only be a listener of the same load balancer.

* `redirect_url_config` - (Optional, List) The URL config to which the traffic is redirected.
  This parameter is mandatory when `action` is set to **REDIRECT_TO_URL**. The `advanced_forwarding_enabled` of the
  listener must be set to **true**.
  The [redirect_url_config](#redirect_url_config_object) structure is documented below.

* `fixed_response_config` - (Optional, List) The fixed configuration of the page to which the traffic is
  redirected. This parameter is mandatory when `action` is set to **FIXED_RESPONSE**. The `advanced_forwarding_enabled` of
  the listener must be set to **true**.
  The [fixed_response_config](#fixed_response_config_object) structure is documented below.

<a name="redirect_pools_config_object"></a>
The `redirect_pools_config` block supports:

* `pool_id` - (Required, String) The ID of the backend server group.

* `weight` - (Optional, Int) The weight of the backend server group.

<a name="redirect_pools_sticky_session_config_object"></a>
The `redirect_pools_sticky_session_config` block supports:

* `enable` - (Optional, Bool) Whether enable config session persistence between backend server groups.

* `timeout` - (Optional, Int) The timeout of the session persistence.

<a name="redirect_pools_extend_config_object"></a>
The `redirect_pools_extend_config` block supports:

* `rewrite_url_enabled` - (Optional, Bool) Whether the rewrite url is enabled.

* `rewrite_url_config` - (Optional, List) The rewrite url config. This parameter is mandatory when `rewrite_url_enabled`
  is set to **true**.
  The [rewrite_url_config](#rewrite_url_config_object) structure is documented below.

* `insert_headers_config` - (Optional, List) The header parameters to be added.
  The [insert_headers_config](#insert_headers_config_object) structure is documented below.

* `remove_headers_config` - (Optional, List) The header parameters to be removed.
  The [remove_headers_config](#remove_headers_config_object) structure is documented below.

* `traffic_limit_config` - (Optional, List) The traffic limit config of the policy.
  The [traffic_limit_config](#traffic_limit_config_object) structure is documented below.

<a name="rewrite_url_config_object"></a>
The `rewrite_url_config` block supports:

* `host` - (Optional, String) The host name of the rewrite URL. The value can contain only letters, digits, and
  periods (.) and must start with a letter or digit. Defaults to **${host}**, indicates inherit original value.

* `path` - (Optional, String) The path of the rewrite URL. The value can contain only letters, digits, and special
  characters _~';@^-%#&$.+?,=!:|\/(), and must start with a slash (/). Defaults to **${path}**, indicates inherit
  original value.

* `query` - (Optional, String) The query of the rewrite URL. The value can contain only letters, digits, and special
  characters !$&'()+,-./:;=?@^_`. Defaults to **${query}**, indicates inherit original value.

<a name="redirect_url_config_object"></a>
The `redirect_url_config` block supports:

* `status_code` - (Required, String) The status code returned after the requests are redirected.
  Value options: **301**, **302**, **303**, **307**, **308**.

* `protocol` - (Optional, String) The protocol for redirection. Value options: **HTTP**, **HTTPS**, **${protocol}**.
  Defaults to **${protocol}**, indicating that the path of the request will be used.

* `host` - (Optional, String) The host name that requests are redirected to. The value can contain only letters,
  digits, hyphens (-), and periods (.) and must start with a letter or digit. Defaults to **${host}**, indicating
  that the host of the request will be used.

* `port` - (Optional, String) The  port that requests are redirected to. Defaults to **${port}**, indicating that
  the port of the request will be used.

* `path` - (Optional, String) The path that requests are redirected to. The value can contain only letters, digits,
  and special characters _~';@^- %#&$.*+?,=!:|/()[]{} and must start with a slash (/).
  Defaults to **${path}**, indicating that the path of the request will be used.

* `query` - (Optional, String) The query string set in the URL for redirection. The value is case-sensitive and can
  contain only letters, digits, and special characters !$&'()*+,-./:;=?@^_\`. Defaults to **${query}**, indicating that
  the query string of the request will be used.
  For example, in the URL `https://www.xxx.com:8080/elb?type=loadbalancer`, **${query}** indicates **type=loadbalancer**.
  If this parameter is set to **${query}&name=my_name**, the URL will be redirected to
  URL `https://www.xxx.com:8080/elb?type=loadbalancer&name=my_name`.

* `insert_headers_config` - (Optional, List) The header parameters to be added.
  The [insert_headers_config](#insert_headers_config_object) structure is documented below.

* `remove_headers_config` - (Optional, List) The header parameters to be removed.
  The [remove_headers_config](#remove_headers_config_object) structure is documented below.

<a name="fixed_response_config_object"></a>
The `fixed_response_config` block supports:

* `status_code` - (Required, String) The fixed HTTP status code configured in the forwarding rule. The value can be
  any integer in the range of **200–299**, **400–499**, or **500–599**.

* `content_type` - (Optional, String) The format of the response body. Value options: **text/plain**, **text/css**,
  **text/html**, **application/javascript**, or **application/json**. Defaults to: **text/plain**.

* `message_body` - (Optional, String) The content of the response message body.

* `insert_headers_config` - (Optional, List) The header parameters to be added.
  The [insert_headers_config](#insert_headers_config_object) structure is documented below.

* `remove_headers_config` - (Optional, List) The header parameters to be removed.
  The [remove_headers_config](#remove_headers_config_object) structure is documented below.

* `traffic_limit_config` - (Optional, List) The traffic limit config of the policy.
  The [traffic_limit_config](#traffic_limit_config_object) structure is documented below.

<a name="insert_headers_config_object"></a>
The `insert_headers_config` block supports:

* `configs` - (Required, List) The list of request header parameters to be added.
  The [insert_header_configs](#insert_header_configs_object) structure is documented below.

<a name="insert_header_configs_object"></a>
The `insert_header_configs` block supports:

* `key` - (Required, String) The parameter name of the added request header. The value can contain `1` to `40`
  characters, only a-z, digits, hyphens (-) and underscore (_) are allowed, and it can not be the following characters:
  **connection**, **upgrade**, **content-length**, **transfer-encoding**, **keep-alive**, **te**, **host**, **cookie**,
  **remoteip**, **authority**, **x-forwarded-host**, **x-forwarded-for**, **x-forwarded-for-port**,
  **x-forwarded-tls-certificate-id**, **x-forwarded-tls-protocol**, **x-forwarded-tls-cipher**, **x-forwarded-elb-ip**,
  **x-forwarded-port**, **x-forwarded-elb-id**, **x-forwarded-elb-vip**, **x-real-ip**, **x-forwarded-proto**,
  **x-nuwa-trace-ne-in**, **x-nuwa-trace-ne-out**.

* `value_type` - (Required, String) The value type of the parameter. Value options: **USER_DEFINED**,
  **REFERENCE_HEADER**, **SYSTEM_DEFINED**.

* `value` - (Required, String) The value of the parameter. The value can contain `1` to `128`, only printable
  characters in the range of ASCII code value 32<=ch<=127, asterisks (*) and question marks (?) are allowed, and it
  cannot start or end with a space characters. If the value of `value_type` is **SYSTEM_DEFINED**, the value options is:
  **CLIENT-PORT**, **CLIENT-IP**, **ELB-PROTOCOL**, **ELB-ID**, **ELB-PORT**, **ELB-EIP**, **ELB-VIP**.

<a name="remove_headers_config_object"></a>
The `remove_headers_config` block supports:

* `configs` - (Required, List) The list of request header parameters to be removed.
  The [remove_header_configs](#remove_header_configs_object) structure is documented below.

<a name="remove_header_configs_object"></a>
The `remove_header_configs` block supports:

* `key` - (Required, String) The parameter name of the removed request header. The value can contain `1` to `40`
  characters, only a-z, digits, hyphens (-) and underscore (_) are allowed, and it can not be the following characters:
  **connection**, **upgrade**, **content-length**, **transfer-encoding**, **keep-alive**, **te**, **host**, **cookie**,
  **remoteip**, **authority**, **x-forwarded-host**, **x-forwarded-for**, **x-forwarded-for-port**,
  **x-forwarded-tls-certificate-id**, **x-forwarded-tls-protocol**, **x-forwarded-tls-cipher**, **x-forwarded-elb-ip**,
  **x-forwarded-port**, **x-forwarded-elb-id**, **x-forwarded-elb-vip**, **x-real-ip**, **x-forwarded-proto**,
  **x-nuwa-trace-ne-in**, **x-nuwa-trace-ne-out**.

<a name="traffic_limit_config_object"></a>
The `traffic_limit_config` block supports:

* `qps` - (Optional, Int) The overall qps of the policy.  
  The valid value is range form `0` to `100,000`, `0` indicates no limit.

* `per_source_ip_qps` - (Optional, Int) The single source qps of the policy.  
  The valid value is range form `0` to `100,000`, `0` indicates no limit.  
  If the value of `qps` is not `0`, then the value of `per_source_ip_qps` must less than the value of `qps`.
  If the `protocol` of the listener that the policy associated with is **QUIC**, then `per_source_ip_qps` is not
  supported, the value should be `0` or empty.

* `burst` - (Optional, Int) The qps buffer.  
  The valid value is range form `0` to `100,000`. When qps exceeds the limit, 503 will not be
  returned, and requests that allow local burst size increases are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the L7 policy.

* `provisioning_status` - The provisioning status of the forwarding policy.

* `enterprise_project_id` - The ID of the enterprise project.

* `created_at` - The creation time of the L7 policy.

* `updated_at` - The update time of the L7 policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_policy.policy_1 <id>
```
