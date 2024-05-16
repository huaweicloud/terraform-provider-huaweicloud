---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_elb_loadbalancers"
description: |-
  Use this data source to get the list of loadbalancer the CSS supported.
---

# huaweicloud_css_elb_loadbalancers

Use this data source to get the list of loadbalancer the CSS supported.

## Example Usage

```hcl
variable "cluster_id" {}
variable "name" {}

data "huaweicloud_css_elb_loadbalancers" "test" {
  cluster_id = var.cluster_id
  name       = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CSS cluster.

* `loadbalancer_id` - (Optional, String) Specifies the ID of the loadbalancer.

* `name` - (Optional, String) Specifies the name of the loadbalancer.

* `protocol_id` - (Optional, String) Specifies the layer 7 protocol ID of the loadbalancer.

* `is_cross` - (Optional, Bool) Specifies whether to enable cross-VPC backend.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `loadbalancers` - The list of the loadbalancer.

  The [loadbalancers](#loadbalancers_struct) structure is documented below.

<a name="loadbalancers_struct"></a>
The `loadbalancers` block supports:

* `id` - The loadbalancer ID.

* `name` - The loadbalancer name.

* `l7_flavor_id` - The layer 7 protocol ID of the loadbalancer.

* `ip_target_enable` - Whether to enable cross-VPC backend.
