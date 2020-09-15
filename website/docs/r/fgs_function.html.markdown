---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function"
sidebar_current: "docs-huaweicloud-resource-fgs-function"
description: |-
  Manages a Function resource within HuaweiCloud.
---

# huaweicloud\_fgs\_function

Manages a Function resource within HuaweiCloud.
This is an alternative to `huaweicloud_fgs_function_v2`

## Example Usage

```hcl
resource "huaweicloud_fgs_function" "f_1" {
  name        = "func_1"
  package     = "default"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the function. Changing this creates a new function.

* `package` - (Required) Group to which the function belongs. Changing this creates a new function.

* `code_type` - (Required) Function code type, which can be inline: inline code, zip: ZIP file,
	jar: JAR file or java functions, obs: function code stored in an OBS bucket. Changing this
	creates a new function.

* `code_url` - (Optional) This parameter is mandatory when code_type is set to obs. Changing this
	creates a new function.

* `description` - (Optional) Description of the function. Changing this creates a new function.

* `code_filename` - (Optional) Name of a function file, This field is mandatory only when coe_type is
	set to jar or zip. Changing this creates a new function.

* `handler` - (Required) Entry point of the function. Changing this creates a new function.

* `memory_size` - (Required) Memory size(MB) allocated to the function. Changing this creates a new function.

* `runtime` - (Required) Environment for executing the function. Changing this creates a new function.

* `timeout` - (Required) Timeout interval of the function, ranges from 3s to 900s. Changing this creates a new function.

* `user_data` - (Optional) Key/Value information defined for the function. Changing this creates a new function.

* `xrole` - (Optional) This parameter is mandatory if the function needs to access other cloud services.
	Changing this creates a new function.

* `func_code` - (Required) Function code. When code_type is set to inline, zip, or jar, this parameter is mandatory,
	and the code must be encoded using Base64. Changing this creates a new function.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `package` - See Argument Reference above.
* `code_type` - See Argument Reference above.
* `code_url` - See Argument Reference above.
* `description` - See Argument Reference above.
* `code_filename` - See Argument Reference above.
* `handler` - See Argument Reference above.
* `memory_size` - See Argument Reference above.
* `runtime` - See Argument Reference above.
* `timeout` - See Argument Reference above.
* `user_data` - See Argument Reference above.
* `xrole` - See Argument Reference above.
* `func_code` - See Argument Reference above.

## Import

Functions can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_fgs_function.my-func 7117d38e-4c8f-4624-a505-bd96b97d024c
```
