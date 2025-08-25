---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_flow_block"
description: |-
  Use this data source to get the list of Advanced Anti-DDos flow block information within HuaweiCloud.
---

# huaweicloud_aad_flow_block

Use this data source to get the list of Advanced Anti-DDos flow block information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_flow_block" "test" {}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `ips` - The ips list.  
  The [ips](#ips_struct) structure is documented below.

<a name="ips_struct"></a>
The `ips` block supports:

* `ip_id` - The IP ID.

* `ip` - The IP.

* `isp` - The isp.

* `data_center` - The data center.

* `foreign_switch_status` - The overseas region ban status. `0` represents closed, `1` represents open.

* `udp_switch_status` - The UDP protocol disabled. `0` represents closed, `1` represents open.
