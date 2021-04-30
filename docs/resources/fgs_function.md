---
subcategory: "FunctionGraph"
---

# huaweicloud\_fgs\_function

Manages a Function resource within HuaweiCloud.
This is an alternative to `huaweicloud_fgs_function_v2`

## Example Usage

```hcl
resource "huaweicloud_fgs_function" "f_1" {
  name        = "func_1"
  app         = "default"
  agency      = "test"
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

* `region` - (Optional, String, ForceNew) The region in which to create the Function resource. If omitted, the provider-level region will be used. Changing this creates a new Function resource.

* `name` - (Required, String, ForceNew) A unique name for the function. Changing this creates a new function.

* `app` - (Required, String) Group to which the function belongs. Changing this creates a new function.

* `code_type` - (Required, String, ForceNew) Function code type, which can be inline: inline code, zip: ZIP file,
	jar: JAR file or java functions, obs: function code stored in an OBS bucket. Changing this
	creates a new function.

* `code_url` - (Optional, String, ForceNew) This parameter is mandatory when code_type is set to obs. Changing this
	creates a new function.

* `description` - (Optional, String, ForceNew) Description of the function. Changing this creates a new function.

* `code_filename` - (Optional, String, ForceNew) Name of a function file, This field is mandatory only when coe_type is
	set to jar or zip. Changing this creates a new function.

* `handler` - (Required, String, ForceNew) Entry point of the function. Changing this creates a new function.

* `memory_size` - (Required, Int, ForceNew) Memory size(MB) allocated to the function. Changing this creates a new function.

* `runtime` - (Required, String, ForceNew) Environment for executing the function. Changing this creates a new function.

* `timeout` - (Required, Int, ForceNew) Timeout interval of the function, ranges from 3s to 900s. Changing this creates a new function.

* `user_data` - (Optional, String, ForceNew) Key/Value information defined for the function. Changing this creates a new function.

* `agency` - (Optional, String, ForceNew) This parameter is mandatory if the function needs to access other cloud services.
	Changing this creates a new function.

* `func_code` - (Required, String, ForceNew) Function code. When code_type is set to inline, zip, or jar, this parameter is mandatory,
	and the code must be encoded using Base64. Changing this creates a new function.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Functions can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_fgs_function.my-func 7117d38e-4c8f-4624-a505-bd96b97d024c
```
