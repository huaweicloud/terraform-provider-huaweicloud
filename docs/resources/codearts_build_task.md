---
subcategory: "CodeArts Build"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_build_task"
description: |-
  Manages a CodeArts Build task resource within HuaweiCloud.
---

# huaweicloud_codearts_build_task

Manages a CodeArts Build task resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "name" {}
variable "url" {}
variable "web_url" {}
variable "scm_type" {}
variable "repo_id" {}

resource "huaweicloud_codearts_build_task" "test" {
  project_id = var.codearts_project_id
  name       = var.name
  arch       = "x86-64"

  scms {
    url      = var.url
    scm_type = var.scm_type
    web_url  = var.web_url
    repo_id  = var.repo_id
    branch   = "master"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.
  Changing this creates a new resource.

* `arch` - (Required, String) Specifies the architecture of the build machine.

* `name` - (Required, String) Specifies the name of the build task.

* `scms` - (Optional, List) Specifies the build execution SCM.
  The [scms](#block--scms) structure is documented below.

* `auto_update_sub_module` - (Optional, String) Specifies whether to automatically update submodules.

* `build_config_type` - (Optional, String) Specifies the build task configuration type.

* `build_if_code_updated` - (Optional, Bool) Specifies whether to enable the code commit trigger build switch.
  Defaults to **false**.

* `flavor` - (Optional, String) Specifies the specification of the execution machine.

* `group_id` - (Optional, String) Specifies the task group ID.

* `host_type` - (Optional, String) Specifies the host type.

* `parameters` - (Optional, List) Specifies the build execution parameter list.
  The [parameters](#block--parameters) structure is documented below.

* `steps` - (Optional, List) Specifies the build execution steps.
  The [steps](#block--steps) structure is documented below.

* `triggers` - (Optional, List) Specifies the collection of timed task triggers.
  The [triggers](#block--triggers) structure is documented below.

<a name="block--parameters"></a>
The `parameters` block supports:

* `name` - (Optional, String) Specifies the parameter definition name.
  Defaults to **hudson.model.StringParameterDefinition**.

* `params` - (Optional, List) Specifies the build execution sub-parameters.
  The [params](#block--parameters--params) structure is documented below.

<a name="block--parameters--params"></a>
The `params` block supports:

* `limits` - (Optional, List) Specifies the enumeration parameter restrictions.
  The [limits](#block--parameters--params--limits) structure is documented below.

* `name` - (Optional, String) Specifies the parameter field name.

* `value` - (Optional, String) Specifies the parameter field value.

<a name="block--parameters--params--limits"></a>
The `limits` block supports:

* `disable` - (Optional, String) Specifies whether it is effective. Defaults to **0**, which is effective.

* `display_name` - (Optional, String) Specifies the displayed name of the parameter.

* `name` - (Optional, String) Specifies the parameter name.

<a name="block--scms"></a>
The `scms` block supports:

* `repo_id` - (Required, String) Specifies the repository ID.

* `scm_type` - (Required, String) Specifies the source code management type.

* `url` - (Required, String) Specifies the repository URL.

* `web_url` - (Required, String) Specifies the web URL of the repository.

* `branch` - (Optional, String) Specifies the branch name.

* `build_type` - (Optional, String) Specifies the build type.

* `depth` - (Optional, String) Specifies the depth.

* `enable_git_lfs` - (Optional, Bool) Specifies whether to enable Git LFS. Defaults to **false**.

* `endpoint_id` - (Optional, String) Specifies the endpoint ID.

* `group_name` - (Optional, String) Specifies the group name.

* `is_auto_build` - (Optional, Bool) Specifies whether to automatically build. Defaults to **false**.

* `repo_name` - (Optional, String) Specifies the repository name.

* `source` - (Optional, String) Specifies the source type.

<a name="block--steps"></a>
The `steps` block supports:

* `module_id` - (Required, String) Specifies the build step module ID.

* `name` - (Required, String) Specifies the build step name.

* `enable` - (Optional, Bool) Specifies whether to enable the step. Defaults to **false**.

* `properties` - (Optional, Map) Specifies the build step properties. Value is JSON format string.

* `version` - (Optional, String) Specifies the build step version.

<a name="block--triggers"></a>
The `triggers` block supports:

* `name` - (Required, String) Specifies the trigger type.

* `parameters` - (Required, List) Specifies the custom parameters.
  The [parameters](#block--triggers--parameters) structure is documented below.

<a name="block--triggers--parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `steps` - Indicates the build execution steps.
  The [steps](#attrblock--steps) structure is documented below.

<a name="attrblock--steps"></a>
The `steps` block supports:

* `properties_all` - Indicates the build step properties.

## Import

The build task can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_build_task.test <id>
```
