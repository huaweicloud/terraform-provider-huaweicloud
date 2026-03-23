---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_inspection_vpcs"
description: |-
  Use this data source to get the CFW east-west associated VPC.
---

# huaweicloud_cfw_inspection_vpcs

Use this data source to get the CFW east-west associated VPC.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_inspection_vpcs" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `inspection_vpc_list` - The list of VPCs using traffic redirection.

  The [inspection_vpc_list](#data_inspection_vpc_list_struct) structure is documented below.

<a name="data_inspection_vpc_list_struct"></a>
The `inspection_vpc_list` block supports:

* `inspection_vpc_id` - The VPC ID for referral.

* `name` - The VPC name for referral.
