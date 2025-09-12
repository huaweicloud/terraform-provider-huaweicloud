---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channel_member_group"
description: |-
  Use this resource to manage a member group within HuaweiCloud.
---

# huaweicloud_apig_channel_member_group

Use this resource to manage a member group within HuaweiCloud.

## Example Usage

### Basic usage

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_group_name" {}

resource "huaweicloud_apig_channel_member_group" "test" {
  instance_id    = var.instance_id
  vpc_channel_id = var.vpc_channel_id
  name           = var.member_group_name
  weight         = 10
}
```

### Create APIG channel member group for microservice

```hcl
variable "instance_id" {}
variable "vpc_channel_id" {}
variable "member_group_name" {}

resource "huaweicloud_apig_channel_member_group" "test" {
  instance_id          = var.instance_id
  vpc_channel_id       = var.vpc_channel_id
  name                 = var.member_group_name
  weight               = 10
  microservice_version = "v1.0"
  microservice_port    = 8080

  microservice_labels {
    label_name  = "terraform_test"
    label_value = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the member group is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the member group
  belongs.

* `vpc_channel_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC channel.

* `name` - (Required, String) Specifies the name of the member group.  
  Only the Chinese characters, English letters, numbers, underscores(_), hyphens(-) and dots(.) are allowed and the
  valid value is range from `3` to `64`.
  Only start with English letters and Chinese characters.

* `description` - (Optional, String) Specifies the description of the member group.  
  The description contain a maximum of `255` characters.

* `weight` - (Optional, Int) Specifies the weight value of the member group.  
  This weight value is automatically used for weight allocation, and the valid value is range from `0` to `100`.

* `microservice_version` - (Optional, String) Specifies the version of the member group.  
  Only supported when VPC channel type is microservice.
  The microservice version contain a maximum of `255` characters.

* `microservice_port` - (Optional, Int) Specifies the port number of the member group.  
  Only supported when VPC channel type is microservice.
  When the port number is `0`, all addresses under the backend server group follow the original load balancing
  configuration.

* `microservice_labels` - (Optional, List) Specifies the microservice labels of the member group.  
  Only supported when VPC channel type is microservice.
  The [microservice_labels](#apig_channel_member_group_microservice_labels) structure is documented below.

* `reference_vpc_channel_id` - (Optional, String) Specifies the ID of the reference load balance channel.

<a name="apig_channel_member_group_microservice_labels"></a>
The `microservice_labels` block supports:

* `name` - (Required, String) The name of the microservice label.

* `value` - (Required, String) The value of the microservice label.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the member group, in RFC3339 format.

* `update_time` - The update time of the member group, in RFC3339 format.

## Import

Member groups can be imported using their `id`, the ID of the related dedicated instance and the ID of the related VPC
channel, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_apig_channel_member_groups.test <instance_id>/<vpc_channel_id>/<id>
```
