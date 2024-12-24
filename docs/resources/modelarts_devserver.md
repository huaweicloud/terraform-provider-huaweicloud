---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_devserver"
description: |-
  Manages a ModelArts DevServer resource within HuaweiCloud.
---

# huaweicloud_modelarts_devserver

Manages a ModelArts DevServer resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_name" {}
variable "server_flavor" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "image_id" {}
variable "admin_pass" {}

resource "huaweicloud_modelarts_devserver" "test" {
  name              = var.server_name
  flavor            = var.server_flavor
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  image_id          = var.image_id
  admin_pass        = var.admin_pass

  root_volume {
    size = 100
    type = "SSD"
  }

  charging_mode = "PRE_PAID"
  period        = 1
  period_unit   = "MONTH"
  auto_renew    = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DevServer.
  Changing this creates a new resource.  
  The name valid length is limited from `1` to `64`, only English letters, digits, underscores (_) and hyphens (-) are
  allowed.

* `flavor` - (Required, String, ForceNew) Specifies the flavor of the DevServer.
  Changing this creates a new resource.  
  For the flavor, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html)
  to submit a service ticket to apply for the flavor.

* `architecture` - (Optional, String, ForceNew) Specifies the architecture of the DevServer.
  Changing this creates a new resource.  
  This parameter value is related to the `flavor` parameter.  
  The valid values are as follows:
  + **X86**
  + **ARM**
  
* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the DevServer belongs.
  Changing this creates a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet to which the DevServer belongs.
  Changing this creates a new resource.

* `security_group_id` - (Required, String, ForceNew) Specifies the ID of security group to which the DevServer belongs.
  Changing this creates a new resource.

* `image_id` - (Required, String, ForceNew) Specifies the image ID of the DevServer.
  Changing this creates a new resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the DevServer.
  Changing this creates a new resource.  
  This parameter value is related to the `flavor` parameter.  
  The valid values are as follows:
  + **BMS**
  + **ECS**

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone where the DevServer is located.
  Changing this creates a new resource.

* `admin_pass` - (Optional, String, ForceNew) Specifies the login password for logging in to the server.
  Changing this creates a new resource.  
  The password format must meet the following conditions:
  + Must be `8` to `26` characters.
  + The password must contain at least three types of the following characters: digit, uppercase letter, lowercase letter
    and special characters (!@%-_=+[{}]:,./?).
  + The password cannot be the username or the username spelled backwards
  + The password cannot contain root, administrator or their reverse order.

-> Exactly one of `admin_pass` and `key_pair_name` must be provided.

* `key_pair_name` - (Optional, String, ForceNew) Specifies the key pair name for logging in to the server.
  Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the DevServer belongs.
  Changing this creates a new resource.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.  

* `root_volume` - (Optional, List, ForceNew) Specifies the system disk configuration of the DevServer.
  Changing this creates a new resource.  
  This parameter is related to the `flavor` parameter.  
  The [root_volume](#devServer_root_volume) structure is documented below.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether to enable IPv6.
  Changing this creates a new resource.  
  This parameter is available only when the current subnet, flavor and image all support IPv6.

* `roce_id` - (Optional, String, ForceNew) Specifies the RoCE network ID of the DevServer.
  Changing this creates a new resource.  
  This parameter value is related to the `flavor` parameter.

* `user_data` - (Optional, String, ForceNew) Specifies the user data defined for the server.
  Changing this creates a new resource.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the DevServer.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **PRE_PAID**: The yearly/monthly billing mode.
  + **POST_PAID**: The pay-per-use billing mode.
  
* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the DevServer.
  Changing this creates a new resource.  
  This parameter is required and available if `charging_mode` is set to **PRE_PAID**.  
  The valid values are as follows:
  + **MONTH**
  + **YEAR**

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the DevServer.
  Changing this creates a new resource.  
  This parameter is required and available if `charging_mode` is set to **PRE_PAID**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Defaults to **false**.  
  This parameter is available if `charging_mode` is set to **PRE_PAID**.  
  The valid values are **true** and **false**.

<a name="devServer_root_volume"></a>
The `root_volume` block supports:

* `type` - (Optional, String, ForceNew) Specifies the type of system disk.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **ESSD**: Extreme SSD type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **SAS**: High I/O type.
  + **SATA**: Common I/O type.

* `size` - (Optional, Int, ForceNew) Specifies the size of system disk.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also DevServer ID.

* `created_at` - The creation time of the DevServer, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 30 minutes.

## Import

The DevServer resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_devserver.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `subnet_id`, `security_group_id`, `availability_zone`, `admin_pass`,
`enterprise_project_id`, `root_volume`, `ipv6_enable`, `roce_id`, `user_data`, `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_devserver" "test" {
  ...

  lifecycle {
    ignore_changes = [
      subnet_id, security_group_id, availability_zone, admin_pass, enterprise_project_id, root_volume, ipv6_enable,
      roce_id, user_data, period_unit, period, auto_renew,
    ]
  }
}
```
