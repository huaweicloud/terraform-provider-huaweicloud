---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_site_network_capabilities"
description: |-
  Use this data source to get the list of CC site network capabilities.
---

# huaweicloud_cc_site_network_capabilities

Use this data source to get the list of CC site network capabilities.

## Example Usage

```hcl
data "huaweicloud_cc_site_network_capabilities" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `specification` - (Optional, List) Specifies the site network capabilities. Multiple capabilities can be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `capabilities` - Indicates the list of site network capabilities.

  The [capabilities](#capabilities_struct) structure is documented below.

<a name="capabilities_struct"></a>
The `capabilities` block supports:

* `id` - Indicates the instance ID.

* `specification` - Indicates the site network capabilities.
  The value can be:
  + **site-network.is-support**: site networks
  + **site-network.is-support-enterprise-project**: enterprise projects for site networks
  + **site-network.is-support-tag**: site network tagging
  + **site-network.is-support-intra-region**: site networks in the same region
  + **site-network.support-topologies**: site network topologies
  + **site-network.support-regions**: list of the regions that support site networks
  + **site-network.support-dscp-regions**: list of the regions where DSCP is supported on a site network
  + **site-network.support-freeze-regions**: list of the regions where site networks can be frozen
  + **site-network.support-locations**: list of site access points
  + **site-connection-bandwidth.size-range**: capacity of site-to-site connection bandwidth
  + **site-connection-bandwidth.charge-mode**: billing options of bandwidth used by a site-to-site connection
  + **site-connection-bandwidth.free-line**: free lines for cross-site connections

* `is_support_enterprise_project` - Indicates whether enterprise projects are supported for site networks.

* `is_support_tag` - Indicates whether site network tagging is supported.

* `is_support_intra_region` - Indicates whether site networks in the same region can be created.

* `is_support` - Indicates whether site networks are supported.

* `support_locations` - Indicates the support locations list of a site network.

* `support_regions` - Indicates the region list of a site network.

* `support_freeze_regions` - Indicates the freeze regions list of a site network.

* `support_topologies` - Indicates the topology list of a site network.

* `support_dscp_regions` - Indicates the dscp regions list of a site network.

* `charge_mode` - Indicates the charge mode list of a site network.

* `size_range` - Indicates the size_range.

  The [size_range](#capabilities_size_range_struct) structure is documented below.

<a name="capabilities_size_range_struct"></a>
The `size_range` block supports:

* `min` - Indicates the minimum value.

* `max` - Indicates the maximum value.
