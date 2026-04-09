---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_vpcs_protection"
description: |-
  Use this data source to get the VPC protection information within HuaweiCloud.
---

# huaweicloud_cfw_vpcs_protection

Use this data source to get the VPC protection information within HuaweiCloud.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_vpcs_protection" "test" { 
  object_id = var.object_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The data of VPC protection information.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `other_protect_vpcs` - The list of other protected VPCs.

  The [other_protect_vpcs](#other_protect_vpcs_struct) structure is documented below.

* `other_total` - The total number of other protected VPCs.

* `protect_vpcs` - The list of protected VPCs.

  The [protect_vpcs](#protect_vpcs_struct) structure is documented below.

* `self_protect_vpcs` - The list of self-protected VPCs.

  The [self_protect_vpcs](#self_protect_vpcs_struct) structure is documented below.

* `self_total` - The total number of self-protected VPCs.

* `total` - The total number of protected VPCs.

* `total_assets` - The total number of assets.

<a name="other_protect_vpcs_struct"></a>
The `other_protect_vpcs` block supports:

* `vpc_id` - The VPC ID.

<a name="protect_vpcs_struct"></a>
The `protect_vpcs` block supports:

* `vpc_id` - The VPC ID.

<a name="self_protect_vpcs_struct"></a>
The `self_protect_vpcs` block supports:

* `vpc_id` - The VPC ID.
