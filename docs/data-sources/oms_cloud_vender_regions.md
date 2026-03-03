---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_cloud_vender_regions"
description: |-
  Use this data source to query regions supported for a cloud vender.
---

# huaweicloud_oms_cloud_vender_regions

Use this data source to query regions supported for a cloud vender.

## Example Usage

```hcl
data "huaweicloud_oms_cloud_vender_regions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The data source ID.

* `region_info` - The list of supported regions.

  The [region_info](#oms_cloud_vender_regions_region_info_struct) structure is documented below.

<a name="oms_cloud_vender_regions_region_info_struct"></a>
The `region_info` block supports:

* `service_name` - The cloud vender name.

* `region_list` - The region details.

  The [region_list](#oms_cloud_vender_regions_region_list_struct) structure is documented below.

<a name="oms_cloud_vender_regions_region_list_struct"></a>
The `region_list` block supports:

* `cloud_type` - The cloud service name.

* `value` - The region name.

* `description` - The region description.
