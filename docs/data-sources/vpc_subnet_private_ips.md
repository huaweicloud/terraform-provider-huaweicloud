---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet_private_ips"
description: |-
  Use this data source to get a list of private IPs within HuaweiCloud.
---

# huaweicloud_vpc_subnet_private_ips

Use this data source to get a list of private IPs within HuaweiCloud.

## Example Usage

```hcl
variable "subnet_id" {}

data "huaweicloud_vpc_subnet_private_ips" "test" {
  subnet_id = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `subnet_id` - (Required, String) Specifies the ID of the subnet that the private IP address belongs to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `private_ips` - The list of private IP addresses.

  The [private_ips](#private_ips_struct) structure is documented below.

<a name="private_ips_struct"></a>
The `private_ips` block supports:

* `device_owner` - The resource using the private IP address. The parameter is left blank if it is not used.

* `ip_address` - The private IP address.

* `status` - The status of the private IP address.
  Possible values are **ACTIVE** and **DOWN**.

* `id` - The private IP address ID

* `subnet_id` - The ID of the subnet from which IP addresses are assigned.
