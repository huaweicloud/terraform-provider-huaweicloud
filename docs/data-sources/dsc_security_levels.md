---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_security_levels"
description: |-
  Use this data source to get the security level distribution of sensitive data identification results within HuaweiCloud.
---

# huaweicloud_dsc_security_levels

Use this data source to get the security level distribution of sensitive data identification results within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_dsc_security_levels" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the scan job ID.

* `keyword` - (Optional, String) Specifies the keyword for fuzzy search on object names.

* `asset_type` - (Optional, String) Specifies the asset type for filtering.

* `asset_id` - (Optional, String) Specifies the asset ID for filtering.

* `security_level_ids` - (Optional, List) Specifies the security level IDs for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `security_level_list` - The security level distribution list.
  The [security_level_list](#dsc_security_level_list_struct) structure is documented below.

<a name="dsc_security_level_list_struct"></a>
The `security_level_list` block supports:

* `security_level_id` - The security level ID.

* `security_level_name` - The security level name.

* `security_level_color` - The security level color.

* `count` - The number of risk objects under this security level.

* `percent` - The percentage of this security level among all risks.
