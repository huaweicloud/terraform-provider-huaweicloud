---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_batch_execute"
description: |-
  Manages a COC script execution on multiple ECS instances within HuaweiCloud.
---

# huaweicloud_coc_script_batch_execute

Manages a COC script execution on multiple ECS instances within HuaweiCloud.

-> Please make sure each ECS instance has installed the [UniAgent](https://support.huaweicloud.com/intl/en-us/usermanual-aom2/agent_01_0005.html).

## Example Usage

```hcl
variable "script_id" {}
variable "execute_user" {}
variable "execute_batches" {
  type = list(object({
    batch_index  = number
    instance_ids = list(string)
  }))
}
variable "parameters" {
  type = list(object({
    name  = string
    value = string
  }))
}

resource "huaweicloud_coc_script_batch_execute" "test" {
  script_id    = var.script_id
  timeout      = 600
  execute_user = var.execute_user

  dynamic "execute_batches" {
    for_each = var.execute_batches

    content {
      batch_index  = execute_batches.value.batch_index
      instance_ids = execute_batches.value.instance_ids
    }
  }

  dynamic "parameters" {
    for_each = var.parameters

    content {
      name  = parameters.value.name
      value = parameters.value.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `script_id` - (Required, String, NonUpdatable) Specifies the ID of the COC script.

* `execute_batches` - (Required, List, NonUpdatable) Specifies the batch information of the target instances.  
  The [execute_batches](#coc_script_batch_execute_execute_batches) structure is documented below.  
  The maximum number of batches allowed is `20`.

* `timeout` - (Required, Int, NonUpdatable) Specifies the maximum time to execute the script, in seconds.

* `execute_user` - (Required, String, NonUpdatable) Specifies the user to execute the script.

* `parameters` - (Optional, List, NonUpdatable) Specifies the input parameters of the script.  
  The [parameters](#coc_script_batch_execute_parameters) structure is documented below.

* `is_sync` - (Optional, Bool, NonUpdatable) Specifies whether to sync data before executing the script.  
  Defaults to **true**.

<a name="coc_script_batch_execute_execute_batches"></a>
The `execute_batches` block supports:

* `batch_index` - (Required, Int, NonUpdatable) Specifies the batch index.  
  The minimum value is `1`.

* `instance_ids` - (Required, List, NonUpdatable) Specifies the ID list of the ECS instances in this batch.  
  A maximum of `10` instances can be operated in batches.

<a name="coc_script_batch_execute_parameters"></a>
The `parameters` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the parameter.

* `value` - (Required, String, NonUpdatable) Specifies the value of the parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `script_name` - The name of the script.

* `status` - The status of the script execution.

* `created_at` - The start time of the script execution, in RFC3339 format.

* `finished_at` - The end time of the script execution, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_coc_script_batch_execute.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `execute_batches`, `parameters`, `is_sync`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_script_batch_execute" "test" {
  ...

  lifecycle {
    ignore_changes = [
      execute_batches, parameters, is_sync,
    ]
  }
}
```
