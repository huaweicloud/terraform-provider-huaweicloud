---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_availability_zones"
description: |-
  Use this data source to get the list of DMS rocketMQ availability zones within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_availability_zones

Use this data source to get the list of DMS rocketMQ availability zones within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dms_rocketmq_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the availability zones are located.  
  If omitted, the provider-level region will be used.

## Attributes Reference

The following attributes are exported:

* `availability_zones` - The list of availability zones.  
  The [availability_zones](#rocketmq_availability_zones_attr) structure is documented below.

<a name="rocketmq_availability_zones_attr"></a>
The `availability_zones` block supports:

* `id` - The ID of the availability zone.

* `name` - The name of the availability zone. e.g. `AZ1`.

* `code` - The code of the availability zone. e.g. `cn-north-4a`.

* `port` - The port of the availability zone.

* `sold_out` - Whether the availability zone is sold out.

* `resource_availability` - The resource availability of the availability zone.

* `default_az` - Whether the availability zone is the default availability zone.

* `remain_time` - The remaining time of the availability zone.

* `ipv6_enable` - Whether the availability zone supports IPv6.
