---
subcategory: "FunctionGraph"
---

# huaweicloud_fgs_function

Manages a Function resource within HuaweiCloud.

## Example Usage

### With base64 func code

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

### With text code

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
  func_code   = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}
```

### Create function using SWR image and enable the configuration of the asynchronous invocation

```hcl
variable "function_name" {}
# Authorize the SWR, OBS, and SMN service operation permissions for the FunctionGraph service.
variable "agency_name" {}
variable "image_url" {}
variable "bucket_name" {}
variable "topic_urn" {}

resource "huaweicloud_fgs_function" "by_swr_image" {
  name        = var.function_name
  agency      = var.agency_name
  handler     = "-"
  app         = "default"
  runtime     = "Custom Image"
  memory_size = 128
  timeout     = 3

  custom_image {
    url = var.image_url
  }

  async_invoke {
    max_async_event_age_in_seconds = 3500
    max_async_retry_attempts       = 2
    enable_async_status_log        = true

    on_success {
      destination = "OBS"
      param = jsonencode({
        bucket  = var.bucket_name
        prefix  = "/success"
        expires = 5
      })
    }

    on_failure {
      destination = "SMN"
      param       = jsonencode({
        topic_urn = var.topic_urn
      })
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Function resource. If omitted, the
  provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the function.
  Changing this will create a new resource.

* `app` - (Required, String) Specifies the group to which the function belongs.

* `memory_size` - (Required, Int) Specifies the memory size(MB) allocated to the function.

* `runtime` - (Required, String, ForceNew) Specifies the environment for executing the function.
  If the function is created using an SWR image, set this parameter to `Custom Image`.
  Changing this will create a new resource.

* `timeout` - (Required, Int) Specifies the timeout interval of the function, ranges from 3s to 900s.

* `code_type` - (Optional, String) Specifies the function code type, which can be:
  + **inline**: inline code.
  + **zip**: ZIP file.
  + **jar**: JAR file or java functions.
  + **obs**: function code stored in an OBS bucket.

* `handler` - (Required, String) Specifies the entry point of the function.

-> If the function is created using an SWR image, keep `code_type` empty and use **-** to set the handler.

* `functiongraph_version` - (Optional, String, ForceNew) Specifies the FunctionGraph version, defaults to **v1**.
  + **v1**: Hosts event-driven functions in a serverless context.
  + **v2**: Next-generation function hosting service powered by Huawei YuanRong architecture.

  Changing this will create a new resource.

* `func_code` - (Optional, String) Specifies the function code. When code_type is set to inline, zip, or jar, this
  parameter is mandatory, and the code can be encoded using Base64 or just with the text code.

* `code_url` - (Optional, String) Specifies the code url. This parameter is mandatory when code_type is set to obs.

* `code_filename` - (Optional, String) Specifies the name of a function file, This field is mandatory only when coe_type
  is set to jar or zip.

* `depend_list` - (Optional, List) Specifies the ID list of the dependencies.

* `user_data` - (Optional, String) Specifies the Key/Value information defined for the function. Key/value data might be
  parsed with [Terraform `jsonencode()` function]('https://www.terraform.io/docs/language/functions/jsonencode.html').

* `encrypted_user_data` - (Optional, String) Specifies the key/value information defined to be encrypted for the
  function. The format is the same as `user_data`.

* `agency` - (Optional, String) Specifies the agency. This parameter is mandatory if the function needs to access other
  cloud services.

* `app_agency` - (Optional, String) Specifies An execution agency enables you to obtain a token or an AK/SK for
  accessing other cloud services.

* `description` - (Optional, String) Specifies the description of the function.

* `initializer_handler` - (Optional, String) Specifies the initializer of the function.

* `initializer_timeout` - (Optional, Int) Specifies the maximum duration the function can be initialized. Value range:
  1s to 300s.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the function.
  Changing this will create a new resource.

* `vpc_id`  - (Optional, String) Specifies the ID of VPC.

* `network_id`  - (Optional, String) Specifies the network ID of subnet.

-> **NOTE:** An agency with VPC management permissions must be specified for the function.

* `mount_user_id` - (Optional, Int) Specifies the user ID, a non-0 integer from –1 to 65534. Default to -1.

* `mount_user_group_id` - (Optional, Int) Specifies the user group ID, a non-0 integer from –1 to 65534. Default to
  -1.

* `func_mounts` - (Optional, List) Specifies the file system list. The `func_mounts` object structure is documented
  below.

* `async_invoke` - (Optional, List) Specifies the configuration of the asynchronous execution notification.  
  The [object](#functiongraph_async_invoke) structure is documented below.

* `custom_image` - (Optional, List, ForceNew) Specifies the custom image configuration for creating function.
  The [object](#functiongraph_custom_image) structure is documented below.
  Changing this will create a new resource.

The `func_mounts` block supports:

* `mount_type` - (Required, String) Specifies the mount type. Options: sfs, sfsTurbo, and ecs.

* `mount_resource` - (Required, String) Specifies the ID of the mounted resource (corresponding cloud service).

* `mount_share_path` - (Required, String) Specifies the remote mount path. Example: 192.168.0.12:/data.

* `local_mount_path` - (Required, String) Specifies the function access path.

<a name="functiongraph_async_invoke"></a>
The `async_invoke` block supports:

* `max_async_event_age_in_seconds` - (Optional, Int) Specifies the maximum validity period of a message.

* `max_async_retry_attempts` - (Optional, Int) Specifies the maximum number of retry attempts to be made if
  asynchronous invocation fails.

* `on_success` - (Optional, List) Specifies the target to be invoked when a function is successfully executed.  
  The [object](#functiongraph_async_invoke) structure is documented below.

* `on_failure` - (Optional, List) Specifies the target to be invoked when a function fails to be executed due to a
  system error or an internal error.  
  The [object](#functiongraph_destination_config) structure is documented below.

* `enable_async_status_log` - (Optional, Bool) Specifies the URL of SWR image, the URL must start with `swr.`.

<a name="functiongraph_destination_config"></a>
The `on_success` and the `on_failure` blocks support:

* `destination` - (Optional, String) Specifies the object type.  
  The valid values are as follows:
  + **OBS**
  + **SMN**
  + **DIS**
  + **FunctionGraph**

* `param` - (Optional, String) Specifies the parameters (map object in JSON format) corresponding to the target service.
  + The **OBS** objects include: `bucket` (bucket name), `prefix` (object directory prefix) and `expires` (object
    expiration time, the valid value ranges from `0` to `365`. If the value is `0`, the object will not expire.).
  + The **SMN** objects include: `topic_urn`.
  + The **DIS** objects include: `stream_name`.
  + The **FunctionGraph** objects include: `func_urn` (function URN).

-> If you enable the destination function, you must be ensured that the agent contains the operation authority of the
   corresponding service.

<a name="functiongraph_custom_image"></a>
The `custom_image` block supports:

* `url` - (Required, String, ForceNew) Specifies the URL of SWR image, the URL must start with `swr.`.
  Changing this will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `func_mounts/status` - The status of file system.
* `urn` - Uniform Resource Name
* `version` - The version of the function

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Functions can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_fgs_function.my-func 7117d38e-4c8f-4624-a505-bd96b97d024c
```
