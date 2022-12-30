---
subcategory: "Cloud Bastion Host (CBH)"
---

# huaweicloud_cbh_instance

Manages CBH instance resources within HuaweiCloud.

## Example Usage

```HCL
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "availability_zone" {}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = var.flavor_id
  name              = "cbh_instance_test"
  vpc_id            = var.vpc_id
  availability_zone = var.availability_zone
  hx_password       = "test_123456",
  bastion_type      = "OEM"
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `flavor_id` - (Required, String, ForceNew) Specifies the product ID of the CBH server.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the CBH instance.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String) Specifies the ID of a subnet.

* `security_group_id` - (Required, String) Specifies the ID list of the security group.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone name.

  Changing this parameter will create a new resource.

* `hx_password` - (Required, String, ForceNew) Specifies the front end login password.

  Changing this parameter will create a new resource.

* `bastion_type` - (Required, String, ForceNew) Specifies the type of the bastion.

  Changing this parameter will create a new resource.

* `cloud_service_type` - (Required, String, ForceNew) Specifies the type of the CBH service.

  Changing this parameter will create a new resource.

* `charging_mode` - (Required, String, ForceNew) Specifies the charging mode of the read replica instance.
  The options are as follows:
  + **prePaid**: the yearly/monthly billing mode.

  Changing this parameter will create a new resource.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.

  Changing this parameter will create a new resource.

* `period` - (Required, String, ForceNew) Specifies the charging period of the read replica instance.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `auto_renew` - (Required, String, ForceNew) Specifies whether auto renew is enabled.
  Valid values are "true" and "false". Defaults to **false**.

  Changing this parameter will create a new resource.

* `image_id` - (Optional, String, ForceNew) Specifies a image ID.

  Changing this parameter will create a new resource.

* `user_data` - (Optional, String, ForceNew) Specifies the inject user data.

  Changing this parameter will create a new resource.

* `password` - (Optional, String) Specifies the initial password.

* `key_name` - (Optional, String, ForceNew) Specifies the secret key of the admin.

  Changing this parameter will create a new resource.

* `vpc_id` - (Optional, String, ForceNew) Specifies the ID of a VPC.

  Changing this parameter will create a new resource.

* `ip_address` - (Optional, String) Specifies the IP address of the subnet.

* `public_ip` - (Optional, List)
  The [PublicIP](#CBHInstance_PublicIP) structure is documented below.

* `root_volume` - (Optional, List, ForceNew)

  Changing this parameter will create a new resource.
  The [RootVolume](#CBHInstance_RootVolume) structure is documented below.

* `data_volume` - (Optional, List, ForceNew)

  Changing this parameter will create a new resource.
  The [DataVolume](#CBHInstance_DataVolume) structure is documented below.

* `slave_availability_zone` - (Optional, String, ForceNew) Specifies the slave availability zone name.
  The slave machine will be created when this field is not empty.

  Changing this parameter will create a new resource.

* `metadata` - (Optional, String, ForceNew) Specifies the metadata of the service.

  Changing this parameter will create a new resource.

* `resource_spec_code` - (Optional, String, ForceNew) Specifies the resource specification.

  Changing this parameter will create a new resource.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether the IPv6 network is enabled.

  Changing this parameter will create a new resource.

* `end_time` - (Optional, String, ForceNew) Specifies the end time.

  Changing this parameter will create a new resource.

* `relative_resource_id` - (Optional, String, ForceNew) Specifies the new capacity expansion.

  Changing this parameter will create a new resource.

* `product_info` - (Optional, List, ForceNew)

  Changing this parameter will create a new resource.
  The [ProductInfo](#CBHInstance_ProductInfo) structure is documented below.

* `network_type` - (Optional, String) Specifies the type of the network operation.
  The options are as follows: **create**, **renewals** and **change**.

<a name="CBHInstance_PublicIP"></a>
The `PublicIP` block supports:

* `id` - (Optional, String) Specifies the ID of the elastic IP.

* `address` - (Optional, String) Specifies the elastic IP address.

* `eip` - (Optional, List)
  The [Eip](#CBHInstance_Eip) structure is documented below.

<a name="CBHInstance_Eip"></a>
The `PublicIPEip` block supports:

* `type` - (Optional, String) Specifies the type of EIP.

* `flavor_id` - (Optional, String) Specifies the product ID of the IP associated with.

* `bandwidth` - (Optional, List)
  The [Bandwidth](#CBHInstance_Bandwidth) structure is documented below.

<a name="CBHInstance_Bandwidth"></a>
The `EipBandwidth` block supports:

* `size` - (Optional, String) Specifies the size of the bandwidth.

* `share_type` - (Optional, String) Specifies the share type. Only PER is supported noe.

* `charge_mode` - (Optional, String) Specifies the charge type. The value can be traffic or empty.

* `flavor_id` - (Optional, String) Specifies the product ID of the bandwidth associated with.

<a name="CBHInstance_RootVolume"></a>
The `RootVolume` block supports:

* `type` - (Optional, String) Specifies the type of volume.

* `size` - (Optional, String) Specifies the size of the root volume, unit is GB.

* `extend_param` - (Optional, String) Specifies the info of the volume.

<a name="CBHInstance_DataVolume"></a>
The `DataVolume` block supports:

* `type` - (Optional, String) Specifies the type of volume.

* `size` - (Optional, String) Specifies the size of the data volume, unit is GB.

* `extend_param` - (Optional, String) Specifies the info of the volume.

<a name="CBHInstance_ProductInfo"></a>
The `ProductInfo` block supports:

* `product_id` - (Optional, String) Specifies the ID of the product.

* `cloud_service_type` - (Optional, String) Specifies the type of the CBH service.

* `resource_type` - (Optional, String) Specifies the type of the CBH resource.

* `resource_spec_code` - (Optional, String) Specifies the flavor of the CBH service.

* `availability_zone` - (Optional, String) Specifies the availability zone name.

* `resource_size_measure_id` - (Optional, String) Specifies the resource capacity measurement ID.

* `resource_size` - (Optional, String) Specifies the size of the resource capacity.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `publicip_id` - Indicates the ID of the elastic IP.

* `exp_time` - Indicates the expire time of the instance.

* `start_time` - Indicates the start time of the instance.

* `release_time` - Indicates the release time of the instance.

* `instance_id` - Indicates the server id of the instance.

* `private_ip` - Indicates the private ip of the instance.

* `task_status` - Indicates the task status of the instance.

* `status` - Indicates the status of the instance.

* `update` - Indicates whether the instance image can be upgraded.

* `instance_key` - Indicates the ID of the instance.

* `resource_id` - Indicates the ID of the resource.

* `alter_permit` - Indicates whether the front-end displays the capacity expansion button.

* `bastion_version` - Indicates the current version of the instance image.

* `new_bastion_version` - Indicates the latest version of the instance image.

* `instance_status` - Indicates the status of the instance.

* `description` - Indicates the type of the bastion.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `delete` - Default is 10 minute.

## Import

The cbh instance can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_cbh_instance.test 0ce123456a00f2591fabc00385ff1234
```
