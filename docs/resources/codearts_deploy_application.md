---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application"
description: |-
  Manages a CodeArts deploy application resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_application

Manages a CodeArts deploy application resource within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "operation_name" {}
variable "operation_description" {}
variable "operation_code" {}
variable "operation_params" {}
variable "operation_entrance" {}
variable "operation_version" {}
variable "operation_module_id" {}

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id     = var.project_id
  name           = "test_name"
  description    = "test description"
  is_draft       = true
  create_type    = "template"
  trigger_source = "0"

  operation_list {
    name        = var.operation_name
    description = var.operation_description
    code        = var.operation_code
    params      = var.operation_params
    entrance    = var.operation_entrance
    version     = var.operation_version
    module_id   = var.operation_module_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the application name.

* `is_draft` - (Required, Bool) Specifies whether the application is in draft status.
  Valid values:
  + **true**:  Draft state.
  + **false**: Available state.

  -> Only applications in available state can be deployed.
  If `operation_list` is not specified, this field can only be set to **true**.

* `create_type` - (Required, String, ForceNew) Specifies the creation type. Only **template** is supported.
  Changing this parameter will create a new resource.

* `trigger_source` - (Required, String) Specifies where a deployment task can be executed.
  Valid values:
  + **0**: Indicates that all execution requests can be triggered.
  + **1**: Indicates that only pipeline can be triggered.

* `artifact_source_system` - (Optional, String) Specifies the source information transferred by the pipeline.
  This field is only valid when `trigger_source` is set to **1**. Only **CloudArtifact** is supported.
  
* `artifact_type` - (Optional, String) Specifies the artifact type for the pipeline source.
  This field is only valid when `trigger_source` is set to **1**. Valid values are **generic** and **docker**.

* `operation_list` - (Optional, List) Specifies the deployment orchestration list information.

  The [operation_list](#DeployApplication_operation_list) structure is documented below.

* `description` - (Optional, String) Specifies the application description.

* `resource_pool_id` - (Optional, String) Specifies the resource pool ID. A resource pool is a collection
  of physical environments that execute deployment commands when deploying software packages.
  If not specified, the resource pool hosted by HuaweiCloud will be used.
  If you want to use your own servers as resource pools, please fill your own resource pool ID.

  -> Please refer to the following documents to create your own resource pool:
  [Creating an Agent Pool](https://support.huaweicloud.com/intl/en-us/usermanual-devcloud/devcloud_01_0016.html) and
  [Creating an Agent](https://support.huaweicloud.com/intl/en-us/usermanual-devcloud/devcloud_01_0017.html).

* `group_id` - (Optional, String) Specifies the application group ID.
  + When creating the application, if value is empty or **no_grouped**, means the application is ungrouped.
  + If the application is under a specific application group, and you would like to update the application to become
    ungrouped, only specifies it as **no_grouped** is available.

* `is_disable` - (Optional, Bool) Specifies whether to disable the application. Defaults to **false**.
  
  -> When value is **true**, it's unable to update other parameters.

* `permission_level` - (Optional, String) Specifies the permission level.
  Valid values are **instance** and **project**. Defaults to **project**.

<a name="DeployApplication_operation_list"></a>
The `operation_list` block supports:

* `name` - (Optional, String) Specifies the step name.

* `description` - (Optional, String) Specifies the step description.

* `code` - (Optional, String) Specifies the download URL.

* `params` - (Optional, String) Specifies the parameter.

* `entrance` - (Optional, String) Specifies the entry function.

* `version` - (Optional, String) Specifies the version.

* `module_id` - (Optional, String) Specifies the module ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time.

* `updated_at` - The update time.

* `project_name` - The project name.

* `can_modify` - Indicates whether the user has the editing permission.

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_view` - Indicates whether the user has the view permission.

* `can_execute` - Indicates whether the user has the deployment permission

* `can_copy` - Indicates whether the user has the copy permission.

* `can_manage` - Indicates whether the user has the management permission, including adding, deleting, modifying,
  querying deployment and permission modification.

* `can_create_env` - Indicates whether the user has the permission to create an environment.

* `can_disable` - Indicates whether the user has the permission to disable the application.

* `is_care` - Indicates whether the user has favorited the application.

* `task_id` - The deployment task ID.

* `task_name` - The deployment task name.

* `steps` - The deployment steps. The example value is `{"step1":"XXX", "step2":"XXX"}`.

* `permission_matrix` - Indicates the permission matrix.
  The [permission_matrix](#attrblock--permission_matrix) structure is documented below.

<a name="attrblock--permission_matrix"></a>
The `permission_matrix` block supports:

* `role_id` - Indicates the role ID.

* `role_name` - Indicates the role name.

* `role_type` - Indicates the role type.

* `can_modify` - Indicates whether the role has the editing permission.

* `can_delete` - Indicates whether the role has the deletion permission.

* `can_view` - Indicates whether the role has the view permission.

* `can_execute` - Indicates whether the role has the deployment permission.

* `can_copy` - Indicates whether the role has the copy permission.

* `can_manage` - Check whether the role has the management permission, including adding, deleting, modifying,
  querying deployment and permission modification.

* `can_create_env` - Indicates whether the role has the permission to create an environment.

* `can_disable` - Indicates whether the role has the permission to disable the application.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.

## Import

The CodeArts deploy application resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_application.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `is_draft`, `trigger_source`,
`artifact_source_system`, `artifact_type`, `operation_list` and `group_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_deploy_application" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      is_draft,
      trigger_source,
      artifact_source_system,
      artifact_type,
      operation_list,
    ]
  }
}
```
