---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_routes"
description: ""
---

# huaweicloud_apig_instance_routes

Using this resource to manage the instance routes within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_apig_instance_routes" "test" {
  instance_id = var.instance_id
  nexthops    = ["172.16.3.0/24", "172.16.7.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dedicated instance and routes are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the routes belong.
  Changing this will create a new resource.

* `nexthops` - (Required, List) Specifies the configuration of the next-hop routes.

-> The network segment of the next hop cannot overlap with the network segment of the APIG instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (instance ID).

## Import

Routes can be imported using their related dedicated instance ID (`instance_id`), e.g.

```bash
$ terraform import huaweicloud_apig_instance_routes.test 128001b3c5eb4d3e91a8da9c0f46420f
```
