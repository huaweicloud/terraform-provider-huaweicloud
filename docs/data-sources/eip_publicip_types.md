---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eip_publicip_types"
description: |-
  Use this data source to get a list of EIP public IP types.
---

# huaweicloud_eip_publicip_types

Use this data source to get a list of EIP public IP types.

## Example Usage

```hcl
data "huaweicloud_eip_publicip_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `publicip_types` - Indicates the public IP types.

  The [publicip_types](#publicip_types_struct) structure is documented below.

<a name="publicip_types_struct"></a>
The `publicip_types` block supports:

* `id` - Indicates the elastic public IP pool ID.

* `type` - Indicates the elastic public IP pool type. Valid values are:
  + **5_bgp**: Dynamic BGP public IP pool.
  + **5_sbgp**: Static BGP public IP pool.
  + **5_testbgp**: IP test pool.
