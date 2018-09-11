---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_group_v1"
sidebar_current: "docs-huaweicloud-resource-as-group-v1"
description: |-
  Manages a V1 Autoscaling Group resource within HuaweiCloud.
---

# huaweicloud\_as\_group_v1

Manages a V1 Autoscaling Group resource within HuaweiCloud.

## Example Usage

### Basic Autoscaling Group

```hcl
resource "huaweicloud_as_group_v1" "my_as_group" {
  scaling_group_name = "my_as_group"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number = 2
  min_instance_number = 0
  max_instance_number = 10
  networks = [{id = "ad091b52-742f-469e-8f3c-fd81cadf0743"}]
  security_groups = [{id = "45e4c6de-6bf0-4843-8953-2babde3d4810"}]
  vpc_id = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip = true
  delete_instances = "yes"
}
```

### Autoscaling Group Only Remove Members When Scaling Down

```hcl
resource "huaweicloud_as_group_v1" "my_as_group_only_remove_members" {
  scaling_group_name = "my_as_group_only_remove_members"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number = 2
  min_instance_number = 0
  max_instance_number = 10
  networks = [{id = "ad091b52-742f-469e-8f3c-fd81cadf0743"}]
  security_groups = [{id = "45e4c6de-6bf0-4843-8953-2babde3d4810"}]
  vpc_id = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip = true
  delete_instances = "no"
}
```

### Autoscaling Group With ELB Listener

```hcl
resource "huaweicloud_as_group_v1" "my_as_group_with_elb" {
  scaling_group_name = "my_as_group_with_elb"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number = 2
  min_instance_number = 0
  max_instance_number = 10
  networks = [{id = "ad091b52-742f-469e-8f3c-fd81cadf0743"}]
  security_groups = [{id = "45e4c6de-6bf0-4843-8953-2babde3d4810"}]
  vpc_id = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  lb_listener_id = "${huaweicloud_elb_listener.my_listener.id}"
  delete_publicip = true
  delete_instances = "yes"
}

resource "huaweicloud_elb_listener" "my_listener" {
  name = "my_listener"
  description = "my test listener"
  protocol = "TCP"
  backend_protocol = "TCP"
  port = 12345
  backend_port = 21345
  lb_algorithm = "roundrobin"
  loadbalancer_id = "cba48790-baf5-4446-adb3-02069a916e97"
  timeouts {
        create = "5m"
        update = "5m"
        delete = "5m"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the AS group. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new AS group.

* `scaling_group_name` - (Required) The name of the scaling group. The name can contain letters,
    digits, underscores(_), and hyphens(-),and cannot exceed 64 characters.

* `scaling_configuration_id` - (Optional) The configuration ID which defines
    configurations of instances in the AS group.

* `desire_instance_number` - (Optional) The expected number of instances. The default
    value is the minimum number of instances. The value ranges from the minimum number of
    instances to the maximum number of instances.

* `min_instance_number` - (Optional) The minimum number of instances.
    The default value is 0.

* `max_instance_number` - (Optional) The maximum number of instances.
    The default value is 0.

* `cool_down_time` - (Optional) The cooling duration (in seconds). The value ranges
    from 0 to 86400, and is 900 by default.

* `lb_listener_id` - (Optional) The ELB listener IDs. The system supports up to
    three ELB listeners, the IDs of which are separated using a comma (,).

* `available_zones` - (Optional) The availability zones in which to create
    the instances in the autoscaling group.

* `networks` - (Required) An array of one or more network IDs.
    The system supports up to five networks. The networks object structure
    is documented below.

* `security_groups` - (Required) An array of one or more security group IDs
    to associate with the group. The security_groups object structure is
    documented below.

* `vpc_id` - (Required) The VPC ID. Changing this creates a new group.

* `health_periodic_audit_method` - (Optional) The health check method for instances
    in the AS group. The health check methods include `ELB_AUDIT` and `NOVA_AUDIT`.
    If load balancing is configured, the default value of this parameter is `ELB_AUDIT`.
    Otherwise, the default value is `NOVA_AUDIT`.

* `health_periodic_audit_time` - (Optional) The health check period for instances.
    The period has four options: 5 minutes (default), 15 minutes, 60 minutes, and 180 minutes.

* `instance_terminate_policy` - (Optional) The instance removal policy. The policy has
    four options: `OLD_CONFIG_OLD_INSTANCE` (default), `OLD_CONFIG_NEW_INSTANCE`,
    `OLD_INSTANCE`, and `NEW_INSTANCE`.

* `notifications` - (Optional) The notification mode. The system only supports `EMAIL`
    mode which refers to notification by email.

* `delete_publicip` - (Optional) Whether to delete the elastic IP address bound to the
    instances of AS group when deleting the instances. The options are `true` and `false`.

* `delete_instances` - (Optional) Whether to delete the instances in the AS group
    when deleting the AS group. The options are `yes` and `no`.

The `networks` block supports:

* `id` - (Required) The network UUID.

The `security_groups` block supports:

* `id` - (Required) The UUID of the security group.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `scaling_group_name` - See Argument Reference above.
* `desire_instance_number` - See Argument Reference above.
* `min_instance_number` - See Argument Reference above.
* `max_instance_number` - See Argument Reference above.
* `cool_down_time` - See Argument Reference above.
* `lb_listener_id` - See Argument Reference above.
* `health_periodic_audit_method` - See Argument Reference above.
* `health_periodic_audit_time` - See Argument Reference above.
* `instance_terminate_policy` - See Argument Reference above.
* `scaling_configuration_id` - See Argument Reference above.
* `delete_publicip` - See Argument Reference above.
* `notifications` - See Argument Reference above.
* `instances` - The instances IDs of the AS group.
