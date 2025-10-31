---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_gateway_specs"
description: |-
  Use this data source to get NAT Private Gateway specs within HuaweiCloud.
---

# huaweicloud_nat_private_gateway_specs

Use this data source to get NAT Private Gateway specs within HuaweiCloud.

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

* `specs` - The spec list.

  The [specs](#specs_struct) structure is documented below.

<a name="specs_struct"></a>
The `specs` block supports:

* `name` - The spec name.
* `code` - The spec code.
  The values can be:
  + **1**: cbc type privatenat_small.
  + **2**: cbc type privatenat_medium.
  + **3**: cbc type privatenat_large.
  + **4**: cbc type privatenat_xlarge.
  + **5**: cbc type privatenat_xxlarge.
  + **6**: cbc type private-nat.basic (traffic billing specifications).
* `cbc_code` - The spec cbc code.
  The values can be:
  + **privatenat_small**: corresponding code 1, small-sized.
  + **privatenat_medium**: corresponding code 2, medium-sized.
  + **privatenat_large**: corresponding code 3, large-sized.
  + **privatenat_xlarge**: corresponding code 4, .
  + **privatenat_xxlarge**: corresponding code 5, xxlarge-sized.
  + **private-nat.basic**: corresponding code 6 (traffic billing specifications).
* `rule_max` - The maximum number of rules.
* `sess_max` - The maximum number of connections.
* `bps_max` - The maximum bits per second.
* `pps_max` - The maximum packets per second.
* `qps_max` - The maximum queries per second.
