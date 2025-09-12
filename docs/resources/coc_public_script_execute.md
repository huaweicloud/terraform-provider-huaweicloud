---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_public_script_execute"
description: |-
  Manages a COC public script execute resource within HuaweiCloud.
---

# huaweicloud_coc_public_script_execute

Manages a COC public script execute resource within HuaweiCloud.

~> Deleting public script execute resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "script_uuid" {}
variable "resource_id" {}
variable "region_id" {}

resource "huaweicloud_coc_public_script_execute" "test" {
  script_uuid  = var.script_uuid
  timeout      = 300
  success_rate = 100
  execute_user = "root"
  script_params {
    param_name  = "action"
    param_value = "stop"
  }
  execute_batches {
    batch_index = 1
    target_instances {
      resource_id        = var.resource_id
      region_id          = var.region_id
      cloud_service_name = "ECS"
    }
    rotation_strategy = "CONTINUE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `script_uuid` - (Required, String, NonUpdatable) Specifies the public script UUID.

* `timeout` - (Required, Int, NonUpdatable) Specifies the timeout period.

* `success_rate` - (Required, Float, NonUpdatable) Specifies the success rate.

* `execute_user` - (Required, String, NonUpdatable) Specifies the script execution user, such as root.

* `execute_batches` - (Required, List, NonUpdatable) Specifies the target instance batch information.
  The [execute_batches](#execute_batches_struct) structure is documented below.

* `script_params` - (Optional, List, NonUpdatable) Specifies the script input parameter list.
  The [script_params](#script_params_struct) structure is documented below.

<a name="execute_batches_struct"></a>
The `execute_batches` block supports:

* `batch_index` - (Required, Int, NonUpdatable) Specifies the batch index, starting from **1**.

* `target_instances` - (Required, List, NonUpdatable) Specifies the list of target nodes.
  The [target_instances](#execute_batches_target_instances_struct) structure is documented below.

* `rotation_strategy` - (Required, String, NonUpdatable) Specifies the pause-resume policy.
  Values can be **CONTINUE** or **PAUSE**.

<a name="execute_batches_target_instances_struct"></a>
The `target_instances` block supports:

* `resource_id` - (Required, String, NonUpdatable) Specifies the ECS cloud server ID.

* `region_id` - (Required, String, NonUpdatable) Specifies the region ID.

* `cloud_service_name` - (Optional, String, NonUpdatable) Specifies the resource provider. The default value is **ECS**.

* `type` - (Optional, String, NonUpdatable) Specifies the resource type under the resource provider. The default value
  is **CLOUDSERVER**.

* `custom_attributes` - (Optional, List, NonUpdatable) Specifies the five user-defined attributes in the key_value
  format are supported.
  The [custom_attributes](#execute_batches_target_instances_custom_attributes_struct) structure is documented below.

<a name="execute_batches_target_instances_custom_attributes_struct"></a>
The `custom_attributes` block supports:

* `key` - (Required, String, NonUpdatable) Specifies the custom attribute key.

* `value` - (Required, String, NonUpdatable) Specifies the custom attribute value.

<a name="script_params_struct"></a>
The `script_params` block supports:

* `param_name` - (Required, String, NonUpdatable) Specifies the name of the script input parameter.
  The parameter name cannot be repeated in the same script.

* `param_value` - (Required, String, NonUpdatable) Specifies the value of the script input parameter.

  -> This is required by default. It can be empty when **param_refer** is not empty.

* `param_refer` - (Optional, List, NonUpdatable) Specifies the parameter reference.
  The [param_refer](#script_params_param_refer_struct) structure is documented below.

<a name="script_params_param_refer_struct"></a>
The `param_refer` block supports:

* `refer_type` - (Required, String, NonUpdatable) Specifies the parameter reference type: PARAM_STORE.

* `param_id` - (Required, String, NonUpdatable) Specifies the unique primary key id of the reference parameter.
  Values can be **LOW**, **MEDIUM** and **HIGH**.

* `param_version` - (Optional, String, NonUpdatable) Specifies the version number of the reference parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `gmt_created` - Indicates the execution creation time.

* `gmt_finished` - Indicates the execution completion time.

* `execute_costs` - Indicates the execution time, the unit is seconds.

* `creator` - Indicates the creator.

* `status` - Indicates the execution status.

* `script_name` - Indicates the public script name.

* `script_version_uuid` - Indicates the public script version UUID.

* `script_version_name` - Indicates the public script version name.

* `current_execute_batch_index` - Indicates the index of the currently executed batch.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

The COC public script execution can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_public_script_execute.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include `execute_batches`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_public_script_execute" "test" {
    ...

  lifecycle {
    ignore_changes = [
      execute_batches
    ]
  }
}
```
