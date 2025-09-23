---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_policy_black_white_lists"
description: |-
  Use this data source to get the list of Advanced Anti-DDos policy black white list within HuaweiCloud.
---

# huaweicloud_aad_policy_black_white_lists

Use this data source to get the list of Advanced Anti-DDos policy black white list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_policy_black_white_lists" "test" {}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name.

* `overseas_type` - (Required, Int) Specifies protection zone.  
  The valid values are as follows:
  + **0**: Mainland.
  + **1**: Overseas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `black` - The black list detail.  
  The [black](#black_struct) structure is documented below.

* `white` - The white list detail.  
  The [white](#white_struct) structure is documented below.

<a name="black_struct"></a>
The `black` block supports:

* `id` - The ID.

* `type` - The type. `0` indicates blacklist, and `1` indicates whitelist.

* `ip` - The IP address.

* `domain_id` - The domain ID.

<a name="white_struct"></a>
The `white` block supports:

* `id` - The ID.

* `type` - The type. `0` indicates blacklist, and `1` indicates whitelist.

* `ip` - The IP address.

* `domain_id` - The domain ID.
