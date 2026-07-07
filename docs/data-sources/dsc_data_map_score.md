---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_data_map_score"
description: |-
  Use this data source to get the data map score details of the project within HuaweiCloud.
---

# huaweicloud_dsc_data_map_score

Use this data source to get the data map score details of the project within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_data_map_score" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `analysis_status` - The analysis status.
  The valid values are as follows:
  + `0`: Completed.
  + `1`: Analyzing.
  + `2`: Analysis failed.

* `last_analyze_time` - The timestamp of the last analysis completion.

* `level` - The risk score level.

* `score` - The comprehensive security score of the project.
