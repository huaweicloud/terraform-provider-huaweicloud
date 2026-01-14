---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_source_ip_list"
description: |-
  Use this data source to get the Advanced Anti-DDos source IP list.
---

# huaweicloud_aad_source_ip_list

Use this data source to get the Advanced Anti-DDos source IP list.

## Example Usage

```hcl
data "huaweicloud_aad_source_ip_list" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ips` - The IP list.

  The [ips](#ips_struct) structure is documented below.

<a name="ips_struct"></a>
The `ips` block supports:

* `data_center` - The data center.

* `isp` - The line.

* `ip` - The instance IP address.
