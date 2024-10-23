---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_quotas"
description: |-
  Use this data source to get a list of the DC resource quotas.
---

# huaweicloud_dc_quotas

Use this data source to get a list of the DC resource quotas.

## Example Usage

```hcl
variable type {}

data "huaweicloud_dc_quotas" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, List) Specifies the quota type.
  The valid values are as follows:
  + **physicalConnect**: Quotas and usage for the connection.
  + **virtualInterface**: Quotas and usage for the virtual interface.
  + **connectGateway**: Quotas and usage for the connection gateway.
  + **geip**: Quotas and usage for the GEIP that each tenant can be associated.
  + **globalDcGateway**: Quotas and usage for the global DC gateway.
  + **peerLinkPerGdgw**: Quotas and usage for the peer links established with a global DC gateway.
  + **localGateway**: Quotas and usage for the local gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of the DC resource quotas.

  The [quotas](#quotas_quotas_struct) structure is documented below.

<a name="quotas_quotas_struct"></a>
The `quotas` block supports:

* `type` - The quota type.

* `quota` - The number of available quotas. The value **-1** indicates no quota limit.

* `used` - The number of used quotas.
