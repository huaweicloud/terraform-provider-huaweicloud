---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_groups"
description: |-
  Use this data source to get a list of AS groups.
---

# huaweicloud_as_groups

Use this data source to get a list of AS groups.

```hcl
data "huaweicloud_as_groups" "groups" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the AS groups.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the AS group name. Fuzzy search is supported.

* `scaling_configuration_id` - (Optional, String) Specifies the AS configuration ID, which can be obtained using
  the API for listing AS configurations.

* `status` - (Optional, String) Specifies the AS group status. The options are as follows:
  - **INSERVICE**: indicates that the AS group is functional.
  - **PAUSED**: indicates that the AS group is paused.
  - **ERROR**: indicates that the AS group malfunctions.
  - **DELETING**: indicates that the AS group is being deleted.
  - **FREEZED**: indicates that the AS group has been frozen.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the list.

* `groups` - A list of AS groups.

The `groups` block supports:

* `scaling_group_name` - The name of the AS group.

* `scaling_group_id` - The AS group ID.

* `status` - The AS group status.

* `scaling_configuration_id` - The AS configuration ID.

* `scaling_configuration_name` - The AS configuration name.

* `current_instance_number` - The number of current instances in the AS group.

* `desire_instance_number` - The expected number of instances in the AS group.

* `min_instance_number` - The minimum number of instances in the AS group.

* `max_instance_number` - The maximum number of instances in the AS group.

* `cool_down_time` - The cooling duration, in seconds..

* `lbaas_listeners` - The enhanced load balancers.
  The [object](#lbaas_listener_object) structure is documented below.

* `availability_zones` - The AZ information.

* `networks` - The network information.
  The [object](#network_object) structure is documented below.

* `security_groups` - The security group information.
  The [object](#security_group_object) structure is documented below.

* `created_at` - The time when an AS group was created. The time format complies with UTC.

* `vpc_id` - The ID of the VPC to which the AS group belongs.

* `detail` - Details about the AS group. If a scaling action fails, this parameter is used to record errors.

* `is_scaling` - The scaling flag of the AS group.

* `health_periodic_audit_method` - The health check method.

* `health_periodic_audit_time` - The health check interval.

* `health_periodic_audit_grace_period` - The grace period for health check.

* `instance_terminate_policy` - The instance removal policy.

* `delete_publicip` - Whether to delete the EIP bound to the ECS when deleting the ECS.

* `delete_volume` - Whether to delete the data disks attached to the ECS when deleting the ECS.

* `enterprise_project_id` - The enterprise project ID.

* `activity_type` - The type of the AS action.

* `multi_az_scaling_policy` - The priority policy used to select target AZs when adjusting the number of
  instances in an AS group.

* `description` - The description of the AS group.

* `iam_agency_name` - The agency name.

* `tags` - The tag of AS group.

* `instances` - The scaling group instances ids.

<a name="lbaas_listener_object"></a>
The `lbaas_listeners` block supports:

* `pool_id` - The backend ECS group ID.

* `protocol_port` - The backend protocol ID, which is the port on which a backend ECS listens for traffic.

* `weight` - The weight, which determines the portion of requests a backend ECS processes
  compared to other backend ECSs added to the same listener.

* `protocol_version` - The version of IP addresses of backend servers to be bound with the ELB.

* `listener_id` - The ID of the listener associate with the ELB.

<a name="network_object"></a>
The `networks` block supports:

* `id` - The subnet ID.

* `ipv6_enable` - Specifies whether to support IPv6 addresses.

* `ipv6_bandwidth_id` - The ID of the shared bandwidth of an IPv6 address.

* `source_dest_check` - Whether processing only traffic that is destined specifically for it.

<a name="security_group_object"></a>
The `security_groups` block supports:

* `id` - The ID of the security group.
