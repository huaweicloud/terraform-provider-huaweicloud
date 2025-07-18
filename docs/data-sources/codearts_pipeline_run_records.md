---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_run_records"
description: |-
  Use this data source to get the CodeArts pipeline run records.
---

# huaweicloud_codearts_pipeline_run_records

Use this data source to get the CodeArts pipeline run records.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_run_records" "test" {
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `pipeline_id` - (Required, String) Specifies the pipeline ID.

* `status` - (Optional, List) Specifies the list of status.
  Value can be as follows:
  + **COMPLETED**: completed
  + **RUNNING**: running
  + **FAILED**: failed
  + **CANCELED**: canceled
  + **PAUSED**: paused
  + **SUSPEND**: suspended
  + **IGNORED**: ignored

* `start_time` - (Optional, String) Specifies the start time.

* `end_time` - (Optional, String) Specifies the end time.

* `sort_dir` - (Optional, String) Specifies the sorting sequence. Value can be **asc** and **desc**.

* `sort_key` - (Optional, String) Specifies the sorting attribute. Value can be **start_time**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the pipeline record list.
  The [records](#attrblock--records) structure is documented below.

<a name="attrblock--records"></a>
The `records` block supports:

* `artifact_params` - Indicates the artifacts after running a pipeline.
  The [artifact_params](#attrblock--records--artifact_params) structure is documented below.

* `build_params` - Indicates the build parameters.
  The [build_params](#attrblock--records--build_params) structure is documented below.

* `detail_url` - Indicates the address of the details page.

* `end_time` - Indicates the end time.

* `executor_id` - Indicates the executor ID.

* `executor_name` - Indicates the executor name.

* `modify_url` - Indicates the address of the editing page.

* `pipeline_run_id` - Indicates the pipeline run ID.

* `run_number` - Indicates the pipeline running sequence number.

* `stage_status_list` - Indicates the stage information list.
  The [stage_status_list](#attrblock--records--stage_status_list) structure is documented below.

* `start_time` - Indicates the start time.

* `status` - Indicates the status of pipeline run.

* `trigger_type` - Indicates the trigger type.

<a name="attrblock--records--artifact_params"></a>
The `artifact_params` block supports:

* `branch_filter` - Indicates the branch filter.

* `organization` - Indicates the docker organization.

* `package_name` - Indicates the package name.

* `version` - Indicates the package version.

<a name="attrblock--records--build_params"></a>
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

<a name="attrblock--records--stage_status_list"></a>
The `stage_status_list` block supports:

* `end_time` - Indicates the end time.

* `id` - Indicates the stage ID.

* `name` - Indicates the stage name.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `status` - Indicates the stage status.
