---
subcategory: "CodeArts Build"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_build_task_action"
description: |-
  Manages a CodeArts Build task action resource within HuaweiCloud.
---

# huaweicloud_codearts_build_task_action

Manages a CodeArts Build task action resource within HuaweiCloud.

## Example Usage

### execute a build task

```hcl
variable "job_id" {}

resource "huaweicloud_codearts_build_task_action" "execute" {
  job_id = var.job_id
  action = "execute"
}
```

### stop a build task

```hcl
variable "job_id" {}
variable "build_no" {}

resource "huaweicloud_codearts_build_task_action" "stop" {
  job_id   = var.job_id
  action   = "stop"
  build_no = var.build_no
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the action. Value can be **execute** and **stop**.

* `job_id` - (Required, String, NonUpdatable) Specifies the build task ID.

* `build_no` - (Optional, String, NonUpdatable) Specifies the build task number, start from 1.
  Only valid when `action` is **stop**.

* `parameter` - (Optional, List, NonUpdatable) Specifies the parameter list. Only valid when `action` is **execute**.
  The [parameter](#block--parameter) structure is documented below.

* `scm` - (Optional, List, NonUpdatable) Specifies the build execution SCM. Only valid when `action` is **execute**.
  The [scm](#block--scm) structure is documented below.

<a name="block--parameter"></a>
The `parameter` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the parameter name.

* `value` - (Required, String, NonUpdatable) Specifies the parameter value.

<a name="block--scm"></a>
The `scm` block supports:

* `build_commit_id` - (Optional, String, NonUpdatable) Specifies the build commit ID.

* `build_tag` - (Optional, String, NonUpdatable) Specifies the build tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `daily_build_number` - Indicates the daily build number.
