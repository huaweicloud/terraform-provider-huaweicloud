---
subcategory: "CodeArts Deploy"
---

# huaweicloud_codearts_deploy_application

Manages a CodeArts deploy application resource within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "project_name" {}
variable "template_id" {}

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id   = var.project_id
  project_name = var.project_name
  template_id  = var.template_id
  name         = "test_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.

  Changing this parameter will create a new resource.

* `project_name` - (Required, String, ForceNew) Specifies the project name for CodeArts service.

  Changing this parameter will create a new resource.

* `template_id` - (Required, String, ForceNew) Specifies the deployment template ID.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the application name. The name consists of 3 to 128 characters,
  including letters, digits, chinese characters or `-_` symbols.

  Changing this parameter will create a new resource.

* `resource_pool_id` - (Optional, String, ForceNew) Specifies the resource pool ID. A resource pool is a collection
  of physical environments that execute deployment commands when deploying software packages.
  If not specified, the resource pool hosted by HuaweiCloud will be used.
  If you want to use your own servers as resource pools, please fill your own resource pool ID.

  Changing this parameter will create a new resource.

  -> Please refer to the following documents to create your own resource pool:
  [Creating an Agent Pool](https://support.huaweicloud.com/intl/en-us/usermanual-devcloud/devcloud_01_0016.html) and
  [Creating an Agent](https://support.huaweicloud.com/intl/en-us/usermanual-devcloud/devcloud_01_0017.html).

* `configs` - (Optional, List, ForceNew) Specifies the deployment parameters.
  If not specified, the application will generate some default preset parameters.
  If specified, the default generated preset parameters will be overridden.

  Changing this parameter will create a new resource.

  The [configs](#DeployApplication_Configs) structure is documented below.

<a name="DeployApplication_Configs"></a>
The `configs` block supports:

* `name` - (Optional, String, ForceNew) Specifies the deployment parameter name, which can be customized.

  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the deployment parameter type, valid values are: **text**,
  **host_group**, **enum** and **encrypt**. Defaults to **text**.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the deployment parameter description.

  Changing this parameter will create a new resource.

* `value` - (Optional, String, ForceNew) Specifies the deployment parameter value.

  Changing this parameter will create a new resource.

* `static_status` - (Optional, Int, ForceNew) Specifies whether the parameter is a static parameter.
  If the value is **1**, the parameter cannot be changed during deployment.
  If the value is **0**, the parameter can be changed and reported to the pipeline. Defaults to **1**.
  This field cannot be set to **0** when `type` is set to **encrypt**.

  Changing this parameter will create a new resource.

* `limits` - (Optional, List, ForceNew) Specifies the enum values. This filed is required when `type` is set to **enum**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time.

* `updated_at` - The update time.

* `state` - The application state. Valid values are as follows:
  + **Available**: Indicates the available state.
  + **Draft**: Indicates the draft state.

* `description` - The description.

* `owner_name` - The username of the application creator.

* `owner_id` - The ID of the application creator.

* `can_modify` - Indicates whether the user has the editing permission.

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_view` - Indicates whether the user has the view permission.

* `can_execute` - Indicates whether the user has the deployment permission

* `can_copy` - Indicates whether the user has the copy permission.

* `can_manage` - Check whether the user has the management permission, including adding, deleting, modifying,
  querying deployment and permission modification.

* `role_id` - The role ID. Valid values are as follows:
  + **-1**: project creator.
  + **0**: application creator.
  + **3**: project manager.
  + **4**: developer.
  + **5**: test manager.
  + **6**: tester.
  + **7**: participant.
  + **8**: viewer.

## Import

The CodeArts deploy application resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_application.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `configs` and `template_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_deploy_application" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      configs,
      template_id,
    ]
  }
}
```
