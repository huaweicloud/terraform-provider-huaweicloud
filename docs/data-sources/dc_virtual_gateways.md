---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_virtual_gateways"
description: ""
---

# huaweicloud_dc_virtual_gateways

Use this data source to get the list of DC virtual gateways.

## Example Usage

```hcl
data "huaweicloud_dc_virtual_gateways" "test" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `virtual_gateway_id` - (Optional, String) Specifies the ID of the virtual gateway.

* `name` - (Optional, String) Specifies the name of the virtual gateway.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC connected to the virtual gateway.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `virtual_gateways` - Indicates the virtual gateways list.
  The [virtual_gateways](#DC_virtual_gateways) structure is documented below.

<a name="DC_virtual_gateways"></a>
The `virtual_gateways` block supports:

* `id` - The virtual gateway ID.

* `vpc_id` - Indicates the ID of the VPC connected by the virtual gateway.

* `name` - Indicates the virtual gateway name.

* `type` - Indicates the virtual gateway type.

* `status` - Indicates the virtual gateway status.

* `asn` - Indicates the local BGP ASN of the virtual gateway.

* `local_ep_group` - Indicates the IPv4 subnets connected by the virtual gateway.

* `description` - Indicates the virtual gateway description.

* `enterprise_project_id` - Indicates the ID of the enterprise project that the virtual gateway belongs to.
