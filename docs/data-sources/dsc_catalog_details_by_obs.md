---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_catalog_details_by_obs"
description: |-
  Use this data source to get the catalog column details by OBS dimension within HuaweiCloud.
---

# huaweicloud_dsc_catalog_details_by_obs

Use this data source to get the catalog column details by OBS dimension within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_catalog_details_by_obs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type_id` - (Optional, String) Specifies the type ID used to filter OBS bucket dimension of a specific type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The OBS dimension result list.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `bucket_name` - The bucket name.

* `classification_tags` - The classification tag list.

* `security_level_name` - The security level name.
