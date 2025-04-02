---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_active_standby_pools"
description: ""
---

# huaweicloud_elb_active_standby_pools

Use this data source to get the list of active standby ELB pools.

## Example Usage

```hcl
variable "pool_name" {}

data "huaweicloud_elb_active_standby_pools" "test" {
  name = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the active-standby pool.

* `pool_id` - (Optional, String) Specifies the ID of the active-standby pool.

* `description` - (Optional, String) Specifies supplementary information about the active-standby pool.

* `connection_drain` - (Optional, String) Specifies whether delayed logout is enabled. Value options:
  + **false**: Disable this option.
  + **true**: Enable this option.

* `ip_version` - (Optional, String) Specifies the IP address version supported by the pool.

* `lb_algorithm` - (Optional, String) Specifies the load balancing algorithm used by the load balancer to route requests
  to backend servers in the associated pool. Value options:
  + **ROUND_ROBIN**: weighted round robin.
  + **LEAST_CONNECTIONS**: weighted least connections.
  + **SOURCE_IP**: source IP hash.
  + **QUIC_CID**: connection ID.

* `loadbalancer_id` - (Optional, String) Specifies the ID of the load balancer with which the active-standby pool is
  associated.

* `healthmonitor_id` - (Optional, String) Specifies the ID of the health check configured for the active-standby pool.

* `protocol` - (Optional, String) Specifies the protocol used by the active-standby pool to receive requests from the
  load balancer. Value options: **TCP**, **UDP**, **QUIC** or **TLS**.

* `member_address` - (Optional, String) Specifies the private IP address bound to the member. This parameter is used
  only as a query condition and is not included in the response.

* `member_instance_id` - (Optional, String) Specifies the ID of the ECS used as the member. This parameter is used only
  as a query condition and is not included in the response.

* `listener_id` - (Optional, String) Specifies the ID of the listener to which the forwarding policy is added.

* `type` - (Optional, String) Specifies the type of the active-standby pool.
  The valid values are as follows:
  + **instance**: Any type of backend servers can be added.
  + **ip**: Only IP as backend servers can be added.

  If the `type` is empty, any type of backend servers can be added.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC where the active-standby pool works.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pools` - The pool list. For details, see data structure of the pool field.
  The [pools](#elb_pools) structure is documented below.

<a name="elb_pools"></a>
The `pools` block supports:

* `id` - The ID of the active-standby pool.

* `name` - The name of the active-standby pool.

* `description` - The description of the active-standby pool.

* `protocol` - The protocol used by the active-standby pool to receive requests.

* `type` - The type of the active-standby pool.

* `any_port_enable` - Whether to enable Forward to same Port for a pool.

* `enterprise_project_id` - The ID of the enterprise project.

* `ip_version` - The IP address version supported by the pool.

* `lb_algorithm` - The load balancing algorithm used by the load balancer to route requests to backend servers in the
  associated pool.

* `vpc_id` - The ID of the VPC where the active-standby pool works.

* `connection_drain_enabled` - Whether to enable delayed logout.

* `connection_drain_timeout` - The timeout of the delayed logout in seconds.

* `listeners` - The IDs of the listeners with which the active-standby pool is associated.
  The [listeners](#elb_listeners) structure is documented below.

* `loadbalancers` - The IDs of the load balancers with which the active-standby pool is associated.
  The [loadbalancers](#elb_loadbalancers) structure is documented below.

* `members` - The backend servers in the active-standby pool.
  The [members](#elb_members) structure is documented below.

* `healthmonitor` - The health check configured for the active-standby pool.
  The [healthmonitor](#elb_healthmonitor) structure is documented below.

* `quic_cid_hash_strategy` - The multi-path distribution configuration based on destination connection IDs.
  The [quic_cid_hash_strategy](#elb_quic_cid_hash_strategy) structure is documented below.

* `created_at` - The time when the backend server group was created.

* `updated_at` - The time when the backend server group was updated.

<a name="elb_listeners"></a>
The `listeners` block supports:

* `id` - The listener ID.

<a name="elb_loadbalancers"></a>
The `loadbalancers` block supports:

* `id` - The loadbalancer ID.

<a name="elb_members"></a>
The `members` block supports:

* `id` - The ID of the member.

* `name` - The name of the member.

* `subnet_id` - The ID of the IPv4 or IPv6 subnet where the member resides.

* `protocol_port` - The port used by the member to receive requests.

* `address` - The private IP address bound to the member.

* `ip_version` - The IP version supported by the member.

* `operating_status` - The health status of the member.

* `member_type` - The type of the member.

* `instance_id` - The ID of the ECS used as the member.

* `role` - The active-standby status of the member.

* `reason` - Why health check fails.
  The [reason](#ELB_member_reason) structure is documented below.

* `status` - The health status of the backend server if `listener_id` under status is specified. If `listener_id` under
  status is not specified, operating_status of member takes precedence.
  The [status](#ELB_member_status) structure is documented below.

<a name="ELB_member_reason"></a>
The `reason` block supports:

* `reason_code` - The code of the health check failures. The value can be:
  + **CONNECT_TIMEOUT**: The connection with the backend server times out during a health check.
  + **CONNECT_REFUSED**: The load balancer rejects connections with the backend server during a health check.
  + **CONNECT_FAILED**: The load balancer fails to establish connections with the backend server during a health check.
  + **CONNECT_INTERRUPT**: The load balancer is disconnected from the backend server during a health check.
  + **SSL_HANDSHAKE_ERROR**: The SSL handshakes with the backend server fail during a health check.
  + **RECV_RESPONSE_FAILED**: The load balancer fails to receive responses from the backend server during a health check.
  + **RECV_RESPONSE_TIMEOUT**: The load balancer does not receive responses from the backend server within the timeout
    duration during a health check.
  + **SEND_REQUEST_FAILED**: The load balancer fails to send a health check request to the backend server during a health
    check.
  + **SEND_REQUEST_TIMEOUT**: The load balancer fails to send a health check request to the backend server within the
    timeout duration.
  + **RESPONSE_FORMAT_ERROR**: The load balancer receives invalid responses from the backend server during a health check.
  + **RESPONSE_MISMATCH**: The response code received from the backend server is different from the preset code.

* `expected_response` - The expected HTTP status code. This parameter will take effect only when `type` is set to **HTTP**,
  **HTTPS** or **GRPC**.
  + A specific status code. If `type` is set to **GRPC**, the status code ranges from **0** to **99**. If `type` is set
    to other values, the status code ranges from **200** to **599**.
  + A list of status codes that are separated with commas (,). A maximum of five status codes are supported.
  + A status code range. Different ranges are separated with commas (,). A maximum of five ranges are supported.

* `healthcheck_response` - The returned HTTP status code in the response. This parameter will take effect only when `type`
  is set to **HTTP**, **HTTPS** or **GRPC**.
  + A specific status code. If type is set to **GRPC**, the status code ranges from **0** to **99**. If `type` is set to
    other values, the status code ranges from **200** to **599**.

<a name="ELB_member_status"></a>
The `status` block supports:

* `listener_id` - The ID of the listener associated with the backend server.

* `operating_status` - The health status of the backend server. The value can be:
  + **ONLINE**: The backend server is running normally.
  + **NO_MONITOR**: No health check is configured for the backend server group to which the backend server belongs.
  + **OFFLINE**: The cloud server used as the backend server is stopped or does not exist.

<a name="elb_healthmonitor"></a>
The `healthmonitor` block supports:

* `id` - The health check ID.

* `name` - The health check name.

* `delay` - The interval between health checks, in seconds.

* `domain_name` - The domain name that HTTP requests are sent to during the health check.

* `expected_codes` - The expected HTTP status code.

* `http_method` - The HTTP method.

* `max_retries_down` - The number of consecutive health checks when the health check result of a backend server changes
  from **ONLINE** to **OFFLINE**.

* `max_retries` - The number of consecutive health checks when the health check result of a backend server changes from
  **OFFLINE** to **ONLINE**.

* `monitor_port` - The port used for the health check.

* `timeout` - The maximum time required for waiting for a response from the health check, in seconds.

* `type` - The health check protocol.

* `url_path` - The HTTP request path for the health check.

<a name="elb_quic_cid_hash_strategy"></a>
The `quic_cid_hash_strategy` block supports:

* `len` - The length of the hash factor in the connection ID, in byte.

* `offset` - The start position in the connection ID as the hash factor, in byte.
