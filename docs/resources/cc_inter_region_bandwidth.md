---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_inter_region_bandwidth"
description: ""
---

# huaweicloud_cc_inter_region_bandwidth

Manages an inter-region bandwidth resource of Cloud Connect within HuaweiCloud.  

-> If network instances are in the same region, they can communicate with each other by default after they are loaded
   to one cloud connection. If network instances are in different regions, you need to assign inter-region bandwidths
   to ensure normal network communications between the instances.
   By default, the system allocates 10 kbit/s of bandwidth for testing network connectivity across regions.

## Example Usage

```hcl
variable "cloud_connection_id" {}
variable "bandwidth_package_id" {}
variable "region_local" {}
variable "region_remote" {}

resource "huaweicloud_cc_inter_region_bandwidth" "test" {
  cloud_connection_id  = var.cloud_connection_id
  bandwidth_package_id = var.bandwidth_package_id
  bandwidth            = 5
  inter_region_ids     = [var.region_local, var.region_remote]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cloud_connection_id` - (Required, String, ForceNew) Cloud connection ID.

  Changing this parameter will create a new resource.

* `bandwidth_package_id` - (Required, String, ForceNew) Bandwidth package ID.

  Changing this parameter will create a new resource.

* `bandwidth` - (Required, Int) Inter-region bandwidth size.  

* `inter_region_ids` - (Required, List, ForceNew) Two regions to which bandwidth is allocated.  

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `inter_regions` - Details about regions of the inter-region bandwidth.
  The [inter_regions](#interRegionBandwidth_Inter_Regions) structure is documented below.

<a name="interRegionBandwidth_Inter_Regions"></a>
The `inter_regions` block supports:

* `id` - Inter-region bandwidth ID.

* `project_id` - Project ID of a region where the inter-region bandwidth is used.

* `local_region_id` - ID of the local region where the inter-region bandwidth is used.

* `remote_region_id` - ID of the remote region where the inter-region bandwidth is used.

## Import

The inter-region bandwidth can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_inter_region_bandwidth.test 0ce123456a00f2591fabc00385ff1234
```
