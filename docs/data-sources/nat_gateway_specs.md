---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateway_specs"
description: |-
  Use this data source to get NAT gateway specifications within HuaweiCloud.
---

# huaweicloud_nat_gateway_specs

Use this data source to get NAT gateway specifications within HuaweiCloud.

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

`id` - The data source ID.

* `specs` - The list of the NAT Gateway specifications.
  The valid values are as follows::
  + **1**: Small, which supports up to 10,000 SNAT connections.
  + **2**: Medium, which supports up to 50,000 SNAT connections.
  + **3**: Large, which supports up to 200,000 SNAT connections.
  + **4**: Extra large, which supports up to 1,000,000 SNAT connections.
  + **5**: Enterprise-class, which supports up to 10,000,000 SNAT connections.
