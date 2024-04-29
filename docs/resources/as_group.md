---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_group"
description: ""
---

# huaweicloud_as_group

Manages an AS group resource within HuaweiCloud.

## Example Usage

### Basic Autoscaling Group

```hcl
variable "configuration_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_as_group" "my_as_group" {
  scaling_group_name       = "my_as_group"
  scaling_configuration_id = var.configuration_id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = var.vpc_id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = var.subnet_id
  }
}
```

### Autoscaling Group with tags

```hcl
variable "configuration_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_as_group" "my_as_group_tags" {
  scaling_group_name       = "my_as_group_tags"
  scaling_configuration_id = var.configuration_id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = var.vpc_id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = var.subnet_id
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Autoscaling Group Only Remove Members When Scaling Down

```hcl
variable "configuration_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_as_group" "my_as_group_only_remove_members" {
  scaling_group_name       = "my_as_group_only_remove_members"
  scaling_configuration_id = var.configuration_id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = var.vpc_id
  delete_publicip          = true
  delete_instances         = "no"

  networks {
    id = var.subnet_id
  }
}
```

### Autoscaling Group With Elastic Load Balancer Listener

```hcl
variable "configuration_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "ipv4_subnet_id" {}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = var.ipv4_subnet_id
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
  scaling_configuration_id = var.configuration_id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = var.vpc_id

  networks {
    id = var.subnet_id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the AS group.
  If omitted, the provider-level region will be used. Changing this creates a new group.

* `scaling_group_name` - (Required, String) Specifies the name of the scaling group. The name can contain
  letters, digits, underscores(_), and hyphens(-),and cannot exceed 64 characters.

* `scaling_configuration_id` - (Required, String) Specifies the configuration ID which defines configurations
  of instances in the AS group.

* `desire_instance_number` - (Optional, Int) Specifies the expected number of instances. The default value is the
  minimum number of instances. The value ranges from the minimum number of instances to the maximum number of instances.

* `min_instance_number` - (Optional, Int) Specifies the minimum number of instances. The default value is 0.

* `max_instance_number` - (Optional, Int) Specifies the maximum number of instances. The default value is 0.

* `cool_down_time` - (Optional, Int) Specifies the cooling duration (in seconds). The value ranges from 0 to 86400,
  and is 300 by default.

* `availability_zones` - (Optional, List) Specifies the availability zones in which to create the instances in the
  autoscaling group.

* `multi_az_scaling_policy` - (Optional, String) Specifies the priority policy used to select target AZs when adjusting
  the number of instances in an AS group. The value can be `EQUILIBRIUM_DISTRIBUTE` and `PICK_FIRST`.

  + **EQUILIBRIUM_DISTRIBUTE** (default): When adjusting the number of instances, ensure that instances in each AZ in the
    availability_zones list is evenly distributed. If instances cannot be added in the target AZ, select another AZ based
    on the PICK_FIRST policy.
  + **PICK_FIRST**: When adjusting the number of instances, target AZs are determined in the order in the
    availability_zones list.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this creates a new group.

* `networks` - (Required, List) Specifies an array of one or more network IDs. The system supports up to five networks.
  The [object](#group_network_object) structure is documented below.

* `security_groups` - (Optional, List) Specifies an array of one or more security group IDs to associate with the group.
  The [object](#group_security_group_object) structure is documented below.

  -> If the security group is specified both in the AS configuration and AS group, scaled ECS instances will be added to
  the security group specified in the AS configuration. If the security group is not specified in either of them, scaled
  ECS instances will be added to the default security group. For your convenience, you are advised to specify the security
  group in the AS configuration.

* `lbaas_listeners` - (Optional, List) Specifies an array of one or more enhanced load balancer. The system supports
  the binding of up to six load balancers. The [object](#group_lbaas_listener_object) structure is documented below.

* `health_periodic_audit_method` - (Optional, String) Specifies the health check method for instances in the AS group.
  The health check methods include `ELB_AUDIT` and `NOVA_AUDIT`. If load balancing is configured, the default value of
  this parameter is `ELB_AUDIT`. Otherwise, the default value is `NOVA_AUDIT`.

* `health_periodic_audit_time` - (Optional, Int) Specifies the health check period for instances. The unit is minute
  and value includes 0, 1, 5 (default), 15, 60, and 180. If the value is set to 0, health check is performed every 10 seconds.

* `health_periodic_audit_grace_period` - (Optional, Int) Specifies the health check grace period for instances.
  The unit is second and the value ranges from 0 to 86400. The default value is 600.

  -> This parameter is valid only when the instance health check method of the AS group is `ELB_AUDIT`.

* `instance_terminate_policy` - (Optional, String) Specifies the instance removal policy. The policy has four
  options: `OLD_CONFIG_OLD_INSTANCE` (default), `OLD_CONFIG_NEW_INSTANCE`, `OLD_INSTANCE`, and `NEW_INSTANCE`.

  + **OLD_CONFIG_OLD_INSTANCE** (default): The earlier-created instances based on the earlier-created AS configurations
    are removed first.
  + **OLD_CONFIG_NEW_INSTANCE**: The later-created instances based on the earlier-created AS configurations are removed first.
  + **OLD_INSTANCE**: The earlier-created instances are removed first.
  + **NEW_INSTANCE**: The later-created instances are removed first.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the AS group.

* `description` - (Optional, String) Specifies the description of the AS group.
  The value can contain 0 to 256 characters.

* `agency_name` - (Optional, String) Specifies the IAM agency name. If you change the agency,
  the new agency will be available for ECSs scaled out after the change.

* `delete_publicip` - (Optional, Bool) Specifies whether to delete the elastic IP address bound to the instances of
  AS group when deleting the instances. The options are `true` and `false`.

* `delete_instances` - (Optional, String) Specifies whether to delete the instances in the AS group when deleting
  the AS group. The options are `yes` and `no`.

* `force_delete` - (Optional, Bool) Specifies whether to forcibly delete the AS group, remove the ECS instances and
  release them. The default value is `false`.

* `enable` - (Optional, Bool) Specifies whether to enable the AS Group. The options are `true` and `false`.
  The default value is `true`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the AS group.

<a name="group_network_object"></a>
The `networks` block supports:

* `id` - (Required, String) Specifies the subnet ID.

* `ipv6_enable` - (Optional, Bool) Specifies whether to support IPv6 addresses. The default value is `false`.

* `ipv6_bandwidth_id` - (Optional, String) Specifies the ID of the shared bandwidth of an IPv6 address.

* `source_dest_check` - (Optional, Bool) Specifies whether process only traffic that is destined specifically
  for it. Defaults to true.

<a name="group_security_group_object"></a>
The `security_groups` block supports:

* `id` - (Required, String) Specifies the ID of the security group.

<a name="group_lbaas_listener_object"></a>
The `lbaas_listeners` block supports:

* `pool_id` - (Required, String) Specifies the backend ECS group ID.

* `protocol_port` - (Required, Int) Specifies the backend protocol, which is the port on which a backend ECS listens for
  traffic. The number of the port ranges from 1 to 65535.

* `weight` - (Optional, Int) Specifies the weight, which determines the portion of requests a backend ECS processes
  compared to other backend ECSs added to the same listener. The value of this parameter ranges from 0 to 100. The
  default value is 1.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The AS group ID.

* `status` - The status of the AS group.

* `current_instance_number` - The number of current instances in the AS group.

* `instances` - The instances IDs of the AS group.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

AS groups can be imported by their `id`. For example,

```
terraform import huaweicloud_as_group.my_as_group 9ec5bea6-a728-4082-8109-5a7dc5c7af74
```
