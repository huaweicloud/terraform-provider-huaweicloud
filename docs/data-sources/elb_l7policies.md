---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_l7policies"
description: |-
  Use this data source to get the list of ELB L7 policies.
---

# huaweicloud_elb_l7policies

Use this data source to get the list of ELB L7 policies.

## Example Usage

```hcl
variable "policy_name" {}

data "huaweicloud_elb_l7policies" "test" {
  name = var.policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source. If omitted, the provider-level
  region will be used.

* `l7policy_id` - (Optional, String) Specifies the forwarding policy ID.

* `name` - (Optional, String) Specifies the forwarding policy name.

* `description` - (Optional, String) Specifies the supplementary information about the forwarding policy.

* `listener_id` - (Optional, String) Specifies the ID of the listener to which the forwarding policy is added.

* `action` - (Optional, String) Specifies the requests are forwarded. The value can be one of the following:
  + **REDIRECT_TO_POOL**: Requests are forwarded to another backend server group.
  + **REDIRECT_TO_LISTENER**: Requests are redirected to an HTTPS listener.
  + **REDIRECT_TO_URL**: Requests are redirected to another URL.
  + **FIXED_RESPONSE**: A fixed response body is returned.

* `priority` - (Optional, Int) Specifies the forwarding policy priority.

* `provisioning_status` - (Optional, String) Specifies the provisioning status of the forwarding policy.

* `redirect_listener_id` - (Optional, String) Specifies the ID of the listener to which requests are redirected.

* `redirect_pool_id` - (Optional, String) Specifies the ID of the backend server group to which requests will be forwarded.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `l7policies` - Lists the L7 policies.
  The [l7policies](#Elb_l7policies) structure is documented below.

<a name="Elb_l7policies"></a>
The `l7policies` block supports:

* `id` - The forwarding policy ID.

* `name` - The forwarding policy name.

* `description` - The supplementary information about the forwarding policy.

* `enterprise_project_id` - The enterprise project ID.

* `action` - The requests will be forwarded.

* `listener_id` - The ID of the listener to which the forwarding policy is added.

* `priority` - The forwarding policy priority.

* `provisioning_status` - The provisioning status of the forwarding policy.

* `redirect_pool_id` - The ID of the backend server group that requests will be forwarded to.

* `redirect_pools_config` - The list of the backend server groups to which traffic is forwarded.
  The [redirect_pools_config](#Elb_l7policies_redirect_pools_config) structure is documented below.

* `redirect_listener_id` - The ID of the listener to which requests are redirected.

* `rules` - The forwarding rules in the forwarding policy.
  The [rules](#Elb_l7policies_rules) structure is documented below.

* `redirect_url_config` - The URL to which requests are forwarded.
  The [redirect_url_config](#Elb_l7policies_redirect_url_config) structure is documented below.

* `redirect_pools_extend_config` - The backend server group that the requests are forwarded to.
  The [redirect_pools_extend_config](#Elb_l7policies_redirect_pools_extend_config) structure is documented below.

* `redirect_pools_sticky_session_config` - The session persistence between backend server groups which associated with
  the policy.
  The [redirect_pools_sticky_session_config](#Elb_l7policies_redirect_pools_sticky_session_config) structure is documented
  below.

* `fixed_response_config` - The configuration of the page that will be returned.
  The [fixed_response_config](#Elb_l7policies_fixed_response_config) structure is documented below.

* `created_at` - The time when the forwarding policy was created.

* `updated_at` - The time when the forwarding policy was updated.

<a name="Elb_l7policies_redirect_pools_config"></a>
The `redirect_pools_config` block supports:

* `pool_id` - The ID of the backend server group.

* `weight` - The weight of the backend server group.

<a name="Elb_l7policies_rules"></a>
The `rules` block supports:

* `id` - The forwarding rule ID.

<a name="Elb_l7policies_redirect_url_config"></a>
The `redirect_url_config` block supports:

* `protocol` - The protocol for redirection.The value can be **HTTP**, **HTTPS**, or **${protocol}**.

* `host` - The host name that requests are redirected to.

* `port` - The port that requests are redirected to.

* `path` - The path that requests are redirected to.

* `query` - The query string set in the URL for redirection.

* `status_code` - The status code returned after the requests are redirected.

* `insert_headers_config` - The header parameters to be added.
  The [insert_headers_config](#insert_headers_config_object) structure is documented below.

* `remove_headers_config` - The header parameters to be removed.
  The [remove_headers_config](#remove_headers_config_object) structure is documented below.

<a name="Elb_l7policies_redirect_pools_extend_config"></a>
The `redirect_pools_extend_config` block supports:

* `rewrite_url_enabled` - Whether to enable URL redirection.

* `rewrite_url_config` - The URL for the backend server group that requests are forwarded to.
  The [rewrite_url_config](#Elb_l7policies_rewrite_url_config) structure is documented below.

* `insert_headers_config` - The header parameters to be added.
  The [insert_headers_config](#insert_headers_config_object) structure is documented below.

* `remove_headers_config` - The header parameters to be removed.
  The [remove_headers_config](#remove_headers_config_object) structure is documented below.

* `traffic_limit_config` - The traffic limit config of the policy.
  The [traffic_limit_config](#traffic_limit_config_object) structure is documented below.

<a name="Elb_l7policies_redirect_pools_sticky_session_config"></a>
The `redirect_pools_sticky_session_config` block supports:

* `enable` - Whether enable config session persistence between backend server groups.

* `timeout` - The timeout of the session persistence.

<a name="Elb_l7policies_fixed_response_config"></a>
The `fixed_response_config` block supports:

* `status_code` -The HTTP status code configured in the forwarding policy.

* `content_type` - The format of the response body.

* `message_body` - The content of the response message body.

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

* `key` - The parameter name of the added request header.

* `value_type` - The value type of the parameter.

* `value` - The value of the parameter.

<a name="remove_headers_config_object"></a>
The `remove_headers_config` block supports:

* `configs` - The list of request header parameters to be removed.
  The [remove_header_configs](#remove_header_configs_object) structure is documented below.

<a name="remove_header_configs_object"></a>
The `remove_header_configs` block supports:

* `key` - The parameter name of the removed request header.

<a name="traffic_limit_config_object"></a>
The `traffic_limit_config` block supports:

* `qps` - The overall qps of the policy.

* `per_source_ip_qps` - The single source qps of the policy.

* `burst` - (The qps buffer.

<a name="Elb_l7policies_rewrite_url_config"></a>
The `rewrite_url_config` block supports:

* `host` - The url host.

* `path` - The URL path.

* `query` - The URL query character string.
