---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_endpoint_upgrade"
description: -|
  Use this resource to upgrade a basic VPC endpoint to a professional VPC endpoint within HuaweiCloud.
---

# huaweicloud_vpcep_endpoint_upgrade

Use this resource to upgrade a basic VPC endpoint to a professional VPC endpoint within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_endpoint_id" {}

resource "huaweicloud_vpcep_endpoint_upgrade" "test" {
  vpc_endpoint_id = var.vpc_endpoint_id
  action          = "start"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which the resource belongs to.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `vpc_endpoint_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC endpoint.

* `action` - (Required, String, NonUpdatable) Specifies the upgrade action. The default value is **start**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
