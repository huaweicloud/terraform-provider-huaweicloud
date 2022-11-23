---
subcategory: "Auto Scaling"
---

# huaweicloud_as_group

Manages a Autoscaling Group resource within HuaweiCloud.

## Example Usage

### Basic Autoscaling Group

```hcl
resource "huaweicloud_as_group" "my_as_group" {
  scaling_group_name       = "my_as_group"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }
}
```

### Autoscaling Group with tags

```hcl
resource "huaweicloud_as_group" "my_as_group_tags" {
  scaling_group_name       = "my_as_group_tags"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Autoscaling Group Only Remove Members When Scaling Down

```hcl
resource "huaweicloud_as_group" "my_as_group_only_remove_members" {
  scaling_group_name       = "my_as_group_only_remove_members"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip          = true
  delete_instances         = "no"

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }
}
```

### Autoscaling Group With Enhanced Load Balancer Listener

```hcl
resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_as_group" "my_as_group_with_enhanced_lb" {
  scaling_group_name       = "my_as_group_with_enhanced_lb"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the AS group. If omitted, the `region` argument
  of the provider is used. Changing this creates a new AS group.

* `scaling_group_name` - (Required, String) The name of the scaling group. The name can contain letters, digits,
  underscores(_), and hyphens(-),and cannot exceed 64 characters.

* `scaling_configuration_id` - (Required, String) The configuration ID which defines configurations of instances in the
  AS group.

* `desire_instance_number` - (Optional, Int) The expected number of instances. The default value is the minimum number
  of instances. The value ranges from the minimum number of instances to the maximum number of instances.

* `min_instance_number` - (Optional, Int) The minimum number of instances. The default value is 0.

* `max_instance_number` - (Optional, Int) The maximum number of instances. The default value is 0.

* `cool_down_time` - (Optional, Int) The cooling duration (in seconds). The value ranges from 0 to 86400, and is 300 by
  default.

* `lbaas_listeners` - (Optional, List) An array of one or more enhanced load balancer. The system supports the binding
  of up to six load balancers. The object structure is documented below.

* `availability_zones` - (Optional, List) The availability zones in which to create the instances in the autoscaling group.

* `multi_az_scaling_policy` - (Optional, String) The priority policy used to select target AZs when adjusting the number
  of instances in an AS group. The value can be `EQUILIBRIUM_DISTRIBUTE` and `PICK_FIRST`.

  + **EQUILIBRIUM_DISTRIBUTE** (default): When adjusting the number of instances, ensure that instances in each AZ in the
    availability_zones list is evenly distributed. If instances cannot be added in the target AZ, select another AZ based
    on the PICK_FIRST policy.
  + **PICK_FIRST**: When adjusting the number of instances, target AZs are determined in the order in the availability_zones list.

* `vpc_id` - (Required, String, ForceNew) The VPC ID. Changing this creates a new group.

* `networks` - (Required, List) An array of one or more network IDs. The system supports up to five networks. The
  networks object structure is documented below.

* `security_groups` - (Optional, List) An array of one or more security group IDs to associate with the group. The
  security_groups object structure is documented below.

  -> If the security group is specified both in the AS configuration and AS group, scaled ECS instances will be added to
  the security group specified in the AS configuration. If the security group is not specified in either of them, scaled
  ECS instances will be added to the default security group. For your convenience, you are advised to specify the security
  group in the AS configuration.

* `health_periodic_audit_method` - (Optional, String) The health check method for instances in the AS group. The health
  check methods include `ELB_AUDIT` and `NOVA_AUDIT`. If load balancing is configured, the default value of this
  parameter is `ELB_AUDIT`. Otherwise, the default value is `NOVA_AUDIT`.

* `health_periodic_audit_time` - (Optional, Int) The health check period for instances. The period has four options: 5
  minutes (default), 15 minutes, 60 minutes, and 180 minutes.

* `health_periodic_audit_grace_period` - (Optional, Int) The health check grace period for instances. The unit is second
  and the value ranges from 0 to 86400. The default value is 600.

  -> This parameter is valid only when the instance health check method of the AS group is `ELB_AUDIT`.

* `instance_terminate_policy` - (Optional, String) The instance removal policy. The policy has four
  options: `OLD_CONFIG_OLD_INSTANCE` (default), `OLD_CONFIG_NEW_INSTANCE`,
  `OLD_INSTANCE`, and `NEW_INSTANCE`.

* `tags` - (Optional, Map) The key/value pairs to associate with the scaling group.

* `description` (Optional, String) The description of the AS group. The value can contain 0 to 256 characters.

* `notifications` - (Optional, List) The notification mode. The system only supports `EMAIL`
  mode which refers to notification by email.

* `delete_publicip` - (Optional, Bool) Whether to delete the elastic IP address bound to the instances of AS group when
  deleting the instances. The options are `true` and `false`.

* `delete_instances` - (Optional, String) Whether to delete the instances in the AS group when deleting the AS group.
  The options are `yes` and `no`.

* `force_delete` - (Optional, Bool) Whether to forcibly delete the AS group, remove the ECS instances and release them.
  The default value is `false`.

* `enable` - (Optional, Bool) Whether to enable the AS Group. The options are `true` and `false`. The default value
  is `true`.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the AS group.

The `networks` block supports:

* `id` - (Required, String) The network UUID.

* `ipv6_enable` - (Optional, Bool) Whether to support IPv6 addresses. The default value is false.

* `ipv6_bandwidth_id` - (Optional, String) The ID of the shared bandwidth of an IPv6 address.

The `security_groups` block supports:

* `id` - (Required, String) The UUID of the security group.

The `lbaas_listeners` block supports:

* `pool_id` - (Required, String) Specifies the backend ECS group ID.
* `protocol_port` - (Required, Int) Specifies the backend protocol, which is the port on which a backend ECS listens for
  traffic. The number of the port ranges from 1 to 65535.
* `weight` - (Optional, Int) Specifies the weight, which determines the portion of requests a backend ECS processes
  compared to other backend ECSs added to the same listener. The value of this parameter ranges from 0 to 100. The
  default value is 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The AS group ID.

* `status` - Indicates the status of the AS group.

* `current_instance_number` - Indicates the number of current instances in the AS group.

* `instances` - The instances IDs of the AS group.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

AS groups can be imported by their `id`. For example,

```
terraform import huaweicloud_as_group.my_as_group 9ec5bea6-a728-4082-8109-5a7dc5c7af74
```
