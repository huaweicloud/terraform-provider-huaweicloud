---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_across_regions_bandwidth_package_flavors"
description: |-
  Use this data source to get the specifications of bandwidth packages for communications across regions.
---

# huaweicloud_cc_across_regions_bandwidth_package_flavors

Use this data source to get the specifications of bandwidth packages for communications across regions.

## Example Usage

```hcl
data "huaweicloud_cc_across_regions_bandwidth_package_flavors" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `local_region_id` - (Optional, String) Specifies the ID of a region for querying the list of multi-city bandwidth package
  configurations.

* `remote_region_id` - (Optional, String) Specifies the ID of another region for querying the list of multi-city bandwidth
  package configurations.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region_specifications` - Indicates the list of the multi-city bandwidth package specifications.

  The [region_specifications](#region_specifications_struct) structure is documented below.

<a name="region_specifications_struct"></a>
The `region_specifications` block supports:

* `id` - Indicates the specification ID of the bandwidth package for communications between regions.

* `name` - Indicates the specification name of the bandwidth package for communications between regions.

* `en_name` - Indicates the name of the bandwidth package for communications between regions.

* `es_name` - Indicates the specification name in Spanish of the bandwidth package for communications between
  regions.

* `pt_name` - Indicates the specification name in Portuguese of the bandwidth package for communications between
  regions.

* `local_region_id` - Indicates the local region ID.

* `remote_region_id` - Indicates the remote region ID.

* `spec_codes` - Indicates the list of the bandwidth package specifications.

  The [spec_codes](#region_specifications_spec_codes_struct) structure is documented below.

<a name="region_specifications_spec_codes_struct"></a>
The `spec_codes` block supports:

* `spec_code` - Indicates the specification code of the bandwidth package.

* `billing_mode` - Indicates the bandwidth package billing option.

* `max_bandwidth` - Indicates the maximum bandwidth.

* `mim_bandwidth` - Indicates the minimum bandwidth.
