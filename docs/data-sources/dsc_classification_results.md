---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_classification_results"
description: |-
  Use this data source to get the classification results of sensitive data identification within HuaweiCloud.
---

# huaweicloud_dsc_classification_results

Use this data source to get the classification results of sensitive data identification within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_dsc_classification_results" "test" {
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

* `classification_list` - The classification results list.
  The [classification_list](#dsc_classification_list_struct) structure is documented below.

<a name="dsc_classification_list_struct"></a>
The `classification_list` block supports:

* `id` - The result ID.

* `project_id` - The project ID.

* `job_id` - The job ID.

* `task_id` - The task ID.

* `ins_id` - The instance ID.

* `asset_id` - The asset ID.

* `asset_name` - The asset name.

* `asset_type` - The asset type.

* `object_name` - The object name.

* `object_full_path` - The object full path.

* `security_level_name` - The security level name.

* `security_level_id` - The security level ID.

* `security_level_color` - The security level color.

* `create_time` - The creation time.

* `update_time` - The update time.

* `scan_time` - The scan time.

* `match_info` - The match information list.
  The [match_info](#dsc_match_info_struct) structure is documented below.

<a name="dsc_match_info_struct"></a>
The `match_info` block supports:

* `template_id` - The template ID.

* `template_name` - The template name.

* `rule_id` - The rule ID.

* `rule_name` - The rule name.

* `security_level_name` - The security level name.

* `security_level_id` - The security level ID.

* `security_level_color` - The security level color.

* `classification_name` - The classification name.

* `classification_id` - The classification ID.

* `matched_detail` - The matched detail.

* `matched_examples` - The matched examples list.
  The [matched_examples](#dsc_matched_examples_struct) structure is documented below.

<a name="dsc_matched_examples_struct"></a>
The `matched_examples` block supports:

* `line_number` - The line number of the match.

* `matched_content` - The matched content.
