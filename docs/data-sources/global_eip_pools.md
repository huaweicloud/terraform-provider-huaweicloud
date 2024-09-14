---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_pools"
description: ""
---

# huaweicloud_global_eip_pools

Use this data source to get a list of global EIP pools available to your account.

## Example Usage

### Get all available global EIP pools

```hcl
data "huaweicloud_global_eip_pools" "all" {}
```

### Get specific available global EIP pools information

```hcl
data "huaweicloud_global_eip_pools" "test" {
  access_site = "cn-south-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `access_site` - (Optional, String) Specifies the access site to which the global EIP pool belongs.

* `ip_version` - (Optional, Int) Specifies the ip version. Valid values are `4` and `6`.

* `isp` - (Optional, String) Specifies the internet service provider of the global EIP pool.

* `name` - (Optional, String) Specifies the name of the global EIP pool.

* `type` - (Optional, String) Specifies the type of the global EIP pool.

  Valid values are:
  + **GEIP**: global EIP.
  + **GEIP_SEGMENT**: global EIP range.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `geip_pools` - The available global EIP pools list.
  The [geip_pools](#attrblock--geip_pools) structure is documented below.

<a name="attrblock--geip_pools"></a>
The `geip_pools` block supports:

* `id` - The id of the global EIP pool.

* `access_site` - The access site to which the global EIP pool belongs.

* `cn_name` - The Chinese name of the global EIP pool.

* `en_name` - The English name of the global EIP pool.

* `ip_version` - The ip version of the global EIP pool.

* `isp` - The internet service provider of the global EIP pool.

* `name` - The name of the global EIP pool.

* `type` - The type of the global EIP pool.

* `allowed_bandwidth_types` - The allowed bandwidth type of the global EIP pool
  The [allowed_bandwidth_type](#attrblock--geip_pools--allowed_bandwidth_type) structure is documented below.

* `created_at` - The create time of the global EIP pool.

* `updated_at` - The update time of the global EIP pool.

<a name="attrblock--geip_pools--allowed_bandwidth_type"></a>
The `allowed_bandwidth_type` block supports:

* `cn_name` - The Chinese name of the allowed bandwidth type.

* `en_name` - The English name of the allowed bandwidth type.

* `type` - The type of the allowed bandwidth.
