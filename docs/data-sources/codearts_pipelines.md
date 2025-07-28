---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipelines"
description: |-
  Use this data source to get a list of CodeArts pipelines.
---

# huaweicloud_codearts_pipelines

Use this data source to get a list of CodeArts pipelines.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_pipelines" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `by_group` - (Optional, String) Specifies whether to query by group or not. Valid values are **true** and **false**.

* `component_id` - (Optional, String) Specifies the component ID.

* `creator_ids` - (Optional, List) Specifies the creator ID list.

* `end_time` - (Optional, String) Specifies the end time.

* `executor_ids` - (Optional, List) Specifies the executor ID list.

* `group_path_id` - (Optional, String) Specifies the group ID path.

* `is_banned` - (Optional, String) Specifies whether the pipeline is banned. Valid values are **true** and **false**.

* `is_publish` - (Optional, String) Specifies whether the pipeline is a change pipeline.
  Valid values are **true** and **false**.

* `name` - (Optional, String) Specifies the pipeline name.

* `sort_dir` - (Optional, String) Specifies the sorting rule.

* `sort_key` - (Optional, String) Specifies the sorting field name.
  Valid values can be:
  + **name**: pipeline name
  + **create_time**: pipeline creation time
  + **update_time**: pipeline updating time.

* `start_time` - (Optional, String) Specifies the start time.

* `status` - (Optional, List) Specifies the status.
  Valid values can be:
  + **COMPLETED**: completed
  + **RUNNING**: running
  + **FAILED**: failed
  + **CANCELED**: canceled
  + **PAUSED**: paused
  + **SUSPEND**: suspended
  + **IGNORED**: ignored

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pipelines` - Indicates the pipeline list
  The [pipelines](#attrblock--pipelines) structure is documented below.

<a name="attrblock--pipelines"></a>
The `pipelines` block supports:

* `id` - Indicates the pipeline ID.

* `name` - Indicates the pipeline name.

* `component_id` - Indicates the component ID.

* `convert_sign` - Indicates the sign of converting an old version to a new version.

* `create_time` - Indicates the create time.

* `is_collect` - Indicates whether the pipeline is collected.

* `is_publish` - Indicates whether the pipeline is a change pipeline.

* `latest_run` - Indicates the latest running information.
  The [latest_run](#attrblock--pipelines--latest_run) structure is documented below.

* `manifest_version` - Indicates the pipeline version.

<a name="attrblock--pipelines--latest_run"></a>
The `latest_run` block supports:

* `artifact_params` - Indicates the artifacts after running a pipeline.
  The [artifact_params](#attrblock--pipelines--latest_run--artifact_params) structure is documented below.

* `build_params` - Indicates the build parameters.
  The [build_params](#attrblock--pipelines--latest_run--build_params) structure is documented below.

* `detail_url` - Indicates the address of the details page.

* `end_time` - Indicates the end time.

* `executor_id` - Indicates the executor ID.

* `executor_name` - Indicates the executor name.

* `modify_url` - Indicates the address of the editing page.

* `pipeline_run_id` - Indicates the pipeline run ID.

* `run_number` - Indicates the pipeline running sequence number.

* `stage_status_list` - Indicates the stage information list.
  The [stage_status_list](#attrblock--pipelines--latest_run--stage_status_list) structure is documented below.

* `start_time` - Indicates the start time.

* `status` - Indicates the status of pipeline run.

* `trigger_type` - Indicates the trigger type.

<a name="attrblock--pipelines--latest_run--artifact_params"></a>
The `artifact_params` block supports:

* `branch_filter` - Indicates the branch filter.

* `organization` - Indicates the docker organization.

* `package_name` - Indicates the package name.

* `version` - Indicates the package version.

<a name="attrblock--pipelines--latest_run--build_params"></a>
The `build_params` block supports:

* `action` - Indicates the action.

* `build_type` - Indicates the code repository trigger type.

* `codehub_id` - Indicates the CodeArts Repo code repository ID.

* `commit_id` - Indicates the commit ID.

* `event_type` - Indicates the event type.

* `merge_id` - Indicates the merge ID.

* `message` - Indicates the commit message.

* `source_branch` - Indicates the source branch.

* `source_codehub_http_url` - Indicates the HTTP address of the source Repo code repository.

* `source_codehub_id` - Indicates the ID of the source Repo code repository.

* `source_codehub_url` - Indicates the address of the source Repo code repository.

* `tag` - Indicates the tag that triggers the pipeline execution.

* `target_branch` - Indicates the branch that triggers the pipeline execution.

<a name="attrblock--pipelines--latest_run--stage_status_list"></a>
The `stage_status_list` block supports:

* `id` - Indicates the stage ID.

* `name` - Indicates the stage name.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `end_time` - Indicates the end time.

* `status` - Indicates the stage status.
