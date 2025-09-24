---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channel_member"
description: |-
  Use this resource to manage a channel member within HuaweiCloud.
---

# huaweicloud_apig_channel_member

Use this resource to manage a channel member within HuaweiCloud.

## Example Usage

### Create APIG channel member by IP address

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_group_name" {}
variable "member_ip_address" {}
variable "port" {}

resource "huaweicloud_apig_channel_member" "test" {
  instance_id       = var.instance_id
  vpc_channel_id    = var.vpc_channel_id
  member_group_name = var.member_group_name
  member_ip_address = var.member_ip_address
  port              = var.port
}
```

### Create APIG channel member by ECS instance ID

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_group_name" {}
variable "ecs_id" {}
variable "ecs_name" {}
variable "port" {}

resource "huaweicloud_apig_channel_member" "test" {
  instance_id       = var.instance_id
  vpc_channel_id    = var.vpc_channel_id
  member_group_name = var.member_group_name
  ecs_id            = var.ecs_id
  ecs_name          = var.ecs_name
  port              = var.port
  is_backup         = false
  status            = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the channel member is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the channel
  member belongs.

* `vpc_channel_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC channel.

* `member_group_name` - (Required, String, NonUpdatable) Specifies the name of the channel member group.  
  This is used to select a channel member group for easy unified modification of the corresponding server group backend
  attributes.

* `port` - (Required, Int, NonUpdatable) Specifies the port number of the channel member.  
  The valid value is range from `0` to `65,535`.

* `member_ip_address` - (Optional, String) Specifies the IP address of the channel member.  
  The member_ip_address contain a maximum of `255` characters.

  -> Required if the type of vpc channel is **ip**.

* `ecs_id` - (Optional, String) Specifies the ID of the ECS instance.  
  Only the English letters, numbers, underscores(_) and hyphens(-) are allowed and the valid value is range
  from `0` to `255`.

* `ecs_name` - (Optional, String) Specifies the name of the ECS instance.  
  Only the Chinese characters, English letters, numbers, underscores(_), hyphens(-) and dots(.) are allowed
  and the valid value is range from `0` to `255`.

-> Parameter `ecs_id` and `ecs_name` are required if the type of vpc channel is **ecs**.

* `status` - (Optional, Int) Specifies the status of the channel member.  
  The valid values are as follow:
  + **1**: Available
  + **2**: Unavailable

* `is_backup` - (Optional, Bool) Specifies whether this member is the backup member.  
  When enabled, the corresponding backend service is a backup node and only works when all non-backup nodes fail.  
  The default value is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `weight` - The weight value of the channel member.

* `create_time` - The creation time of the channel member, in RFC3339 format.

* `member_group_id` - The ID of the member group.

* `health_status` - The health status of the channel member.

## Import

Channel members can be imported using their `id`, the ID of the related dedicated instance, the ID of the related VPC
channel and the name of the related member group, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_apig_channel_member.test <instance_id>/<vpc_channel_id>/<member_group_name>/<id>
```
