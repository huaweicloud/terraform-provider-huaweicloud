---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_connection_update"
description: -|
  Manages a VPC endpoint service connection description update resource within HuaweiCloud.
---

# huaweicloud_vpcep_service_connection_update

Manages a VPC endpoint service connection description update resource within HuaweiCloud.

## Example Usage

```hcl
variable "service_id" {}
variable "endpoint_id" {}
variable "description" {}

resource "huaweicloud_vpcep_service_connection_update" "test" {
  service_id  = var.service_id
  endpoint_id = var.endpoint_id
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `service_id` - (Required, String, NonUpdatable) Specifies the VPC endpoint service ID.

* `endpoint_id` - (Required, String, NonUpdatable) Specifies the VPC endpoint ID.

* `description` - (Optional, String) Specifies the description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The VPC endpoint service connection with description can be imported using the `service_id` and `endpoint_id`, separated
by a slash, e.g.

```bash
$ terraform import huaweicloud_vpcep_service_connection_update <service_id>/<endpoint_id>
```
