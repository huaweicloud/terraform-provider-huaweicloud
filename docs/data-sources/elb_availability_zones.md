---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_availability_zones"
description: |-
  Use this data source to get the list of available AZs when create a load balancer.
---

# huaweicloud_elb_availability_zones

Use this data source to get the list of available AZs when create a load balancer.

## Example Usage

```hcl
data "huaweicloud_elb_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `public_border_group` - (Optional, String) Specifies the public border group.

* `loadbalancer_id` - (Optional, String) Specifies the load balancer ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - Indicates the AZs that are available during load balancer creation. For example, in [az1,az2] and
  [az2,az3] sets, you can select **az1** and **az2** or **az2** and **az3**, but cannot select **az1** and **az3**.
  The [availability_zones](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `list` - Indicates the AZs list.
  The [list](#list_struct) structure is documented below.

<a name="list_struct"></a>
The `list` block supports:

* `code` - Indicates the AZ code.

* `state` - Indicates the AZ status. The value can only be **ACTIVE**.

* `protocol` - Indicates the type of the flavor that is not sold out. The value can be:
  + **L4**: indicates the flavor at Layer 4 (flavor for network load balancing).
  + **L7**: indicates the flavor at Layer 7 (flavor for application load balancing).

* `public_border_group` - Indicates the public border group, for example, **center**.

* `category` - Indicates the AZ code. The value can be:
  + **0**: indicates center.
  + **21**: indicates homezone.
