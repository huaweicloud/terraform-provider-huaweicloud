---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_site_network_quotas"
description: |-
  Use this data source to get the site network quotas.
---

# huaweicloud_cc_site_network_quotas

Use this data source to get the site network quotas.

## Example Usage

```hcl
data "huaweicloud_cc_site_network_quotas" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `quota_type` - (Optional, List) Specifies the quota type. Multiple quota types can be queried.
  Value options:
  + **site_networks_per_account**: the maximum number of site networks for each account
  + **sites_per_mesh_site_network**: the maximum number of sites on a site network of the mesh type
  + **spoke_sites_per_star_site_network**: the maximum number of spoke sites on a site network of the star type
  + **sites_per_hybrid_site_network**: the maximum number of sites on a hybrid site network

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the quota list of a site network.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `quota_key` - Indicates the quota type of a site network.

* `quota_limit` - Indicates the quotas.

* `used` - Indicates the used quotas.

* `unit` - Indicates the unit of the quota value.
