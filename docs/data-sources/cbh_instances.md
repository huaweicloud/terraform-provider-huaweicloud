---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_instances"
description: ""
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

* `version` - (Optional, String) Specifies the current version of the instance image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of CBH instance.
  The [instances](#CbhInstances_Instance) structure is documented below.

<a name="CbhInstances_Instance"></a>
The `instances` block supports:

* `id` - Indicates the ID of the instance.

* `public_ip_id` - Indicates the ID of the elastic IP.

* `public_ip` - Indicates the elastic IP address.

* `name` - Indicates the instance name.

* `private_ip` - Indicates the private IP address of the instance.

* `status` - Indicates the status of the instance.

* `vpc_id` - Indicates the ID of a VPC.

* `subnet_id` - Indicates the ID of a subnet.

* `security_group_id` - Indicates the ID of a security group.

* `flavor_id` - Indicates the specification of the instance.

* `availability_zone` - Indicates the availability zone name.

* `version` - Indicates the current version of the instance image.
