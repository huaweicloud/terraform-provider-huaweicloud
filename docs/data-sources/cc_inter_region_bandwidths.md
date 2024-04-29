---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_inter_region_bandwidths"
description: ""
---

# huaweicloud_cc_inter_region_bandwidths

Use this data source to get the list of CC inter-region bandwidths.

## Example Usage

```hcl
variable "inter_region_bandwidth_id" {}

data "huaweicloud_cc_inter_region_bandwidths" "test" {
  inter_region_bandwidth_id = var.inter_region_bandwidth_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `inter_region_bandwidth_id` - (Optional, String) Specifies the inter-region bandwidth ID.

* `enterprise_project_id` - (Optional, String) Specifies enterprise project ID.

* `cloud_connection_id` - (Optional, String) Specifies the cloud connection ID.

* `bandwidth_package_id` - (Optional, String) Specifies the bandwidth package ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `inter_region_bandwidths` - The inter-region bandwidth list.

  The [inter_region_bandwidths](#inter_region_bandwidths_struct) structure is documented below.

<a name="inter_region_bandwidths_struct"></a>
The `inter_region_bandwidths` block supports:

* `id` - The inter-region bandwidth ID.

* `name` - The inter-region bandwidth name.

* `bandwidth` - The range of an inter-region bandwidth.

* `description` - The inter-region bandwidth description.

* `cloud_connection_id` - The cloud connection ID.

* `bandwidth_package_id` - The bandwidth package ID.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `inter_regions` - The inter-region bandwidth details.

  The [inter_regions](#inter_region_bandwidths_inter_regions_struct) structure is documented below.

<a name="inter_region_bandwidths_inter_regions_struct"></a>
The `inter_regions` block supports:

* `id` - The inter-region ID.

* `project_id` - The project ID of a region where the inter-region bandwidth is used.

* `remote_region_id` - The ID of another region where an inter-region bandwidth is used.

* `local_region_id` - The ID of one region where an inter-region bandwidth is used.
