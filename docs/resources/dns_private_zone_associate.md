---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_private_zone_associate"
description: |-
  Manages a DNS private zone associate resource within HuaweiCloud.
---


# huaweicloud_dns_private_zone_associate

Manages a DNS private zone associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "zone_id" {}
variable "router_id" {}

resource "huaweicloud_dns_private_zone_associate" "test" {
  zone_id   = var.zone_id
  router_id = var.router_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `zone_id` - (Required, String, NonUpdatable) The ID of the zone to which the record set belongs.
  Changing this creates a new resource.

* `router_id` - (Required, String, NonUpdatable) The ID of the associated VPC.

* `router_region` - (Optional, String, NonUpdatable) The region of the VPC. Default to the region of resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the associated VPC.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The associated relationship of private zone and VPC can be imported using `zone_id` and `router_id`, e.g.

```bash
$ terraform import huaweicloud_dns_private_zone_associate.test <zone_id>/<router_id>
```
