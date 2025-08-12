---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_connect_gateway_geips"
description: |-
  Use this data source to get the list of global EIPs bound to a connect gateway
---

# huaweicloud_dc_connect_gateway_geips

Use this data source to get the list of global EIPs bound to a connect gateway

## Example Usage

```hcl
variable "connect_gateway_id"{}

data "huaweicloud_dc_connect_gateway_geips" "test"{
  connect_gateway_id = var.connect_gateway_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `connect_gateway_id` - (Required, String) Specifies the DC connect gateway ID.

* `global_eip_id` - (Optional, List) Specifies the global EIP ID.

* `global_eip_segment_id` - (Optional, List) Specifies the ID of the global EIP range.

* `status` - (Optional, List) Specifies the status by which instances are queried.

* `sort_key` - (Optional, String) Specifies the sorting field.

* `sort_dir` - (Optional, List) Specifies the sorting order of returned results.
  Value options: **asc (default)** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `global_eips` - Indicates the list of the bound global EIPs.

  The [global_eips](#global_eips_struct) structure is documented below.

<a name="global_eips_struct"></a>
The `global_eips` block supports:

* `global_eip_id` - Indicates the global EIP ID.

* `global_eip_segment_id` - Indicates the ID of the global EIP range.

* `type` - Indicates the global EIP type.
  The value can be **IP_ADDRESS** or **IP_SEGMENT**.

* `error_message` - Indicates the cause of the failure to bind the global EIP.

* `gcb_id` - Indicates the global connection bandwidth ID.

* `status` - Indicates whether the global EIP has been bound.

* `cidr` - Indicates the global EIP and its subnet mask.

* `address_family` - Indicates the address family of the global EIP.

* `ie_vtep_ip` - Indicates the VTEP IP address of the CloudPond cluster.

* `created_time` - Indicates the time when the global EIP was bound.
