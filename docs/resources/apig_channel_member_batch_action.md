---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channel_member_batch_action"
description: |-
   Use this resource to batch operate the status of channel members in HuaweiCloud.
---

# huaweicloud_apig_channel_member_batch_action

Use this resource to batch operate the status of channel members in HuaweiCloud.

-> This resource is only a one-time action resource for performing an enable/disable operation with
   the VPC channel member list. Deleting this resource will not clear the corresponding request record,
   but will only remove the resource information from the tfstate file.

## Example Usage

### Enable VPC Channel Members

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_ids" {
  type = list(string)
}

resource "huaweicloud_apig_channel_member_batch_action" "enable_members" {
  instance_id    = var.instance_id
  vpc_channel_id = var.vpc_channel_id
  action         = "enable"
  member_ids     = var.member_ids
}
```

### Disable VPC Channel Members

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_ids" {
  type = list(string)
}

resource "huaweicloud_apig_channel_member_batch_action" "disable_members" {
  instance_id    = var.instance_id
  vpc_channel_id = var.vpc_channel_id
  action         = "disable"
  member_ids     = var.member_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the channel members to be operated are located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the VPC channel
  belongs.

* `vpc_channel_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC channel to which the members belong.

* `action` - (Required, String, NonUpdatable) Specifies the batch operation for the VPC channel members.  
  The valid values are as follows:
  + **enable**: Enable the specified VPC channel members
  + **disable**: Disable the specified VPC channel members

* `member_ids` - (Required, List, NonUpdatable) The list of member IDs to be batch operated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the batch action resource.
