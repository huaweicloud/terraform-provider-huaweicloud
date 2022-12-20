---
subcategory: "Cloud Bastion Host (CBH)"
---

# huaweicloud_cbh_instances

Use this data source to get the list of CBH instance.

## Example Usage

```hcl
data "huaweicloud_cbh_instance" "test" {
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

* `total` - Indicates the total of instances.

* `quota_detail` -
  The [QuotaDetail](#CbhInstances_QuotaDetail) structure is documented below.

* `instance` - Indicates the list of CBH instance.
  The [Instance](#CbhInstances_Instance) structure is documented below.

<a name="CbhInstances_QuotaDetail"></a>
The `QuotaDetail` block supports:

* `zh_cn` - Indicates the Chinese quota description.

* `en_us` - Indicates the English quota description.

* `remaining` - Indicates the tenant remaining quota quantity.

<a name="CbhInstances_Instance"></a>
The `Instance` block supports:

* `publicip` - Indicates the public ip of the instance.

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

* `specification` - Indicates the specification of the instance.

* `update` - Indicates whether the instance image can be upgraded.

* `instance_key` - Indicates the ID of the instance.

* `order_id` - Indicates the ID of order.

* `period_num` - Indicates the duration of tenant purchase.

* `bastion_type` - Indicates the type of the bastion.

* `public_id` - Indicates the ID of the elastic IP bound by the instance.

* `alter_permit` - Indicates whether the front-end displays the capacity expansion button.

* `bastion_version` - Indicates the current version of the instance image.

* `new_bastion_version` - Indicates the latest version of the instance image.

* `instance_status` - Indicates the status of the instance.

* `description` - Indicates the type of the bastion.

* `auto_renew` - Indicates whether auto renew is enabled.
