---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function"
description: ""
---

# huaweicloud_fgs_function

Manages a Function resource within HuaweiCloud.

## Example Usage

### With base64 func code

```hcl
variable "function_name" {}
variable "function_codes" {}
variable "agency_name" {}

resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = "default"
  agency      = var.agency_name
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.function_codes)
}
```

### With text code

```hcl
variable "function_name" {}
variable "agency_name" {}

resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = "default"
  agency      = var.agency_name
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

### Create function using SWR image

```hcl
variable "function_name" {}
variable "agency_name" {} // The agent name that authorizes FunctionGraph service SWR administrator privilege
variable "image_url" {}

resource "huaweicloud_fgs_function" "by_swr_image" {
  name        = var.function_name
  agency      = var.agency_name
  handler     = "-"
  app         = "default"
  runtime     = "Custom Image"
  code_type   = "Custom-Image-Swr"
  memory_size = 128
  timeout     = 3

  custom_image {
    url = var.image_url
  }
}
```

### Create function with an alias for latest version

```hcl
variable "function_name" {}
variable "function_codes" {}

resource "huaweicloud_fgs_function" "with_alias" {
  name        = var.function_name
  app         = "default"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.function_codes)

  versions {
    name = "latest"

    aliases {
      name = "demo"
    }
  }
}
```

### Create function with VPC access and DNS configuration

```hcl
variable "function_name" {}
variable "agency_name" {} # Allow VPC and DNS permissions for FunctionGraph service
variable "function_codes" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_dns_zone" "test" {
  count = 3

  zone_type = "private"
  name      = format("functiondebug.example%d.com.", count.index)

  router {
    router_id = var.vpc_id
  }
}

resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  memory_size = 128
  runtime     = "Python3.10"
  timeout     = 3
  func_code   = base64encode(var.function_codes)

  # VPC access and DNS configuration
  agency     = var.agency_name
  vpc_id     = var.vpc_id
  network_id = var.network_id
  dns_list   = jsonencode(
    [for v in huaweicloud_dns_zone.test[*] : tomap({id=v.id, domain_name=v.name})]
  )
}
```

### Create function with log group and stream

```hcl
variable "function_name" {}
variable "function_codes" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "log_group_name" {}
variable "log_stream_name" {}

resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = "default"
  agency      = "test"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.function_codes)

  log_group_id    = var.log_group_id
  log_stream_id   = var.log_stream_id
  log_group_name  = var.log_group_name
  log_stream_name = var.log_stream_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Function resource.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the function.
  Changing this will create a new resource.

* `app` - (Required, String) Specifies the group to which the function belongs.

* `memory_size` - (Required, Int) Specifies the memory size allocated to the function, in MByte (MB).

* `runtime` - (Required, String, ForceNew) Specifies the environment for executing the function.
  The valid values are as follows:
  + **Java8**
  + **Java11**
  + **Node.js6.10**
  + **Node.js8.10**
  + **Node.js10.16**
  + **Node.js12.13**
  + **Node.js14.18**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Go1.8**
  + **Go1.x**
  + **C#(.NET Core 2.0)**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **PHP7.3**
  + **Custom**
  + **http**

  If the function is created using an SWR image, set this parameter to `Custom Image`.
  Changing this will create a new resource.

* `timeout` - (Required, Int) Specifies the timeout interval of the function, in seconds.  
  The value ranges from `3` to `900`.

* `code_type` - (Required, String) Specifies the function code type, which can be:
  + **inline**: inline code.
  + **zip**: ZIP file.
  + **jar**: JAR file or java functions.
  + **obs**: function code stored in an OBS bucket.
  + **Custom-Image-Swr**: function code comes from the SWR custom image.

* `handler` - (Required, String) Specifies the entry point of the function.

-> If the function is created using an SWR image, keep `code_type` empty and use **-** to set the handler.

* `functiongraph_version` - (Optional, String, ForceNew) Specifies the FunctionGraph version, default value is **v2**.
  Some regions support only v1, the default value is **v1**.
  + **v1**: Hosts event-driven functions in a serverless context.
  + **v2**: Next-generation function hosting service powered by Huawei YuanRong architecture.

  Changing this will create a new resource.

* `func_code` - (Optional, String) Specifies the function code.  
  The code value can be encoded using **Base64** or just with the text code.  
  Required if the `code_type` is set to **inline**, **zip**, or **jar**.

* `code_url` - (Optional, String) Specifies the code url.  
  Required if the `code_type` is set to **obs**.

* `code_filename` - (Optional, String) Specifies the name of a function file.  
  Required if the `code_type` is set to **jar** or **zip**.

* `depend_list` - (Optional, List) Specifies the ID list of the dependencies.

* `user_data` - (Optional, String) Specifies the Key/Value information defined for the function. Key/value data might be
  parsed with [Terraform `jsonencode()` function]('https://www.terraform.io/docs/language/functions/jsonencode.html').

* `encrypted_user_data` - (Optional, String) Specifies the key/value information defined to be encrypted for the
  function. The format is the same as `user_data`.

* `agency` - (Optional, String) Specifies the agency. This parameter is mandatory if the function needs to access other
  cloud services.

* `app_agency` - (Optional, String) Specifies the execution agency enables you to obtain a token or an AK/SK for
  accessing other cloud services.

* `description` - (Optional, String) Specifies the description of the function.

* `initializer_handler` - (Optional, String) Specifies the initializer of the function.

* `initializer_timeout` - (Optional, Int) Specifies the maximum duration the function can be initialized. Value range:
  1s to 300s.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  function belongs.  
  Changing this will create a new resource.

* `vpc_id` - (Optional, String) Specifies the ID of VPC.

* `network_id` - (Optional, String) Specifies the network ID of subnet.

  -> An agency with VPC management permissions must be specified for the function.

* `dns_list` - (Optional, String) Specifies the private DNS configuration of the function network.
  Private DNS list is associated to the function by a string in the following format:  
  `[{\"id\":\"ff8080828a07ffea018a17184aa310f5\","domain_name":"functiondebug.example1.com."}]`

  -> Ensure the agency with DNS management permissions specified before using this parameter.

* `mount_user_id` - (Optional, Int) Specifies the user ID, a non-0 integer from `–1` to `65,534`.
  Defaults to `-1`.

* `mount_user_group_id` - (Optional, Int) Specifies the user group ID, a non-0 integer from `–1` to `65,534`.
  Defaults to `-1`.

* `func_mounts` - (Optional, List) Specifies the file system list. The `func_mounts` object structure is documented
  below.

* `custom_image` - (Optional, List) Specifies the custom image configuration for creating function.
  The [object](#functiongraph_custom_image) structure is documented below.  
  Required if the parameter `code_type` is **Custom-Image-Swr**.

* `max_instance_num` - (Optional, String) Specifies the maximum number of instances of the function.  
  The valid value ranges from `-1` to `1,000`, defaults to `400`.
  + The minimum value is `-1` and means the number of instances is unlimited.
  + `0` means this function is disabled.
  + The empty value means to keep the default (latest updated) value.

  -> This parameter is only supported by the `v2` version of the function.

* `versions` - (Optional, List) Specifies the versions management of the function.
  The [object](#functiongraph_versions_management) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the function.

* `log_group_id` - (Optional, String) Specifies the ID of the LTS log group.

* `log_group_name` - (Optional, String) Specifies the name of the LTS log group.

* `log_stream_id` - (Optional, String) Specifies the ID of the LTS log stream.

* `log_stream_name` - (Optional, String) Specifies the name of the LTS stream.

* `reserved_instances` - (Optional, List) Specifies the reserved instance policies of the function.
  The [reserved_instances](#functiongraph_reserved_instances) structure is documented below.

* `concurrency_num` - (Optional, Int) Specifies the number of concurrent requests of the function.
  The valid value ranges from `1` to `1,000`, the default value is `1`.
  
  -> This parameter is only supported by the `v2` version of the function.

* `gpu_memory` - (Optional, Int) Specifies the GPU memory size allocated to the function, in MByte (MB).
  The valid value ranges form `1,024` to `16,384`, the value must be a multiple of `1,024`.
  If not specified, the GPU function is disabled.

* `gpu_type` - (Optional, String) Specifies the GPU type to the function.
  Currently, only **nvidia-t4** is supported.

  -> If you want to use the `gpu_memory` and `gpu_type` parameters, the `runtime` parameter must be set to `Custom` or
  `Custom Image`. And you need to submit a service ticket to open this function, please refer to
  the documentation [how to submit a service ticket](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html).

The `func_mounts` block supports:

* `mount_type` - (Required, String) Specifies the mount type.
  + **sfs**
  + **sfsTurbo**
  + **ecs**

* `mount_resource` - (Required, String) Specifies the ID of the mounted resource (corresponding cloud service).

* `mount_share_path` - (Required, String) Specifies the remote mount path. Example: 192.168.0.12:/data.

* `local_mount_path` - (Required, String) Specifies the function access path.

<a name="functiongraph_custom_image"></a>
The `custom_image` block supports:

* `url` - (Required, String) Specifies the URL of SWR image, the URL must start with `swr.`.

* `command` - (Optional, String) Specifies the startup commands of the SWR image.
  Multiple commands are separated by commas (,). e.g. `/bin/sh`.
  If this parameter is not specified, the entrypoint or CMD in the image configuration will be used by default.

* `args` - (Optional, String) Specifies the command line arguments used to start the SWR image.
  If multiple arguments are separated by commas (,). e.g. `-args,value`.
  If this parameter is not specified, the CMD in the image configuration will be used by default.

* `working_dir` - (Optional, String) Specifies the working directory of the SWR image.
  If not specified, the default value is `/`.
  Currently, the folder path can only be set to `/` and it cannot be created or modified.

<a name="functiongraph_versions_management"></a>
The `versions` block supports:

* `name` - (Required, String) Specifies the version name.

* `aliases` - (Optional, List) Specifies the aliases management for specified version.
  The [object](#functiongraph_aliases_management) structure is documented below.

  -> A version can configure at most **one** alias.

<a name="functiongraph_aliases_management"></a>
The `aliases` block supports:

* `name` - (Required, String) Specifies the name of the version alias.

* `description` - (Optional, String) Specifies the description of the version alias.

<a name="functiongraph_reserved_instances"></a>
The `reserved_instances` block supports:

* `qualifier_type` - (Required, String) Specifies qualifier type of reserved instance. The valid values are as follows:
  + **version**
  + **alias**

  -> Reserved instances cannot be configured for both a function alias and the corresponding version. For example,
  if the alias of the `latest` version is `1.0` and reserved instances have been configured for this version,
  no more instances can be configured for alias `1.0`.

* `qualifier_name` - (Required, String) Specifies the version name or alias name.

* `count` - (Required, Int) Specifies the number of reserved instance.
  The valid value ranges from `0` to `1,000`.
  If this parameter is set to `0`, the reserved instance will not run.

* `idle_mode` - (Optional, Bool) Specifies whether to enable the idle mode. The default value is `false`.
  If this parameter is enabled, reserved instances are initialized and the mode change needs some time to take effect.
  You will still be billed at the price of reserved instances for non-idle mode in this period.

* `tactics_config` - (Optional, List) Specifies the auto scaling policies for reserved instance.
  The [tactics_config](#functiongraph_tactics_config) structure is documented below.

<a name="functiongraph_tactics_config"></a>
The `tactics_config` block supports:

* `cron_configs` - (Optional, List) Specifies the list of scheduled policy configurations.
  The [cron_configs](#functiongraph_cron_configs) structure is documented below.

* `metric_configs` - (Optional, List) Specifies the list of metric policy configurations.
  The [metric_configs](#functiongraph_metric_configs) structure is documented below.

  -> If you want to use the `metric_configs` parameter, please open a service ticket to enable this function. Refer to
  the documentation [how to submit a service ticket](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html).

<a name="functiongraph_cron_configs"></a>
The `cron_configs` block supports:

* `name` - (Required, String) Specifies the name of scheduled policy configuration.
  The valid length is limited from `1` to `60` characters, only letters, digits, hyphens (-), and underscores (_) are allowed.
  The name must start with a letter and ending with a letter or digit.

* `cron` - (Required, String) Specifies the cron expression.
  Set this parameter, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-functiongraph/functiongraph_01_0908.html).

* `count` - (Required, Int) Specifies the number of reserved instance to which the policy belongs.
  The valid value ranges from `0` to `1,000`.

  -> The number of reserved instances must be greater than or equal to the number of reserved instances in the basic configuration.

* `start_time` - (Required, Int) Specifies the effective timestamp of policy. The unit is `s`, e.g. **1740560074**.

* `expired_time` - (Required, Int) Specifies the expiration timestamp of the policy. The unit is `s`, e.g. **1740560074**.

<a name="functiongraph_metric_configs"></a>
The `metric_configs` block supports:

* `name` - (Required, String) Specifies the name of metric policy.
  The valid length is limited from `1` to `60` characters, only letters, digits, hyphens (-), and underscores (_) are allowed.
  The name must start with a letter and ending with a letter or digit.

* `type` - (Required, String) Specifies the type of metric policy.
  The valid value is as follows:
  + **Concurrency**: Reserved instance usage.

* `threshold` - (Required, Int) Specifies the metric policy threshold. The valid value ranges from `1` to `99`.

* `min` - (Required, Int) Specifies the minimun of traffic. The valid value ranges from `0` to `1,000`.

  -> The number of reserved instances must be greater than or equal to the number of reserved instances in the basic configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, comsist of `urn` and current `version`, the format is `<urn>:<version>`.

* `func_mounts/status` - The status of file system.

* `urn` - Uniform Resource Name.

* `version` - The version of the function.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Functions can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_fgs_function.test <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attributes are:
`app`, `func_code`, `agency`, `tags"`, `package`.
It is generally recommended running `terraform plan` after importing a function.
You can then decide if changes should be applied to the function, or the resource definition should be updated to align
with the function. Also you can ignore changes as below.

```hcl
resource "huaweicloud_fgs_function" "test" {
  ...
  lifecycle {
    ignore_changes = [
      app, func_code, agency, tags, package,
    ]
  }
}
```
