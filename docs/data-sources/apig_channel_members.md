---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channel_members"
description: |-
  Use this data source to query the list of channel members within HuaweiCloud.
---

# huaweicloud_apig_channel_members

Use this data source to query the list of channel members within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_name" {}

data "huaweicloud_apig_channel_members" "test" {
  instance_id    = var.instance_id
  vpc_channel_id = var.vpc_channel_id
  name           = var.member_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the channel members are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the channel members belong.

* `vpc_channel_id` - (Required, String) Specifies the ID of the VPC channel to which the members belong.

* `name` - (Optional, String) Specifies the name of the channel member to be queried for fuzzy matching.

* `member_group_name` - (Optional, String) Specifies the name of the channel member group to be queried for fuzzy
  matching.

* `member_group_id` - (Optional, String) Specifies the ID of the channel member group to be queried.

* `precise_search` - (Optional, String) Specifies the parameter name for exact matching to be queried.  
  When this parameter contains the field(s) to be queried, the corresponding field query will become an exact-match
  query.  
  The valid values are as follows:
  + **name**
  + **member_group_name**
  When performing an exact query with multiple fields, separate the fields with commas(,).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `members` - The list of the channel members that matched filter parameters.  
  The [members](#apig_channel_members) structure is documented below.

<a name="apig_channel_members"></a>
The `members` block supports:

* `id` - The ID of the channel member.

* `vpc_channel_id` - The ID of the VPC channel.

* `member_group_name` - The name of the channel member group.

* `member_group_id` - The ID of the channel member group.

* `member_ip_address` - The IP address of the channel member.

* `ecs_id` - The ID of the ECS instance.

* `ecs_name` - The name of the ECS instance.

* `port` - The port of the channel member.

* `is_backup` - Whether the channel member is a backup node.

* `status` - The status of the channel member.

* `weight` - The weight value of the channel member.

* `health_status` - The health status of the channel member.

* `create_time` - The time when the channel member was added to the VPC channel, in RFC3339 format.
