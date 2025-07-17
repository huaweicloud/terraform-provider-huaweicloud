---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline"
description: |-
  Manages a CodeArts pipeline resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline

Manages a CodeArts pipeline resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "name" {}

resource "huaweicloud_codearts_pipeline" "test" {
  project_id  = var.codearts_project_id
  name        = var.name
  description = "demo"
  is_publish  = false
  definition  = jsonencode({
    "stages": [
      {
        "name": "Stage_1",
        "identifier": "xxx",
        "run_condition": null,
        "type": null,
        "sequence": 0,
        "parallel": null,
        "pre": [
          {
            "runtime_attribution": null,
            "multi_step_editable": 0,
            "official_task_version": null,
            "icon_url": null,
            "name": null,
            "task": "official_devcloud_autoTrigger",
            "business_type": null,
            "inputs": null,
            "env": null,
            "sequence": 0,
            "identifier": null,
            "endpoint_ids": null
          }
        ],
        "post": null,
        "jobs": [
          {
            "id": "",
            "identifier_old": null,
            "stage_index": null,
            "type": null,
            "name": "ManualReview",
            "async": null,
            "identifier": "JOB_EyJYf",
            "sequence": 0,
            "condition": "$${{ default() }}",
            "strategy": {
              "select_strategy": "selected"
            },
            "timeout": "",
            "resource": "{\"type\":\"system\",\"arch\":\"x86\"}",
            "steps": [
              {
                "runtime_attribution": "agentless",
                "multi_step_editable": 0,
                "official_task_version": "0.0.5",
                "icon_url": "xxx",
                "name": "ManualReview",
                "task": "official_devcloud_checkpoint",
                "business_type": "Normal",
                "inputs": [
                  {
                    "key": "audit_source",
                    "value": "members"
                  },
                  {
                    "key": "approvers",
                    "value": "xxx"
                  },
                  {
                    "key": "audit_role",
                    "value": ""
                  },
                  {
                    "key": "check_strategy",
                    "value": "all"
                  },
                  {
                    "key": "timeout_strategy",
                    "value": "reject"
                  },
                  {
                    "key": "timeout",
                    "value": 3600
                  },
                  {
                    "key": "comment",
                    "value": ""
                  }
                ],
                "env": [],
                "sequence": 0,
                "identifier": "xxx",
                "endpoint_ids": []
              }
            ],
            "unfinished_steps": [],
            "condition_tag": "",
            "exec_type": "AGENTLESS_JOB",
            "depends_on": [],
            "reusable_job_id": null
          }
        ],
      "depends_on": [],
      "run_always": false
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `name` - (Required, String) Specifies the pipeline name. The value can contain only letters, digits, hyphens (-), and
  underscores (_). The length is **1** to **128** characters.

* `definition` - (Required, String) Specifies the pipeline definition JSON.

* `is_publish` - (Required, Bool) Specifies whether it is a change-triggered pipeline. Defaults to **false**.

* `component_id` - (Optional, String, NonUpdatable) Specifies the microservice ID.

* `sources` - (Optional, List) Specifies the pipeline source information.
  The [sources](#block--sources) structure is documented below.

* `concurrency_control` - (Optional, List) Specifies the pipeline concurrency control information.
  The [concurrency_control](#block--concurrency_control) structure is documented below.

* `schedules` - (Optional, List) Specifies the pipeline schedule settings.
  The [schedules](#block--schedules) structure is documented below.

* `triggers` - (Optional, List) Specifies the pipeline trigger settings.
  The [triggers](#block--triggers) structure is documented below.

* `variables` - (Optional, List) Specifies the custom variables.
  The [variables](#block--variables) structure is documented below.

* `description` - (Optional, String) Specifies the pipeline description.

* `group_id` - (Optional, String) Specifies the pipeline group ID.

* `manifest_version` - (Optional, String) Specifies the pipeline structure definition version. Defaults to **3.0**.

* `project_name` - (Optional, String) Specifies the project name.

* `banned` - (Optional, Bool) Specifies whether the pipeline is banned. Defaults to **false**.

* `parameter_groups` - (Optional, List) Specifies the parameter groups associated with.

<a name="block--sources"></a>
The `sources` block supports:

* `type` - (Optional, String) Specifies the pipeline source type. Value can be **code**.

* `params` - (Optional, List) Specifies the pipeline source parameters.
  The [params](#block--sources--params) structure is documented below.

<a name="block--sources--params"></a>
The `params` block supports:

* `git_type` - (Optional, String) Specifies the code repository type. Value can be **CodeArts Repo**, **Gitee**,
  **GitHub**, **GitCode**, and **GitLab**.

* `alias` - (Optional, String) Specifies the code repository alias. The value can contain a maximum of **128** characters,
  including letters, digits, and underscores (_).

* `codehub_id` - (Optional, String) Specifies the CodeArts Repo code repository ID.

* `default_branch` - (Optional, String) Specifies the default branch.

* `endpoint_id` - (Optional, String) Specifies the code source endpoint ID.

* `git_url` - (Optional, String) Specifies the HTTPS address of the Git repository.

* `repo_name` - (Optional, String) Specifies the pipeline source name.

* `ssh_git_url` - (Optional, String) Specifies the SSH Git address,

* `web_url` - (Optional, String) Specifies the web page URL.

<a name="block--concurrency_control"></a>
The `concurrency_control` block supports:

* `concurrency_number` - (Optional, Int) Specifies the number of concurrent instances.

* `enable` - (Optional, Bool) Specifies whether to enable the strategy. Defaults to **false**.

* `exceed_action` - (Optional, String) Specifies the policy when the threshold is exceeded.
  Value can be as follows:
  + **ABORT**: ignore
  + **QUEUE**: wait in queue

<a name="block--schedules"></a>
The `schedules` block supports:

* `days_of_week` - (Optional, List) Specifies the execution day in a week. Sunday to Saturday: **1** to **7**.

* `enable` - (Optional, Bool) Specifies whether to enable the schedule job. Defaults to **false**.

* `end_time` - (Optional, String) Specifies the end time.

* `interval_time` - (Optional, String) Specifies the interval time.

* `interval_unit` - (Optional, String) Specifies the interval unit.

* `name` - (Optional, String) Specifies the schedule job name.

* `start_time` - (Optional, String) Specifies the start time.

* `time_zone` - (Optional, String) Specifies the time zone. Value can be **China Standard Time**, **GMT Standard Time**,
  **South Africa Standard Time**, **Russian Standard Time**,**SE Asia Standard Time**, **Singapore Standard Time**,
  **Pacific SA Standard Time**, **E. South America Standard Time**, **Central Standard Time (Mexico)**,
  **Egypt Standard Time**, **Saudi Arabia Standard Time**.

* `type` - (Optional, String) Specifies the schedule job type.

<a name="block--triggers"></a>
The `triggers` block supports:

* `callback_url` - (Optional, String) Specifies the callback URL.

* `endpoint_id` - (Optional, String) Specifies the code source endpoint ID.

* `events` - (Optional, List) Specifies the trigger event list.
  The [events](#block--triggers--events) structure is documented below.

* `git_type` - (Optional, String) Specifies the Git repository type. The options include **CodeHub**, **Gitee**,
  **GitHub**, **GitCode**, and **GitLab**.

* `git_url` - (Optional, String) Specifies the Git URL.

* `is_auto_commit` - (Optional, Bool) Specifies whether to automatically commit code. Defaults to **false**.

* `repo_id` - (Optional, String) Specifies the repository ID.

* `security_token` - (Optional, String) Specifies the User token.

<a name="block--triggers--events"></a>
The `events` block supports:

* `enable` - (Optional, Bool) Specifies whether it is available. Defaults to **false**.

* `type` - (Optional, String) Specifies the event type.
  Value can be as follows:
  + **merge_request**: MR
  + **push**: code push
  + **tag_push**: tag
  + **issue**: Gitee repository issue
  + **note**: Gitee repository comment

<a name="block--variables"></a>
The `variables` block supports:

* `description` - (Optional, String) Specifies the parameter description.

* `is_reset` - (Optional, Bool) Specifies whether to reset. Defaults to **false**.

* `is_runtime` - (Optional, Bool) Specifies whether to set parameters at runtime. Defaults to **false**.

* `is_secret` - (Optional, Bool) Specifies whether it is a private parameter. Defaults to **false**.

* `latest_value` - (Optional, String) Specifies the last parameter value.

* `name` - (Optional, String) Specifies the custom variable name. The value can contain a maximum of **32** characters.

* `runtime_value` - (Optional, String) Specifies the value passed in at runtime.

* `sequence` - (Optional, Int) Specifies the parameter sequence, starting from **1**.

* `type` - (Optional, String) Specifies the custom parameter type.
  Value can be as follows:
  + **autoIncrement**: auto-increment parameter
  + **enum**: enumeration parameter
  + **string**: character string parameter

* `value` - (Optional, String) Specifies the custom parameter default value.

* `limits` - (Optional, List) Specifies the list of enumerated values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - Indicates the creation time.

* `creator_id` - Indicates the creator ID.

* `creator_name` - Indicates the creator name.

* `is_collect` - Indicates whether the current user has collected it.

* `update_time` - Indicates the last update time.

* `updater_id` - Indicates the last updater ID.

* `schedules` - Indicates the pipeline schedule settings.
  The [schedules](#attrblock--schedules) structure is documented below.

* `triggers` - Indicates the pipeline trigger settings.
  The [triggers](#attrblock--triggers) structure is documented below.

<a name="attrblock--schedules"></a>
The `schedules` block supports:

* `uuid` - Indicates the ID of a scheduled task.

<a name="attrblock--triggers"></a>
The `triggers` block supports:

* `hook_id` - Indicates the callback ID.

## Import

The pipeline can be imported using `project_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline.test <project_id>/<id>
```
