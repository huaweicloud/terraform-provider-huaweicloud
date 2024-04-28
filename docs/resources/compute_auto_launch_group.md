---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_auto_launch_group"
description: ""
---

# huaweicloud_compute_auto_launch_group

Manages an ECS auto launch group resource within HuaweiCloud.

## Example Usage

```hcl
variable "auto_launch_group_name" {}
variable "launch_template_id" {}
variable "launch_template_version" {}
variable "availability_zone" {}
variable "flavor_id" {}

resource "huaweicloud_compute_auto_launch_group" "test" {
  name                    = var.auto_launch_group_name
  target_capacity         = 2
  stable_capacity         = 2
  launch_template_id      = var.launch_template_id
  launch_template_version = var.launch_template_version

  overrides {
    availability_zone = var.availability_zone
    flavor_id         = var.flavor_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the auto launch group.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the auto launch group. The valid length is limited
  between `1` to `64`, Only Chinese and English letters, digits, hyphens (-) and underscores (_) are allowed.

* `target_capacity` - (Required, Int) Specifies the target capacity of the auto launch group, the unit is the number of
  vCPU, and the value must be bigger than or equal to `stable_capacity`. The capacity of the spot instance euqals to
  full capacity minus `stable_capacity`.

* `launch_template_id` - (Required, String, ForceNew) Specifies the ID of launch template for instance.
  Changing this creates a new resource.

* `launch_template_version` - (Required, String, ForceNew) Specifies the version of launch template for instance.
  Changing this creates a new resource.

* `overrides` - (Required, List, ForceNew) Specifies the instance details. Supporting mutiple `overrides` to create
  instances of different specification. Changing this creates a new resource.
  The [overrides](#block--overrides) structure is documented below.

* `allocation_strategy` - (Optional, String, ForceNew) Specifies the allocation strategy of the auto launch group.

  Valid values are:
  + **lowest_price**: Lowest price strategy, the sum of the prices of all instances launched by the auto launch group
    is the lowest.
  + **prioritized**: Priority strategy, create instances according to the priority set by the specifications.
  + **capacity_optimized**: Capacity optimization strategy. Instances launched by auto launch group are launched first
    according to large specifications.

  Default is **lowest_price**. Changing this creates a new resource.

* `delete_instances` - (Optional, String) Specifies the interruption behavior of instances when deleting the auto launch
  group.

  Valid values are:
  + **terminate**: Depending on `delete_publicip` and `delete_volume` to determine whether to release the elastic public
    IP and disk.
  + **noTermination**: The elastic public IP and disk are both not released.

  Default is **terminate**.

* `excess_fulfilled_capacity_behavior` - (Optional, String) Specifies the interruption behavior of instances when target
  capacity is exceeded or reduced. Valid values are **terminate** and **noTermination**. Default is **terminate**.

* `instances_behavior_with_expiration` - (Optional, String) Specifies the interruption behavior of running instances
  when requests expire. Valid values are **terminate** and **noTermination**. Default is **terminate**.

* `spot_price` - (Optional, Float) Specifies the highest price a user is willing to pay per hour for a Spot instance.
  If no price is provided in `overrides`, this price can be used.

* `stable_capacity` - (Optional, Int) Specifies the target capacity to on-demand instance, the unit is the number of
  vCPU, and the value must be less than or equal to `target_capacity`. There can be no on-demand instances in the auto
  launch group.

* `supply_option` - (Optional, String, ForceNew) Specifies the selection strategies in resource supply.

  Valid values are:
  + **singlation**: Select a specification to supply.
  + **multiple**: Combine multiple specifications to supply.

  Changing this creates a new resource.

* `type` - (Optional, String, ForceNew) Specifies the request type.

  Valid values are:
  + **request**: One-time. The instance cluster is only delivered at startup and will not be retried after scheduling
    failure.
  + **maintain**: Continuous supply. Try to deliver the instance cluster at startup and monitor the capacity. If the
    target capacity is not reached, try to continue creating ECS ​​instances.

  Changing this creates a new resource.

* `valid_since` - (Optional, String, ForceNew) Specifies the request start time, together with `valid_since`, determines
  the validity period. In the format of **yyyy-MM-dd'T'HH:mm:ssZ**, an ISO 8601 time format. Default is starting now.
  Changing this creates a new resource.

* `valid_until` - (Optional, String, ForceNew) Specifies the request end time, together with `valid_since`, determines
  the validity period. In the format of **yyyy-MM-dd'T'HH:mm:ssZ**, an ISO 8601 time format. Default is never come to an
  end. Changing this creates a new resource.

<a name="block--overrides"></a>
The `overrides` block supports:

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone which the instance in.
  Please refer to the document link [reference](https://developer.huaweicloud.com/intl/en-us/endpoint/?ECS) for values.
  Changing this creates a new resource.

* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID of the instance. You can get available flavor id
  through data source `huaweicloud_compute_flavors`.
  Changing this creates a new resource.

* `priority` - (Optional, Int, ForceNew) Specifies the priority for launching. The smaller the value, the higher the
  priority, and the launch will be given first. Valid value is from zero to max value of integer.
  Changing this creates a new resource.

* `spot_price` - (Optional, Float, ForceNew) Specifies the highest price a user is willing to pay per hour for a Spot
  instance. Changing this creates a new resource.

* `weighted_capacity` - (Optional, Float, ForceNew) Specifies the weight of the instance specification. The higher the
  value, the greater the ability of a single instance to meet computing power requirements, and the smaller the number
  of instances required. It must be bigger than zero. The weight value can be calculated based on the computing power of
  the specified instance specification and the minimum computing power of a single node in the cluster.

  Assuming that the minimum computing power of a single node is 8vCPU and 60GB, the weight of the 8vCPU and 60GB
  instance specification can be set to 1, and the weight of the 16vCPU and 120GB instance specification can be set to 2.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time of the auto launch group.

* `current_capacity` - The total computing power that has been purchased successfully.

* `current_stable_capacity` - The on-demand computing power has been successfully purchased.

* `status` - The status of the auto launch group.

* `task_state` - The status of the auto launch group task. It can be:
  + **HANDLING**: Launching.
  + **FULFILLED**: The auto launch group task is fully equipped.
  + **ERROR**: Error occurs in the auto launch group task.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 5 minutes.

## Import

The auto launch group can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_compute_auto_launch_group.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `task_state`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_compute_auto_launch_group" "test" {
    ...

  lifecycle {
    ignore_changes = [
      task_state,
    ]
  }
}
```
