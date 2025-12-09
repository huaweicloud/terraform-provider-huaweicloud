---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_upgrade"
description: -|
  Manages a VPC endpoint service upgrade resource within HuaweiCloud.
---

# huaweicloud_vpcep_service_upgrade

Manages a VPC endpoint service upgrade resource within HuaweiCloud.

## Example Usage

```hcl
variable "service_id" {}

resource "huaweicloud_vpcep_service_upgrade" "test" {
  service_id = var.service_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `service_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC endpoint service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
