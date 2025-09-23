---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_flavors"
description: |-
  Use this data source to get the available ELB Flavors.
---

# huaweicloud_elb_flavors

Use this data source to get the available ELB Flavors.

## Example Usage

```hcl
data "huaweicloud_elb_flavors" "flavors" {
  type            = "L7"
  max_connections = 200000
  cps             = 2000
  bandwidth       = 50
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors. If omitted, the provider-level region will be
  used.

* `flavor_id` - (Optional, String) Specifies the flavor ID.

* `type` - (Optional, String) Specifies the flavor type. Values options:
  + **L4**: indicates Layer-4 flavor.
  + **L7**: indicates Layer-7 flavor.
  + **L4_elastic**: indicates minimum Layer-4 flavor for elastic scaling.
  + **L7_elastic**: indicates minimum Layer-7 flavor for elastic scaling.
  + **L4_elastic_max**: indicates maximum Layer-4 flavor for elastic scaling.
  + **L7_elastic_max**: indicates maximum Layer-7 flavor for elastic scaling

* `name` - (Optional, String) Specifies the flavor name.

* `shared` - (Optional, String) Specifies whether the flavor is available to all users. Value options:
  + **true**: indicates that the flavor is available to all users.
  + **false**: indicates that the flavor is available only to a specific user.

* `public_border_group` - (Optional, String) Specifies the public border group.

* `category` - (Optional, Int) Specifies the category.

* `flavor_sold_out` - (Optional, String) Specifies whether the flavor is available.
  + **true**: indicates the flavor is unavailable.
  + **false**: indicates the flavor is available.

* `list_all` - (Optional, String) Specifies whether return all maximum elastic specifications. Value options: **true**,
  **false**.
  + If it is set to **true**, all maximum elastic specifications defined by l4_elastic_max and l7_elastic_max are returned.
  + If it is set to **false**, only the largest elastic specifications will be returned.
  + For Layer 4 load balancers, the specification with the highest cps value is returned. If the cps values are the same,
    the specification with the highest bandwidth value is returned.
  + For Layer 7 load balancers, the specification with highest https_cps value is returned. If the https_cps values are
    the same, the specification with highest qps value is returned.

* `max_connections` - (Optional, Int) Specifies the maximum connections in the flavor.

* `bandwidth` - (Optional, Int) Specifies the bandwidth size(Mbit/s) in the flavor.

* `cps` - (Optional, Int) Specifies the cps in the flavor.

* `qps` - (Optional, Int) Specifies the qps in the L7 flavor.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ids` - Indicates the list of flavor IDs.

* `flavors` - Indicates the list of flavors.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the flavor.

* `name` - Indicates the name of the flavor.

* `type` - Indicates the type of the flavor.

* `shared` - Indicates  whether the flavor is available to all users.

* `flavor_sold_out` - Indicates whether the flavor is available.

* `public_border_group` - Indicates the public border group.

* `category` - Indicates the category.

* `max_connections` - Indicates the maximum connections of the flavor.

* `cps` - Indicates the cps of the flavor.

* `qps` - Indicates the qps of the L7 flavor.

* `bandwidth` - Indicates the bandwidth size(Mbit/s) of the flavor.

* `lcu` - Indicates the number of LCUs in the flavor.

* `https_cps` - Indicates the number of new HTTPS connections.
