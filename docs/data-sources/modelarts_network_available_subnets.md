---
subcategory: "ModelArts"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_network_available_subnets"
description: |-
  Use this data source to query the available subnets of a network within HuaweiCloud.
---

# huaweicloud_modelarts_network_available_subnets

Use this data source to query the available subnets of a network within HuaweiCloud.

-> When other resources refer to this data source, any operations other than the creation operation  
  should use `ignore_changes` to ignore the reference to the results of this data source.

## Example Usage

```hcl
variable "network_name" {}
variable "subnet_id" {}

data "huaweicloud_modelarts_network_available_subnets" "test" {
  network_name = var.network_name
  subnet_id    = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

* `network_name` - (Required, String) Specifies the network ID.

* `subnet_id` - (Required, String) Specifies the subnet ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The subnet name.

* `network_id` - The subnet ID.

* `subnets` - The number of available network IP addresses of the subnet.  
  The [subnets](#modelarts_network_available_subnets_struct) structure is documented below.

<a name="modelarts_network_available_subnets_struct"></a>
The `subnets` block supports:

* `cidr` - The CIDR of the subnet.

* `ip_version` - The IP address type.

* `used_ips` - The number of used IP addresses.

* `total_ips` - The total number of IP addresses in the subnet.
