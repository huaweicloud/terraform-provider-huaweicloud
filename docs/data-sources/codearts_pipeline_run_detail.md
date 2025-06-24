---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_run_detail"
description: |-
  Use this data source to get the CodeArts pipeline run detail.
---

# huaweicloud_codearts_pipeline_run_detail

Use this data source to get the CodeArts pipeline run detail.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_run_detail" "test" {
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `pipeline_id` - (Required, String) Specifies the pipeline ID.

* `pipeline_run_id` - (Optional, String) Specifies the pipeline run ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `artifacts` - Indicates the artifacts after running a pipeline.
  The [artifacts](#attrblock--artifacts) structure is documented below.

* `component_id` - Indicates the microservice ID.

* `current_system_time` - Indicates the current system time.

* `description` - Indicates the pipeline running description.

* `detail_url` - Indicates the pipeline detail URL.

* `end_time` - Indicates the end time.

* `executor_id` - Indicates the executor ID.

* `executor_name` - Indicates the executor name.

* `group_id` - Indicates the pipeline group ID.

* `group_name` - Indicates the pipeline group name.

* `is_publish` - Indicates whether the pipeline is a change-triggered pipeline.

* `language` - Indicates the language.

* `manifest_version` - Indicates the pipeline version.

* `name` - Indicates the pipeline name.

* `run_number` - Indicates the pipeline running sequence number.

* `sources` - Indicates the pipeline source information.
  The [sources](#attrblock--sources) structure is documented below.

* `stages` - Indicates the stage running information.
  The [stages](#attrblock--stages) structure is documented below.

* `start_time` - Indicates the start time.

* `status` - Indicates the pipeline run status.

* `subject_id` - Indicates the pipeline run ID.

* `trigger_type` - Indicates the trigger type.

<a name="attrblock--artifacts"></a>
The `artifacts` block supports:

* `download_url` - Indicates the artifact download address.

* `name` - Indicates the artifact name.

* `package_type` - Indicates the artifact type.

* `version` - Indicates the artifact version number.

<a name="attrblock--sources"></a>
The `sources` block supports:

* `params` - Indicates the source parameters.
  The [params](#attrblock--sources--params) structure is documented below.

* `type` - Indicates the source type.

<a name="attrblock--sources--params"></a>
The `params` block supports:

* `alias` - Indicates the code repository alias.

* `build_params` - Indicates the build parameters.
  The [build_params](#attrblock--sources--params--build_params) structure is documented below.

* `codehub_id` - Indicates the CodeArts Repo code repository ID.

* `default_branch` - Indicates the default branch.

* `endpoint_id` - Indicates the ID of the code source endpoint.

* `git_type` - Indicates the code repository type.

* `git_url` - Indicates the HTTPS address of the Git repository.

* `repo_name` - Indicates the code repository name.

* `ssh_git_url` - Indicates the SSH address of the Git repository.

* `web_url` - Indicates the address of the code repository page.

<a name="attrblock--sources--params--build_params"></a>
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

<a name="attrblock--stages"></a>
The `stages` block supports:

* `category` - Indicates the stage type.

* `condition` - Indicates the running conditions.

* `depends_on` - Indicates the dependency.

* `end_time` - Indicates the end time.

* `id` - Indicates the stage ID.

* `identifier` - Indicates the unique identifier of a stage.

* `is_select` - Indicates whether to select.

* `jobs` - Indicates the job running information.
  The [jobs](#attrblock--stages--jobs) structure is documented below.

* `name` - Indicates the stage name.

* `parallel` - Indicates whether to execute jobs in parallel.

* `post` - Indicates the stage exit.
  The [post](#attrblock--stages--post) structure is documented below.

* `pre` - Indicates the stage entry.
  The [pre](#attrblock--stages--pre) structure is documented below.

* `run_always` - Indicates whether to always run.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `status` - Indicates the stage status.

<a name="attrblock--stages--jobs"></a>
The `jobs` block supports:

* `async` - Indicates whether it is asynchronous.

* `category` - Indicates the job type.

* `condition` - Indicates the running conditions.

* `depends_on` - Indicates the dependency.

* `end_time` - Indicates the end time.

* `exec_id` - Indicates the job execution ID.

* `id` - Indicates the job ID.

* `identifier` - Indicates the unique identifier of a job.

* `is_select` - Indicates whether the parameter is selected.

* `last_dispatch_id` - Indicates the ID of the job delivered last time.

* `message` - Indicates the error message.

* `name` - Indicates the job name.

* `resource` - Indicates the execution resources.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `status` - Indicates the job status.

* `steps` - Indicates the step running information.
  The [steps](#attrblock--stages--jobs--steps) structure is documented below.

* `timeout` - Indicates the job timeout settings.

<a name="attrblock--stages--jobs--steps"></a>
The `steps` block supports:

* `business_type` - Indicates the extension type.

* `end_time` - Indicates the end time.

* `endpoint_ids` - Indicates the step name.

* `id` - Indicates the step ID.

* `identifier` - Indicates the unique identifier.

* `inputs` - Indicates the step running information.
  The [inputs](#attrblock--stages--jobs--steps--inputs) structure is documented below.

* `last_dispatch_id` - Indicates the ID of the job delivered last time.

* `message` - Indicates the error message.

* `multi_step_editable` - Indicates whether the parameter is editable.

* `name` - Indicates the step name.

* `official_task_version` - Indicates the official extension version.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `status` - Indicates the step status.

* `task` - Indicates the step extension name.

<a name="attrblock--stages--jobs--steps--inputs"></a>
The `task` block supports:

* `key` - Indicates the parameter name.

* `value` - Indicates the parameter value.

<a name="attrblock--stages--post"></a>
The `post` block supports:

* `business_type` - Indicates the extension type.

* `end_time` - Indicates the end time.

* `endpoint_ids` - Indicates the step name.

* `id` - Indicates the step ID.

* `identifier` - Indicates the unique identifier.

* `inputs` - Indicates the step running information.
  The [inputs](#attrblock--stages--post--inputs) structure is documented below.

* `last_dispatch_id` - Indicates the ID of the job delivered last time.

* `message` - Indicates the error message.

* `multi_step_editable` - Indicates whether the parameter is editable.

* `name` - Indicates the step name.

* `official_task_version` - Indicates the official extension version.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `status` - Indicates the step status.

* `task` - Indicates the step extension name.

<a name="attrblock--stages--post--inputs"></a>
The `inputs` block supports:

* `key` - Indicates the parameter name.

* `value` - Indicates the parameter value.

<a name="attrblock--stages--pre"></a>
The `pre` block supports:

* `business_type` - Indicates the extension type.

* `end_time` - Indicates the end time.

* `endpoint_ids` - Indicates the step name.

* `id` - Indicates the step ID.

* `identifier` - Indicates the unique identifier.

* `inputs` - Indicates the step running information.
  The [inputs](#attrblock--stages--pre--inputs) structure is documented below.

* `last_dispatch_id` - Indicates the ID of the job delivered last time.

* `message` - Indicates the error message.

* `multi_step_editable` - Indicates whether the parameter is editable.

* `name` - Indicates the step name.

* `official_task_version` - Indicates the official extension version.

* `sequence` - Indicates the serial number.

* `start_time` - Indicates the start time.

* `status` - Indicates the step status.

* `task` - Indicates the step extension name.

<a name="attrblock--stages--pre--inputs"></a>
The `inputs` block supports:

* `key` - Indicates the parameter name.

* `value` - Indicates the parameter value.
