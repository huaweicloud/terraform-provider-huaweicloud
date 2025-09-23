---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channel_member_groups"
description: |-
  Use this data source to get the list of member groups within HuaweiCloud.
---

# huaweicloud_apig_channel_member_groups

Use this data source to get the list of member groups within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_group_name" {}

data "huaweicloud_apig_channel_member_groups" "test" {
  instance_id    = var.instance_id
  vpc_channel_id = var.vpc_channel_id
  name           = var.member_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the list of member groups are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the VPC channel belongs.

* `vpc_channel_id` - (Required, String) Specifies the ID of the VPC channel to which the list of member groups belong.

* `name` - (Optional, String) Specifies the name of the member group for fuzzy matching.

* `precise_search` - (Optional, String) Specifies the list of parameter names for exact matching.  
  When this parameter contains the field(s) to be queried, the corresponding field query will become an exact-match
  query.  
  The valid value is **member_group_name**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `member_groups` - The list of the member groups that matched filter parameters.  
  The [member_groups](#apig_channel_member_groups) structure is documented below.

<a name="apig_channel_member_groups"></a>
The `member_groups` block supports:

* `id` - The ID of the member group.

* `name` - The name of the member group.

* `description` - The description of the member group.

* `weight` - The weight value of the member group.

* `microservice_version` - The microservice version of the member group.

* `microservice_port` - The microservice port of the member group.

* `microservice_labels` - The microservice labels of the member group.  
  The [microservice_labels](#apig_channel_member_groups_microservice_labels) structure is documented below.

* `reference_vpc_channel_id` - The ID of the referenced load channel.

* `create_time` - The creation time of the member group, in RFC3339 format.

* `update_time` - The update time of the member group, in RFC3339 format.

<a name="apig_channel_member_groups_microservice_labels"></a>
The `microservice_labels` block supports:

* `name` - The name of the microservice label.

* `value` - The value of the microservice label.
