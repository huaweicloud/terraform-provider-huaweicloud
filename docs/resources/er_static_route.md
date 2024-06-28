---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_static_route"
description: ""
---

# huaweicloud_er_static_route

Manages a static route under the ER route table within HuaweiCloud.

## Example Usage

### Create a static route and cross the VPC

```hcl
variable "route_table_id" {}
variable "destination_vpc_cidr" {}
variable "source_vpc_attachment_id" {}

resource "huaweicloud_er_static_route" "test" {
  route_table_id = var.route_table_id
  destination    = var.destination_vpc_cidr
  attachment_id  = var.source_vpc_attachment_id
}
```

### Create a black hole route

```hcl
variable "route_table_id" {}
variable "destination_vpc_cidr" {}

resource "huaweicloud_er_static_route" "test" {
  route_table_id = var.route_table_id
  destination    = var.destination_vpc_cidr
  is_blackhole   = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the static route and related route table are
  located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `route_table_id` - (Required, String, ForceNew) Specifies the ID of the route table to which the static route
  belongs.  
  Changing this parameter will create a new resource.

* `destination` - (Required, String, ForceNew) Specifies the destination of the static route.  
  Changing this parameter will create a new resource.

* `attachment_id` - (Optional, String) Specifies the ID of the corresponding attachment.

* `is_blackhole` - (Optional, Bool) Specifies whether route is the black hole route, defaults to `false`.  
  + If the value is empty or `false`, the parameter `attachment_id` is required.
  + If the value is `true`, the parameter `attachment_id` must be empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `type` - The type of the static route.

* `status` - The current status of the static route.

* `created_at` - The creation time of the static route.

* `updated_at` - The latest update time of the static route.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

Static routes can be imported using the related `route_table_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_er_static_route.test <route_table_id>/<id>
```
