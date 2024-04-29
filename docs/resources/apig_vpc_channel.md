---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_vpc_channel"
description: ""
---

# huaweicloud_apig_vpc_channel

!> **WARNING:** It has been deprecated.

Manages a VPC channel resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "ecs_id1" {}
variable "ecs_id2" {}

resource "huaweicloud_apig_vpc_channel" "test" {
  instance_id = var.instance_id
  name        = var.app_name
  port        = 8080
  protocol    = "HTTPS"
  path        = "/"
  http_code   = "201,202,203"

  members {
    id     = var.ecs_id1
    weight = 30
  }

  members {
    id     = var.ecs_id2
    weight = 70
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the VPC channel is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the VPC channel
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the VPC channel.  
  The valid length is limited from `3` to `64`, only chinese and english letters, digits, hyphens (-), underscores (_)
  and dots (.) are allowed.  
  The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or **Unicode**
  format.

* `port` - (Required, Int) Specifies the host port of the VPC channel.  
  The valid value ranges from `1` to `65,535`.

* `members` - (Required, List) Specifies the configuration of the backend servers that bind the VPC channel.  
  The [object](#vpc_channel_members) structure is documented below.

* `member_type` - (Optional, String) Specifies the member type of the VPC channel.  
  The valid types are **ECS** and **EIP**, defaults to **ECS**.

* `algorithm` - (Optional, String) Specifies the distribution algorithm.  
  The valid types are **WRR**, **WLC**, **SH** and **URI hashing**, defaults to **WRR**.

* `protocol` - (Optional, String) Specifies the protocol for performing health checks on backend servers in the VPC
  channel.  
  The valid values are **TCP**, **HTTP** and **HTTPS**, defaults to **TCP**.

* `path` - (Optional, String) Specifies the destination path for health checks.  
  Required if the `protocol` is **HTTP** or **HTTPS**.

* `healthy_threshold` - (Optional, Int) Specifies the healthy threshold, which refers to the number of consecutive
  successful checks required for a backend server to be considered healthy.  
  The valid value ranges from `2` to `10`, defaults to `2`.

* `unhealthy_threshold` - (Optional, Int) Specifies the unhealthy threshold, which refers to the number of consecutive
  failed checks required for a backend server to be considered unhealthy.  
  The valid value ranges from `2` to `10`, defaults to `5`.

* `timeout` - (Optional, Int) Specifies the timeout for determining whether a health check fails, in second.  
  The value must be less than the value of the time `interval`.
  The valid value ranges from `2` to `30`, defaults to `5`.

* `interval` - (Optional, Int) Specifies the interval between consecutive checks, in second.  
  The valid value ranges from `5` to `300`, defaults to `10`.

* `http_code` - (Optional, String) Specifies the response codes for determining a successful HTTP response.  
  The valid value ranges from `100` to `599` and the valid formats are as follows:
  + The multiple values, for example, **200,201,202**.
  + The range, for example, **200-299**.
  + Both multiple values and ranges, for example, **201,202,210-299**.

  Required if the `protocol` is **HTTP**.

<a name="vpc_channel_members"></a>
The `members` block supports:

* `id` - (Optional, String) Specifies the ECS ID for each backend servers.
  Required if the `member_type` is **ECS**.
  This parameter and `ip_address` are alternative.

* `ip_address` - (Optional, String) Specifies the IP address each backend servers.
  Required if the `member_type` is **EIP**.

* `weight` - (Optional, Int) Specifies the backend server weight.
  The valid value ranges from `1` to `100`, defaults to `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPC channel.

* `created_at` - The time when the VPC channel was created.

* `status` - The current status of the VPC channel, supports **Normal** and **Abnormal**.

## Import

VPC Channels can be imported using their `name` and the ID of the related dedicated instance, separated by a slash, e.g.

```shell
$ terraform import huaweicloud_apig_vpc_channel.test <instance_id>/<name>
```
