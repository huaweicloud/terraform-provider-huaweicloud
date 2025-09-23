---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_instances"
description: |-
  Use this data source to get a list of DBSS instances.
---

# huaweicloud_dbss_instances

Use this data source to get a list of DBSS instances.

## Example Usage

```hcl
data "huaweicloud_dbss_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The instance information list.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `name` - The instance name.

* `remain_days` - Days to expiry.

* `availability_zone` - The availability zone.

* `created_at` - The creation time.

* `resource_id` - The resource ID.

* `resource_spec_code` - The instance specifications.

* `version` - The current version of the instance.

* `ram` - The memory size.

* `region` - The region where the instance is located.

* `scene` - The scene information.

* `subnet_id` - The subnet ID.

* `vpc_id` - The VPC ID.

* `description` - The remark information.

* `keep_days` - Days until expiration.

* `specification` - The instance specifications.

* `effect` - The effects of freezing. Valid values are as follows:
  + **1**: The resource is frozen and can be released.
  + **2**: The resource is frozen and cannot be released.
  + **3**: The resource is frozen and cannot be renewed.

* `new_version` - The new version. If a value is returned, the upgrade is required.

* `security_group_id` - The security group ID.

* `task` - The task status. Its value can be:
  + **powering-on**: The instance is being started and can be bound or unbound.
  + **powering-off**: The instance is being stopped and can be bound or unbound.
  + **rebooting**: The instance is being restarted and can be bound or unbound.
  + **delete_wait**: The instance is waiting to be deleted and no operations are allowed on the cluster or instance.
  + **NO_TASK**: The instance is not displayed.

* `connect_ip` - The connection IP address.

* `cpu` - The number of CPUs.

* `charge_model` - The payment mode. Its value can be **Period** (yearly/monthly) or **Demand** (pay-per-use).

* `instance_id` - The instance ID.

* `port_id` - The ID of the port that the EIP is bound to.

* `status` - The instance status. Its value can be:
  + **SHUTOFF**: Disabled.
  + **ACTIVE**: Operations allowed.
  + **DELETING**: No operations allowed.
  + **BUILD**: No operations allowed.
  + **DELETED**: Not displayed.
  + **ERROR**: Only deletion allowed.
  + **HAWAIT**: Waiting for the standby to be created; No operations allowed.
  + **FROZEN**: Only renewal, binding, and unbinding allowed.
  + **UPGRADING**: No operations allowed.

* `config_num` - The total number of configured databases.

* `database_limit` - The total number of supported databases.

* `expired_at` - The expired time.

* `connect_ipv6` - The IPv6 address.
