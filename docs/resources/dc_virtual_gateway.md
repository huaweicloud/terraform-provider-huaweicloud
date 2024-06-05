---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_virtual_gateway"
description: ""
---

# huaweicloud_dc_virtual_gateway

Manages a virtual gateway resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "vpc_cidr" {}
variable "gateway_name" {}

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id = var.vpc_id
  name   = var.gateway_name

  local_ep_group = [
    var.vpc_cidr,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the virtual gateway is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC connected to the virtual gateway.
  Changing this will create a new resource.

* `local_ep_group` - (Required, List) Specifies the list of IPv4 subnets from the virtual gateway to access cloud
  services, which is usually the CIDR block of the VPC.

* `name` - (Required, String) Specifies the name of the virtual gateway.
  The valid length is limited from `3` to `64`, only chinese and english letters, digits, hyphens (-), underscores (_)
  and dots (.) are allowed.
  The Chinese characters must be in **UTF-8** or **Unicode** format.

* `description` - (Optional, String) Specifies the description of the virtual gateway.
  The description contain a maximum of 128 characters and the angle brackets (< and >) are not allowed.
  Chinese characters must be in **UTF-8** or **Unicode** format.

* `asn` - (Optional, Int, ForceNew) Specifies the local BGP ASN of the virtual gateway.
  The valid value is range from `1` to `4,294,967,295`.
  Changing this will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the virtual
  gateway belongs.
  Changing this will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the virtual gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the virtual gateway.

* `status` - The current status of the virtual gateway.

## Import

Virtual gateways can be imported using their `id`, e.g.

```shell
$ terraform import huaweicloud_dc_virtual_gateway.test f6f36e69-d980-4b0a-a33d-b9b125b3896c
```
