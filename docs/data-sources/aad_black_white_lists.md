---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_black_white_lists"
description: |-
  Use this data source to get the list of Advanced Anti-DDos black white list within HuaweiCloud.
---

# huaweicloud_aad_black_white_lists

Use this data source to get the list of Advanced Anti-DDos black white list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_black_white_lists" "test" {}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the AAD instance ID.

* `type` - (Required, String) Specifies the rule type. Valid values are **black** and **white**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ips` - The IP address list.

  The [ips](#ips_struct) structure is documented below.

<a name="ips_struct"></a>
The `ips` block supports:

* `desc` - The description.

* `ip` - The black white IP address.
