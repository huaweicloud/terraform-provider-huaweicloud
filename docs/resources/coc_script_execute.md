---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_execute"
description: ""
---

# huaweicloud_coc_script_execute

Execute a COC script on a specified ECS instance within HuaweiCloud.

-> Please make sure the ECS instance has installed the [UniAgent](https://support.huaweicloud.com/intl/en-us/usermanual-aom2/agent_01_0005.html).

## Example Usage

```hcl
variable "script_id" {}
variable "instance_id" {}

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = var.script_id
  instance_id  = var.instance_id
  timeout      = 600
  execute_user = "root"

  parameters {
    name  = "param1"
    value = "value1"
  }
  parameters {
    name  = "param2"
    value = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `script_id` - (Required, String, NonUpdatable) Specifies the COC script ID.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ECS instance ID.

* `timeout` - (Required, Int, NonUpdatable) Specifies the maximum time to execute the script in seconds.

* `execute_user` - (Required, String, NonUpdatable) Specifies the user to execute the script.

* `parameters` - (Optional, List, NonUpdatable) Specifies the input parameters of the script.
  Up to 20 script parameters can be added.
  The [parameters](#block--parameters) structure is documented below.

* `is_sync` - (Optional, Bool, NonUpdatable) Specifies whether sync data before execute the script. Defaults to **true**.

<a name="block--parameters"></a>
The `parameters` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the parameter.

* `value` - (Required, String, NonUpdatable) Specifies the value of the parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `script_name` - The script name.

* `status` - The status of the script execution.

* `created_at` - The start time of the script execution.

* `finished_at` - The end time of the script execution.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The COC script execution can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_script_execute.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include `instance_id`, `parameters` and `is_sync`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_script_execute" "test" {
    ...

  lifecycle {
    ignore_changes = [
      instance_id, parameters, is_sync
    ]
  }
}
```
