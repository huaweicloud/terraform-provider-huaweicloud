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

* `script_id` - (Required, String, ForceNew) Specifies the COC script ID.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ECS instance ID.
  Changing this creates a new resource.

* `timeout` - (Required, Int, ForceNew) Specifies the maximum time to execute the script in seconds.
  Changing this creates a new resource.

* `execute_user` - (Required, String, ForceNew) Specifies the user to execute the script.
  Changing this creates a new resource.

* `parameters` - (Optional, List, ForceNew) Specifies the input parameters of the script.
  Up to 20 script parameters can be added. Changing this creates a new resource.
  The [parameters](#block--parameters) structure is documented below.

<a name="block--parameters"></a>
The `parameters` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of the parameter. Changing this creates a new resource.

* `value` - (Required, String, ForceNew) Specifies the value of the parameter. Changing this creates a new resource.

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
API response, security or some other reason. The missing attributes include `instance_id` and `parameters`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_script_execute" "test" {
    ...

  lifecycle {
    ignore_changes = [
      instance_id, parameters
    ]
  }
}
```
