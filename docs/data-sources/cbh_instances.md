---
subcategory: "Cloud Bastion Host (CBH)"
---

# huaweicloud_cbh_instances

Use this data source to get the list of CBH instance.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

data "huaweicloud_cbh_instances" "test" {
  name              = "test_name"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the instance name.

* `vpc_id` - (Optional, String) Specifies the ID of a VPC.

* `subnet_id` - (Optional, String) Specifies the ID of a subnet.

* `security_group_id` - (Optional, String) Specifies the ID of a security group.

* `flavor_id` - (Optional, String) Specifies the specification of the instance.

* `bastion_type` - (Optional, String) Specifies the type of the bastion.

* `bastion_version` - (Optional, String) Specifies the current version of the instance image

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `instances` - Indicates the list of CBH instance.
  The [Instance](#CbhInstances_Instance) structure is documented below.

<a name="CbhInstances_Instance"></a>
The `Instance` block supports:

* `publicip_id` - Indicates the ID of the elastic IP.

* `exp_time` - Indicates the expire time of the instance.

* `start_time` - Indicates the start time of the instance.

* `end_time` - Indicates the end time of the instance.

* `release_time` - Indicates the release time of the instance.

* `name` - Indicates the instance name.

* `instance_id` - Indicates the server id of the instance.

* `private_ip` - Indicates the private ip of the instance.

* `task_status` - Indicates the task status of the instance.

* `status` - Indicates the status of the instance.

* `vpc_id` - Indicates the ID of a VPC.

* `subnet_id` - Indicates the ID of a subnet.

* `security_group_id` - Indicates the ID of a security group.

* `flavor_id` - Indicates the specification of the instance.

* `update` - Indicates whether the instance image can be upgraded.

* `instance_key` - Indicates the ID of the instance.

* `resource_id` - Indicates the ID of the resource.

* `period` - Indicates the duration of tenant purchase.

* `bastion_type` - Indicates the type of the bastion.

* `alter_permit` - Indicates whether the front-end displays the capacity expansion button.

* `bastion_version` - Indicates the current version of the instance image.

* `new_bastion_version` - Indicates the latest version of the instance image.

* `instance_status` - Indicates the status of the instance.

* `description` - Indicates the type of the bastion.

* `auto_renew` - Indicates whether auto renew is enabled.
