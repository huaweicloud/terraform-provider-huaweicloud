---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_restriction"
description: |-
  Use this data source to query the restricted information of the dedicated instance within HuaweiCloud.
---

# huaweicloud_apig_instance_restriction

Use this data source to query the restricted information of the dedicated instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_restriction" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the dedicated instance is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `restrict_cidrs` - The list of restricted IP CIDR blocks.

* `resource_subnet_cidr` - The CIDR block of the resource subnet.
