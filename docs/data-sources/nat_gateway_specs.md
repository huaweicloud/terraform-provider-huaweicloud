---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateway_specs"
description: |-
  Use this data source to get NAT Gateway specs within HuaweiCloud.
---

# huaweicloud_nat_gateway_specs

Use this data source to get NAT Gateway specs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_nat_gateway_specs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `specs` - The list of the NAT Gateway specs.
  The valid values are as follows::
  + **0**: natgateway_xsmall.
  + **1**: natgateway_small.
  + **2**: natgateway_middle.
  + **3**: natgateway_large.
  + **4**: natgateway_xlarge.
  + **5**: natgateway_xxlarge.
  + **6**: public-nat.basic (traffic billing specifications).
