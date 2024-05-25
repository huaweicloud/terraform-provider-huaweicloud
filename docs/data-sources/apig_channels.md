---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channels"
description: |-
  Use this data source to query the VPC channels within HuaweiCloud.
---

# huaweicloud_apig_channels

Use this data source to query the VPC channels within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "channel_name" {}

data "huaweicloud_apig_channels" "test" {
  instance_id = var.instance_id
  name        = var.channel_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the channels belong.

* `channel_id` - (Optional, String) Specifies the VPC channel ID of the to be queried.

* `name` - (Optional, String) Specifies the name of the channel to be queried.

* `precise_search` - (Optional, String) Specifies the parameter name for exact matching to be queried.

* `member_group_id` - (Optional, String) Specifies the ID of the member group to be queried.

* `member_group_name` - (Optional, String) Specifies the name of the member group to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vpc_channels` - All VPC channels that match the filter parameters.
  The [vpc_channels](#vpc_channels) structure is documented below.

<a name="vpc_channels"></a>
The `vpc_channels` block supports:

* `id` - The ID of the VPC channel.

* `name` - The name of the VPC channel.

* `port` - The port of the backend server.

* `balance_strategy` - The distribution algorithm.

* `member_type` - The member type of the VPC channel.

* `type` - The type of the VPC channel.

* `created_at` - The creation time of channel, in RFC3339 format.

* `member_group` - The parameter member groups of the VPC channels.
  The [member_group](#member_group) structure is documented below.

<a name="member_group"></a>
The `member_group` block supports:

* `id` - The ID of the member group.

* `name` - The name of the member group.

* `description` - The description of the member group.

* `weight` - The weight of the current member group.

* `microservice_version` - The microservice version of the backend server group.

* `microservice_port` - The microservice port of the backend server group.

* `microservice_labels` - The microservice tags of the backend server group.
  The [microservice_labels](microservice_labels) structure is documented below.

<a name="microservice_labels"></a>
The `microservice_labels` block supports:

* `name` - The name of the microservice label.

* `value` - The value of the microservice label.
