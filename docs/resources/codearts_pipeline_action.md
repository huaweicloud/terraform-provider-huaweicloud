---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_action"
description: |-
  Manages a CodeArts pipeline action resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_action

Manages a CodeArts pipeline action resource within HuaweiCloud.

## Example Usage

### start a pipeline

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "codehub_id" {}
variable "git_type" {}
variable "git_url" {}
variable "choose_jobs" {}
variable "stage_identifier" {}
variable "job_identifier" {}

resource "huaweicloud_codearts_pipeline_action" "run" {
  action      = "run"
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id

  sources {
    type = "code"

    params {
      codehub_id     = var.codehub_id
      git_type       = var.git_type
      git_url        = var.git_url
      default_branch = "master"

      build_params {
        build_type    = "branch"
        event_type    = "Manual"
        target_branch = "master"
      }
    }
  }

  choose_jobs   = [var.job_identifier]
  choose_stages = [var.stage_identifier]
  description   = "demo"
}
```

### stop a pipeline

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "stop" {
  action          = "stop"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
}
```

### pass a manual review

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "job_run_id" {}
variable "step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "pass" {
  action          = "pass"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  job_run_id      = var.job_run_id
  step_run_id     = var.step_run_id
}
```

### reject a manual review

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "job_run_id" {}
variable "step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "refuse" {
  action          = "refuse"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  job_run_id      = var.job_run_id
  step_run_id     = var.step_run_id
}
```

### pass a delayed execution job

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "job_run_id" {}
variable "step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "delay-pass" {
  action          = "delay-pass"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  job_run_id      = var.job_run_id
  step_run_id     = var.step_run_id
}
```

### reject a delayed execution job

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "job_run_id" {}
variable "step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "delay-refuse" {
  action          = "delay-refuse"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  job_run_id      = var.job_run_id
  step_run_id     = var.step_run_id
}
```

### delay the execution for one hour

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "job_run_id" {}
variable "step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "delay" {
  action          = "delay"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  job_run_id      = var.job_run_id
  step_run_id     = var.step_run_id
}
```

### pass the manual check point

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "pre_step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "manual-pass" {
  action          = "manual-pass"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  step_run_id     = var.pre_step_run_id
}
```

### reject the manual check point

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "pre_step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "manual-refuse" {
  action          = "manual-refuse"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  step_run_id     = var.pre_step_run_id
}
```

### resume a pipeline

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "job_run_id" {}
variable "step_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "resume" {
  action          = "resume"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  job_run_id      = var.job_run_id
  step_run_id     = var.step_run_id
}
```

### cancel a pipeline queuing

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}
variable "queue_id" {}

resource "huaweicloud_codearts_pipeline_action" "cancel-queuing" {
  action          = "cancel-queuing"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
  queue_id        = var.queue_id
}
```

### retry a pipeline

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}
variable "pipeline_run_id" {}

resource "huaweicloud_codearts_pipeline_action" "retry" {
  action          = "retry"
  project_id      = var.codearts_project_id
  pipeline_id     = var.pipeline_id
  pipeline_run_id = var.pipeline_run_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the action.
  Value can be as follows:
  + **run**: start a pipeline
  + **stop**: stop a pipeline
  + **pass**: pass a manual review
  + **refuse**: reject a manual review
  + **delay-pass**: pass a delayed execution job
  + **delay-refuse**: reject a delayed execution job
  + **delay**: delay the execution for one hour
  + **manual-pass**: pass the manual check point
  + **manual-refuse**: reject the manual check point
  + **resume**: resume a pipeline
  + **cancel-queuing**: cancel a pipeline queuing
  + **retry**: retry a pipeline

* `pipeline_id` - (Required, String, NonUpdatable) Specifies the pipeline ID.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `pipeline_run_id` - (Optional, String, NonUpdatable) Specifies the pipeline run ID.

* `job_run_id` - (Optional, String, NonUpdatable) Specifies the pipeline job run ID.

* `step_run_id` - (Optional, String, NonUpdatable) Specifies the pipeline step run ID.

* `queue_id` - (Optional, String, NonUpdatable) Specifies the queued pipeline step run ID.

* `sources` - (Optional, List, NonUpdatable) Specifies the code source information list. Only valid when `action` is **run**.
  The [sources](#block--sources) structure is documented below.

* `variables` - (Optional, List, NonUpdatable) Specifies the custom parameters used. Only valid when `action` is **run**.
  The [variables](#block--variables) structure is documented below.

* `choose_jobs` - (Optional, List, NonUpdatable) Specifies the selected pipeline jobs. Only valid when `action` is **run**.

* `choose_stages` - (Optional, List, NonUpdatable) Specifies the selected pipeline stages. Only valid when `action` is **run**.

* `description` - (Optional, String, NonUpdatable) Specifies the running description. Only valid when `action` is **run**.

<a name="block--sources"></a>
The `sources` block supports:

* `params` - (Required, List, NonUpdatable) Specifies the source parameters.
  The [params](#block--sources--params) structure is documented below.

* `type` - (Required, String, NonUpdatable) Specifies the pipeline source type.

<a name="block--sources--params"></a>
The `params` block supports:

* `git_type` - (Required, String, NonUpdatable) Specifies the code repository type.

* `git_url` - (Required, String, NonUpdatable) Specifies the HTTPS address of the Git repository.

* `alias` - (Optional, String, NonUpdatable) Specifies the code repository alias.

* `build_params` - (Optional, List, NonUpdatable) Specifies the detailed build parameters.
  The [build_params](#block--sources--params--build_params) structure is documented below.

* `change_request_ids` - (Optional, List, NonUpdatable) Specifies the change IDs of the change-triggered pipeline.

* `codehub_id` - (Optional, String, NonUpdatable) Specifies the CodeArts Repo code repository ID.

* `default_branch` - (Optional, String, NonUpdatable) Specifies the default branch of the code repository for pipeline
  execution.

* `endpoint_id` - (Optional, String, NonUpdatable) Specifies the ID of the code source endpoint.

<a name="block--sources--params--build_params"></a>
The `build_params` block supports:

* `build_type` - (Required, String, NonUpdatable) Specifies the code repository trigger type.

* `event_type` - (Optional, String, NonUpdatable) Specifies the event type that triggers the pipeline execution.

* `tag` - (Optional, String, NonUpdatable) Specifies the tag that triggers the pipeline execution.

* `target_branch` - (Optional, String, NonUpdatable) Specifies the branch that triggers the pipeline execution.

<a name="block--variables"></a>
The `variables` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the parameter name.

* `value` - (Required, String, NonUpdatable) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
