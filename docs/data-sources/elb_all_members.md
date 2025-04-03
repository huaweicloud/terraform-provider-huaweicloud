---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_all_members"
description: |-
  Use this data source to get the list of members under the current project.
---

# huaweicloud_elb_all_members

Use this data source to get the list of members under the current project.

## Example Usage

```hcl
data "huaweicloud_elb_all_members" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, List) Specifies the backend server name.
  Multiple IDs can be queried.

* `weight` - (Optional, List) Specifies the weight of the backend server.
  Multiple weights can be queried.

* `subnet_cidr_id` - (Optional, List) Specifies the ID of the subnet where the backend server works.
  Multiple IDs can be queried.

* `address` - (Optional, List) Specifies the IP address of the backend server.
  Multiple IP addresses can be queried.

* `protocol_port` - (Optional, List) Specifies the port used by the backend servers.
  Multiple ports can be queried.

* `member_id` - (Optional, List) Specifies the backend server ID.
  Multiple IDs can be queried.

* `operating_status` - (Optional, List) Specifies the operating status of the backend server.
  Value options:
  + **ONLINE**: The backend server is running normally.
  + **NO_MONITOR**: No health check is configured for the backend server group to which the backend server belongs.
  + **OFFLINE**: The cloud server used as the backend server is stopped or does not exist.
  Multiple statuses can be queried.

* `enterprise_project_id` - (Optional, List) Specifies the ID of the enterprise project.
  + If `enterprise_project_id` is not specified, resources in all enterprise projects are queried by default.
  + If `enterprise_project_id` is specified, the value can be a specific enterprise project ID or **all_granted_eps**.
  Multiple values can be queried.

* `ip_version` - (Optional, List) Specifies the IP address version supported by the backend server group.
  The value can be **v4** or **v6**.
  Multiple versions can be queried.

* `pool_id` - (Optional, List) Specifies the ID of the backend server group to which the backend server belongs.
  Multiple IDs can be queried.

* `loadbalancer_id` - (Optional, List) Specifies the ID of the load balancer with which the load balancer is associated.
  Multiple IDs can be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `members` - Indicates the list of backend servers.

  The [members](#members_struct) structure is documented below.

<a name="members_struct"></a>
The `members` block supports:

* `id` - Indicates the backend server ID.

* `name` - Indicates the backend server name.

* `member_type` - Indicates the type of the backend server.

* `address` - Indicates the private IP address bound to the backend server.

* `subnet_cidr_id` - Indicates the ID of the IPv4 or IPv6 subnet where the backend server resides.

* `protocol_port` - Indicates the port used by the backend server to receive requests.

* `weight` - Indicates the weight of the backend server.

* `operating_status` - Indicates the health status of the backend server.

* `ip_version` - Indicates the IP version supported by the backend server.

* `instance_id` - Indicates the ID of the instance associated with the backend server.

* `pool_id` - Indicates the ID of the backend server group to which the backend server belongs.

* `loadbalancer_id` - Indicates the ID of the load balancer with which the backend server is associated.

* `project_id` - Indicates the ID of the project where the backend server is used.

* `status` - Indicates the health status of the backend server.

  The [status](#members_status_struct) structure is documented below.

* `reason` - Indicates why health check fails.

  The [reason](#members_reason_struct) structure is documented below.

* `created_at` - Indicates the time when a backend server was added.

* `updated_at` - Indicates the time when a backend server was updated.

<a name="members_status_struct"></a>
The `status` block supports:

* `listener_id` - Indicates the listener ID.

* `operating_status` - Indicates the health status of the backend server.

* `reason` - Indicates why health check fails.

  The [reason](#status_reason_struct) structure is documented below.

<a name="status_reason_struct"></a>
The `reason` block supports:

* `reason_code` - Indicates the code of the health check failures.

* `expected_response` - Indicates the expected HTTP status code.

* `healthcheck_response` - Indicates the returned HTTP status code in the response.

<a name="members_reason_struct"></a>
The `reason` block supports:

* `reason_code` - Indicates the code of the health check failures.

* `expected_response` - Indicates the expected HTTP status code.

* `healthcheck_response` - Indicates the returned HTTP status code in the response.
