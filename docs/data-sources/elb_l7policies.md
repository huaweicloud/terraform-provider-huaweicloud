---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_l7policies"
description: ""
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

* `action` - The requests will be forwarded. The value can be one of the following:
  **REDIRECT_TO_POOL**, **REDIRECT_TO_LISTENER**, **REDIRECT_TO_URL**, **FIXED_RESPONSE**.

* `listener_id` - The ID of the listener to which the forwarding policy is added.

* `priority` - The forwarding policy priority.

* `redirect_pool_id` - The ID of the backend server group that requests will be forwarded to.

* `redirect_listener_id` - The ID of the listener to which requests are redirected.

* `rules` - The forwarding rules in the forwarding policy. The [rules](#Elb_l7policies_rules) structure is documented below.

* `redirect_url_config` - The URL to which requests are forwarded. The [redirect_url_config](#Elb_l7policies_redirect_url_config)
  structure is documented below.

* `redirect_pools_extend_config` - The backend server group that the requests are forwarded to.
  The [redirect_pools_extend_config](#Elb_l7policies_redirect_pools_extend_config) structure is documented below.

* `fixed_response_config` - The configuration of the page that will be returned.
  The [fixed_response_config](#Elb_l7policies_fixed_response_config) structure is documented below.

* `created_at` - The time when the forwarding policy was created.

* `updated_at` - The time when the forwarding policy was updated.

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

<a name="Elb_l7policies_redirect_pools_extend_config"></a>
The `redirect_pools_extend_config` block supports:

* `rewrite_url_enabled` - Whether to enable URL redirection.

* `rewrite_url_config` - The URL for the backend server group that requests are forwarded to.
  The [rewrite_url_config](#Elb_l7policies_rewrite_url_config) structure is documented below.

<a name="Elb_l7policies_fixed_response_config"></a>
The `fixed_response_config` block supports:

* `status_code` -The HTTP status code configured in the forwarding policy.

* `content_type` - The format of the response body.

* `message_body` - The content of the response message body.

<a name="Elb_l7policies_rewrite_url_config"></a>
The `rewrite_url_config` block supports:

* `host` - The url host.

* `path` - The URL path.

* `query` - The URL query character string.
