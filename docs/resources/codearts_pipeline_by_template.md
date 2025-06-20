---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_by_template"
description: |-
  Manages a CodeArts pipeline created by template resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_by_template

Manages a CodeArts pipeline created by template resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "template_id" {}
variable "name" {}
variable "codehub_id" {}
variable "git_type" {}
variable "git_url" {}
variable "ssh_git_url" {}
variable "repo_name" {}

resource "huaweicloud_codearts_pipeline_by_template" "test" {
  project_id  = var.codearts_project_id
  template_id = var.template_id
  name        = var.name
  is_publish  = false
  description = "demo"

  sources {
    type = "code"

    params {
      codehub_id     = var.codehub_id
      git_type       = var.git_type
      git_url        = var.git_url
      ssh_git_url    = var.ssh_git_url
      repo_name      = var.repo_name
      default_branch = "master"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `template_id` - (Required, String, NonUpdatable) Specifies the CodeArts template ID.

* `is_publish` - (Required, Bool) Specifies whether it is a change-triggered pipeline.

* `name` - (Required, String) Specifies the pipeline name.

* `sources` - (Required, List) Specifies the pipeline source information.
  The [sources](#block--sources) structure is documented below.

* `banned` - (Optional, Bool) Specifies whether the pipeline is banned.

* `component_id` - (Optional, String, NonUpdatable) Specifies the microservice ID.

* `concurrency_control` - (Optional, List) Specifies the pipeline concurrency control information.
  The [concurrency_control](#block--concurrency_control) structure is documented below.

* `definition` - (Optional, String) Specifies the pipeline definition JSON.

* `description` - (Optional, String) Specifies the pipeline description.

* `group_id` - (Optional, String) Specifies the pipeline group ID.

* `manifest_version` - (Optional, String) Specifies the pipeline structure definition version.

* `project_name` - (Optional, String) Specifies the project name.

* `schedules` - (Optional, List) Specifies the pipeline schedule settings.
  The [schedules](#block--schedules) structure is documented below.

* `triggers` - (Optional, List) Specifies the pipeline trigger settings.
  The [triggers](#block--triggers) structure is documented below.

* `variables` - (Optional, List) Specifies the custom variables.
  The [variables](#block--variables) structure is documented below.

<a name="block--sources"></a>
The `sources` block supports:

* `params` - (Optional, List) Specifies the pipeline source parameters.
  The [params](#block--sources--params) structure is documented below.

* `type` - (Optional, String) Specifies the pipeline source type.

<a name="block--sources--params"></a>
The `params` block supports:

* `alias` - (Optional, String) Specifies the code repository alias.

* `codehub_id` - (Optional, String) Specifies the CodeArts Repo code repository ID.

* `default_branch` - (Optional, String) Specifies the default branch.

* `endpoint_id` - (Optional, String) Specifies the code source endpoint ID.

* `git_type` - (Optional, String) Specifies the code repository type.

* `git_url` - (Optional, String) Specifies the HTTPS address of the Git repository.

* `repo_name` - (Optional, String) Specifies the pipeline source name.

* `ssh_git_url` - (Optional, String) Specifies the SSH Git address,

* `web_url` - (Optional, String) Specifies the web page URL.

<a name="block--concurrency_control"></a>
The `concurrency_control` block supports:

* `concurrency_number` - (Optional, Int) Specifies the number of concurrent instances.

* `enable` - (Optional, Bool) Specifies whether to enable the strategy.

* `exceed_action` - (Optional, String) Specifies the policy when the threshold is exceeded.

<a name="block--schedules"></a>
The `schedules` block supports:

* `days_of_week` - (Optional, List) Specifies the execution day in a week.

* `enable` - (Optional, Bool) Specifies whether to enable the schedule job.

* `end_time` - (Optional, String) Specifies the end time.

* `interval_time` - (Optional, String) Specifies the interval time.

* `interval_unit` - (Optional, String) Specifies the interval unit.

* `name` - (Optional, String) Specifies the schedule job name.

* `start_time` - (Optional, String) Specifies the start time.

* `time_zone` - (Optional, String) Specifies the time zone.

* `type` - (Optional, String) Specifies the schedule job type.

<a name="block--triggers"></a>
The `triggers` block supports:

* `callback_url` - (Optional, String) Specifies the callback URL.

* `endpoint_id` - (Optional, String) Specifies the code source endpoint ID.

* `events` - (Optional, List) Specifies the trigger event list.
  The [events](#block--triggers--events) structure is documented below.

* `git_type` - (Optional, String) Specifies the Git repository type.

* `git_url` - (Optional, String) Specifies the Git URL.

* `is_auto_commit` - (Optional, Bool) Specifies whether to automatically commit code.

* `repo_id` - (Optional, String) Specifies the repository ID.

* `security_token` - (Optional, String) Specifies the User token.

<a name="block--triggers--events"></a>
The `events` block supports:

* `enable` - (Optional, Bool) Specifies whether it is available.

* `type` - (Optional, String) Specifies the event type.

<a name="block--variables"></a>
The `variables` block supports:

* `description` - (Optional, String) Specifies the parameter description.

* `is_reset` - (Optional, Bool) Specifies whether to reset.

* `is_runtime` - (Optional, Bool) Specifies whether to set parameters at runtime.

* `is_secret` - (Optional, Bool) Specifies whether it is a private parameter.

* `latest_value` - (Optional, String) Specifies the last parameter value.

* `name` - (Optional, String) Specifies the custom variable name.

* `runtime_value` - (Optional, String) Specifies the value passed in at runtime.

* `sequence` - (Optional, Int) Specifies the parameter sequence, starting from 1.

* `type` - (Optional, String) Specifies the custom parameter type.

* `value` - (Optional, String) Specifies the custom parameter default value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - Indicates the creation time.

* `creator_id` - Indicates the creator ID.

* `creator_name` - Indicates the creator name.

* `is_collect` - Indicates whether the current user has collected it.

* `schedules` - Specifies the pipeline schedule settings.
  The [schedules](#attrblock--schedules) structure is documented below.

* `triggers` - Specifies the pipeline trigger settings.
  The [triggers](#attrblock--triggers) structure is documented below.

* `update_time` - Indicates the last update time.

* `updater_id` - Indicates the last updater ID.

<a name="attrblock--schedules"></a>
The `schedules` block supports:

* `uuid` - Indicates the ID of a scheduled task.

<a name="attrblock--triggers"></a>
The `triggers` block supports:

* `hook_id` - Indicates the callback ID.

## Import

The pipeline created by template can be imported using `project_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_by_template.test <project_id>/<id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `template_id`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the pipeline created by template, or the resource definition should
be updated to align with the pipeline created by template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_pipeline_by_template" "test" {
    ...

  lifecycle {
    ignore_changes = [
      template_id,
    ]
  }
}
```
