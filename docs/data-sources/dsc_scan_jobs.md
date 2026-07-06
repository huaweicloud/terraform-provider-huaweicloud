---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_jobs"
description: |-
  Use this data source to get the list of DSC scan jobs within HuaweiCloud.
---

# huaweicloud_dsc_scan_jobs

Use this data source to get the list of DSC scan jobs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_scan_jobs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `content` - (Optional, String) Specifies the task name.

* `is_new` - (Optional, String) Specifies whether to use new version classification.
  Valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `jobs` - The scan job information list.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - The task ID.

* `name` - The task name.

* `rule_groups` - The rule groups used by the task.

* `scan_templates` - The templates used by the task.

* `cycle` - The execution mode of the task.

* `status` - The current status of the task.

* `last_run_time` - The last execution time of the task.

* `create_time` - The creation time of the task.

* `last_scan_risk` - The risk level result of the last scan.

* `use_nlp` - Whether NLP is used for scanning.

* `open` - The enablement status of the task.

* `topic_urn` - The SMN service notification topic.

* `start_time` - The start time of the task.

* `security_level_name` - The name of the identification result risk level.

* `security_level_color` - The value of the identification result risk level.

* `asset_infos` - The asset information list.

  The [asset_infos](#jobs_asset_infos_struct) structure is documented below.

<a name="jobs_asset_infos_struct"></a>
The `asset_infos` block supports:

* `asset_id` - The asset ID.

* `asset_type` - The asset type.
