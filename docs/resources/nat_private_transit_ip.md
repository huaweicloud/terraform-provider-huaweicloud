---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_transit_ip"
description: |-
  Manages a transit IP resource of the **private** NAT within HuaweiCloud.
---

# huaweicloud_nat_private_transit_ip

Manages a transit IP resource of the **private** NAT within HuaweiCloud.

## Example Usage

```hcl
variable "subnet_id" {}
variable "ipv4_address" {}
variable "enterprise_project_id" {}

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id             = var.subnet_id
  ip_address            = var.ipv4_address
  enterprise_project_id = var.enterprise_project_id

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the transit IP is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the transit subnet ID to which the transit IP belongs.  
  Changing this will create a new resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address of the transit subnet.  
  Changing this will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the transit
  IP belongs.  
  Changing this will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the transit IP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `network_interface_id` - The network interface ID of the transit IP for private NAT.

* `gateway_id` - The ID of the private NAT gateway to which the transit IP belongs.

* `status` - The status of the transit IP.

* `created_at` - The creation time of the transit IP for private NAT.

* `updated_at` - The latest update time of the transit IP for private NAT.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_nat_private_transit_ip.test <id>
```
