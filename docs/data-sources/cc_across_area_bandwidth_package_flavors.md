---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_across_area_bandwidth_package_flavors"
description: |-
  Use this data source to get the specifications of bandwidth packages for communications across geographic regions.
---

# huaweicloud_cc_across_area_bandwidth_package_flavors

Use this data source to get the specifications of bandwidth packages for communications across geographic regions.

## Example Usage

```hcl
data "huaweicloud_cc_across_area_bandwidth_package_flavors" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `local_area_id` - (Optional, List) Specifies the IDs of a geographic region for querying the list of bandwidth package
  specifications.

* `remote_area_id` - (Optional, List) Specifies the IDs of another geographic region for querying the list of bandwidth
  package specifications.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `area_specifications` - Indicates the list of bandwidth package specifications in a geographic region.

  The [area_specifications](#area_specifications_struct) structure is documented below.

<a name="area_specifications_struct"></a>
The `area_specifications` block supports:

* `id` - Indicates the bandwidth package specification ID.

* `local_area_id` - Indicates the geographic region.
  Cloud Connect is available in the following geographic regions:
  + **Chinese-Mainland**: Chinese mainland
  + **Asia-Pacific**: Asia Pacific
  + **Africa**
  + **Western-Latin-America**: Western Latin America
  + **Eastern-Latin-America**: Eastern Latin America
  + **Northern-Latin-America**: Northern Latin America

* `remote_area_id` - Indicates the geographic region.
  Cloud Connect is available in the following geographic regions:
  + **Chinese-Mainland**: Chinese mainland
  + **Asia-Pacific**: Asia Pacific
  + **Africa**
  + **Western-Latin-America**: Western Latin America
  + **Eastern-Latin-America**: Eastern Latin America
  + **Northern-Latin-America**: Northern Latin America

* `spec_codes` - Indicates the bandwidth package specifications.

  The [spec_codes](#area_specifications_spec_codes_struct) structure is documented below.

<a name="area_specifications_spec_codes_struct"></a>
The `spec_codes` block supports:

* `spec_code` - Indicates the specification code of the bandwidth package.

* `billing_mode` - Indicates the bandwidth package billing option.

* `max_bandwidth` - Indicates the maximum bandwidth.

* `mim_bandwidth` - Indicates the minimum bandwidth.
