---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_bandwidth_rule"
description: |-
  Manages a VPC EIP bandwidth rule resource within HuaweiCloud.
---

# huaweicloud_vpc_eip_bandwidth_rule

Manages a VPC EIP bandwidth rule resource within HuaweiCloud.

-> This resource is a one-time action resource used to create a bandwidth rule for shared bandwidth.
Please enable the enterprise-level QoS feature before creating the bandwidth rule.
Deleting this resource will not delete the actual bandwidth rule on the cloud, but will only remove the resource
information from the tfstate file.

## Example Usage

```hcl
variable "bandwidth_id" {}
variable "publicip_id" {}
variable "name" {}

resource "huaweicloud_vpc_eip_bandwidth_rule" "test" {
  bandwidth_id          = var.bandwidth_id 
  name                  = var.name 
  egress_size           = 20 
  egress_guarented_size = 10
  
  publicip_info { 
    publicip_id = var.publicip_id 
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC EIP associate resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `bandwidth_id` - (Required, String, NonUpdatable) Specifies the ID of the shared bandwidth to which the rule belongs.

* `name` - (Required, String, NonUpdatable) Specifies the name of the bandwidth rule.  
  The value can contain 1 to 64 characters, including letters, digits, Chinese characters, underscores (_), hyphens (-),
  and periods (.).

* `egress_size` - (Required, Int, NonUpdatable) Specifies the outbound bandwidth size (Mbit/s) for the bandwidth rule.

* `egress_guarented_size` - (Required, Int, NonUpdatable) Specifies the guarented outbound bandwidth size (Mbit/s) for the
  bandwidth rule. This value must be less than or equal to `egress_size`.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the bandwidth rule.

* `publicip_info` - (Required, List, NonUpdatable) Specifies the public IP information associated with the bandwidth rule.

The [publicip_info](#publicip_info_struct) structure is documented below.

<a name="publicip_info_struct"></a>
The `publicip_info` block supports:

* `publicip_id` - (Required, String, ForceNew) Specifies the ID of the public IP address to be associated with the
  bandwidth rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
