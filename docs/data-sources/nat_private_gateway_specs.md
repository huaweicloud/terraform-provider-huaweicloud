---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_gateway_specs"
description: |-
  Use this data source to get the list of NAT private gateway specifications within HuaweiCloud.
---

# huaweicloud_nat_private_gateway_specs

Use this data source to get the list of NAT private gateway specifications within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_nat_private_gateway_specs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

`id` - The data source ID.

* `specs` - The specification list.

  The [specs](#specs_struct) structure is documented below.

<a name="specs_struct"></a>
The `specs` block supports:

* `name` - The specification name.

* `code` - The specification code.

* `cbc_code` - The specification code on Cloud Business Center (CBC).

* `rule_max` - The maximum number of rules.

* `sess_max` - The maximum number of connections.

* `bps_max` - The maximum bandwidth in bit/s.

* `pps_max` - The maximum PPS.

* `qps_max` - The maximum QPS.
