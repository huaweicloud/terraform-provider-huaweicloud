---
subcategory: "CodeArts Build"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_build_task_records"
description: |-
  Use this data source to get a list of CodeArts build task records.
---

# huaweicloud_codearts_build_task_records

Use this data source to get a list of CodeArts build task records.

## Example Usage

```hcl
variable "build_project_id" {}

data "huaweicloud_codearts_build_task_records" "test" {
  build_project_id = var.build_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `build_project_id` - (Required, String) Specifies the build project ID.

* `branches` - (Optional, List) Specifies the list of branches to search.

* `tags` - (Optional, List) Specifies the list of tags to search.

* `from_date` - (Optional, String) Specifies the start date for the query, format is **yyyy-MM-dd HH:mm:ss**.

* `to_date` - (Optional, String) Specifies the end date for the query, format is **yyyy-MM-dd HH:mm:ss**.

* `triggers` - (Optional, List) Specifies the list of triggers to search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the build record list.
  The [records](#attrblock--records) structure is documented below.

<a name="attrblock--records"></a>
The `records` block supports:

* `branch` - Indicates the branch of the build record.

* `build_duration` - Indicates the build duration of the build record.

* `build_no` - Indicates the build number of the build record.

* `build_record_type` - Indicates the parameters of the build record.
  The [build_record_type](#attrblock--records--build_record_type) structure is documented below.

* `build_yml_path` - Indicates the build yaml path of the build record.

* `build_yml_url` - Indicates the build yaml URL of the build record.

* `create_time` - Indicates the creation time of the build record.

* `daily_build_no` - Indicates the daily build number of the build record.

* `daily_build_number` - Indicates the daily build number of the build record.

* `dev_cloud_build_type` - Indicates the build type of the build record.

* `display_name` - Indicates the display name of the build record.

* `duration` - Indicates the duration of the build record.

* `execution_id` - Indicates the execution ID of the build record.

* `finish_time` - Indicates the finish time of the build record.

* `group_name` - Indicates the group name of the build record.

* `id` - Indicates the unique identifier of the build record.

* `parameters` - Indicates the parameters of the build record.
  The [parameters](#attrblock--records--parameters) structure is documented below.

* `pending_duration` - Indicates the pending duration of the build record.

* `project_id` - Indicates the project ID of the build record.

* `queued_time` - Indicates the queued time of the build record.

* `repository` - Indicates the repository of the build record.

* `revision` - Indicates the revision (commitId) of the build record.

* `schedule_time` - Indicates the scheduled time of the build record.

* `scm_type` - Indicates the SCM type of the build record.

* `scm_web_url` - Indicates the SCM web URL of the build record.

* `start_time` - Indicates the start time of the build record.

* `status` - Indicates the status of the build record.

* `status_code` - Indicates the status code of the build record.

* `trigger_name` - Indicates the trigger name of the build record.

* `trigger_type` - Indicates the trigger type of the build record.

* `user_id` - Indicates the user ID of the build record.

<a name="attrblock--records--build_record_type"></a>
The `build_record_type` block supports:

* `is_rerun` - Indicates the whether the record is rerun.

* `record_type` - Indicates the record type.

* `rerun` - Indicates whether the record is rerun.

* `trigger_type` - Indicates the trigger type.

<a name="attrblock--records--parameters"></a>
The `parameters` block supports:

* `name` - Indicates the parameter name.

* `secret` - Indicates whether the parameter is secret.

* `type` - Indicates the parameter type.

* `value` - Indicates the parameter value.
