---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function"
description: |-
  Manages the function resource within HuaweiCloud.
---

# huaweicloud_fgs_function

Manages the function resource within HuaweiCloud.

~> Since version `1.73.1`, the requests of the function resource will send these parameters:<br>
   `enable_dynamic_memory`<br>
   `is_stateful_function`<br>
   `network_controller`<br>
   Since version `1.74.0`, the requests of the function resource will send these parameters:<br>
   `enable_auth_in_header`<br>
   `enable_class_isolation`<br>
   For the regions that do not support this parameter, please use the lower version to deploy this resource.

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

### Create function with a custom version and an alias for the latest version

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
      name        = "demo"
      description = "This is a description of the alias demo under the version latest."
    }
  }
  # The value of the parameter func_code must be modified before each custom version add.
  versions {
    name        = "v1.0"
    description = "This is a description of the version v1.0."

    aliases {
      name        = "v1_0-alias"
      description = "This is a description of the alias v1_0-alias under the version v1.0."
    }
  }
  versions {
    name        = "v2.0"
    description = "This is a description of the version v2.0."

    aliases {
      name        = "v2_0-alias"
      description = "This is a description of the alias v2_0-alias under the version v2.0."

      additional_version_weights = jsonencode({
        "v1.0": 15
      })
    }
  }
  versions {
    name        = "v3.0"
    description = "This is a description of the version v2.0."

    aliases {
      name        = "v3_0-alias"
      description = "This is a description of the alias v2_0-alias under the version v3.0."
      additional_version_strategy = jsonencode({
        "v2.0": {
          "combine_type": "or",
          "rules": [
            {
              "rule_type": "Header",
              "param": "version",
              "op": "=",
              "value": "v2_value"
            },
            {
              "rule_type": "Header",
              "param": "Owner",
              "op": "in",
              "value": "terraform,administrator"
            }
          ]
        }
      })
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

### With advanced configurations

```hcl
variable "function_name" {}
variable "function_codes" {}
variable "agency_name" {}
variable "trigger_access_vpc_ids" {
  type = list(string)
}

resource "huaweicloud_fgs_function" "test" {
  name                  = var.function_name
  app                   = "default"
  agency                = var.agency_name
  description           = "fuction test"
  handler               = "test.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = base64encode(var.function_codes)
  functiongraph_version = "v2"
  enable_dynamic_memory = true
  is_stateful_function  = true

  network_controller {
    disable_public_network = true

    dynamic "trigger_access_vpcs" {
      for_each = var.trigger_access_vpc_ids

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }
}
```

### Create function with Java runtime and corresponding configuration

```hcl
variable "function_name" {}
variable "agency_name" {}

resource "huaweicloud_fgs_function" "test" {
  name          = var.function_name
  memory_size   = 128
  runtime       = "Java11"
  timeout       = 15
  app           = "default"
  handler       = "com.huawei.demo.TriggerTests.apigTest"
  code_type     = "zip"
  code_filename = "java-demo.zip"
  agency        = var.agency_name

  enable_class_isolation = true
  ephemeral_storage      = 512
  heartbeat_handler      = "com.huawei.demo.TriggerTests.heartBeat"
  restore_hook_handler   = "com.huawei.demo.TriggerTests.restoreHook"
  restore_hook_timeout   = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the function is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the function.  
  The valid length is limited from `2` to `60` characters, only letters, digits, underscores (_) and hyphens (-) are
  allowed. The name must start with a letter, and end with a letter or a digit.
  Changing this will create a new resource.

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
  + **Node.js16.17**
  + **Node.js18.15**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Go1.x**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **Custom**
  + **PHP7.3**
  + **http**
  + **Custom Image**
  + **Cangjie1.0**

  Changing this will create a new resource.

* `timeout` - (Required, Int) Specifies the timeout interval of the function, in seconds.  
  The value ranges from `3` to `259,200`.

* `app` - (Required, String) Specifies the group to which the function belongs.

* `code_type` - (Required, String) Specifies the code type of the function.  
  The valid values are as follows:
  + **inline**: inline code.
  + **zip**: ZIP file.
  + **jar**: JAR file or java functions.
  + **obs**: function code stored in an OBS bucket.
  + **Custom-Image-Swr**: function code comes from the SWR custom image.

* `handler` - (Required, String) Specifies the entry point of the function.

-> If the function is created by an SWR image, keep `code_type` empty and use hyphen character (-) to set the handler.

* `description` - (Optional, String) Specifies the description of the function.

* `functiongraph_version` - (Optional, String, ForceNew) Specifies the version of the function framework.  
  The valid values are as follows:
  + **v1**: Hosts event-driven functions in a serverless context.
  + **v2**: Next-generation function hosting service powered by Huawei YuanRong architecture.

  Defaults to **v2**.  
  Changing this will create a new resource.

  -> Some regions support only **v1**, the default value is **v1**.

* `func_code` - (Optional, String) Specifies the function code.  
  The code value can be encoded using **Base64** or just with the text code.  
  Required if the `code_type` is set to **inline**, **zip**, or **jar**.

* `code_url` - (Optional, String) Specifies the URL where the function code is stored in OBS.  
  Required if the `code_type` is set to **obs**.

* `code_filename` - (Optional, String) Specifies the name of the function file.  
  Required if the `code_type` is set to **jar** or **zip**.

* `depend_list` - (Optional, List) Specifies the list of the dependency version IDs.

* `user_data` - (Optional, String) Specifies the key/value information defined for the function.  
  The key/value data might be parsed with [Terraform `jsonencode()` function]('https://www.terraform.io/docs/language/functions/jsonencode.html').

* `encrypted_user_data` - (Optional, String) Specifies the key/value information defined to be encrypted for the
  function.  
  The format is the same as `user_data`.

* `agency` - (Optional, String) Specifies the agency configuration of the function.  
  This parameter is mandatory if the function needs to access other cloud services.

* `app_agency` - (Optional, String) Specifies the execution agency enables you to obtain a token or an AK/SK for
  accessing other cloud services.

  -> After using this parameter, the function execution agency (`app_agency`) and the function configuration
     agency (`agency`) can be independently set, which can reduce unnecessary performance loss. Otherwise, the same
     agency is used for both function execution and function configuration.

* `initializer_handler` - (Optional, String) Specifies the initializer of the function.

* `initializer_timeout` - (Optional, Int) Specifies the maximum duration the function can be initialized, in seconds.  
  The valid value is range from `1` to `300`.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the
  function belongs.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the function belongs.

* `network_id` - (Optional, String) Specifies the network ID of subnet.

  -> An agency with VPC management permissions must be specified for the function.

* `dns_list` - (Optional, String) Specifies the private DNS configuration of the function network.  
  Private DNS list is associated to the function by a string in the following format:  
  `[{\"id\":\"ff8080828a07ffea018a17184aa310f5\","domain_name":"functiondebug.example1.com."}]`

  -> Ensure the agency with DNS management permissions specified before using this parameter.

* `mount_user_id` - (Optional, Int) Specifies the mount user ID.  
  The valid value is range from `–1` to `65,534`, except `0`.  
  Defaults to `-1`.

* `mount_user_group_id` - (Optional, Int) Specifies the mount user group ID.  
  The valid value is range from `–1` to `65,534`, except `0`.  
  Defaults to `-1`.

* `func_mounts` - (Optional, List) Specifies the list of function mount configurations.  
  The [func_mounts](#function_func_mounts) structure is documented below.

* `custom_image` - (Optional, List) Specifies the custom image configuration of the function.  
  The [custom_image](#function_custom_image) structure is documented below.  
  Required if the parameter `code_type` is **Custom-Image-Swr**.

* `max_instance_num` - (Optional, String) Specifies the maximum number of instances of the function.  
  The valid value is range from `-1` to `1,000`, defaults to `400`.
  + The minimum value is `-1` and means the number of instances is unlimited.
  + `0` means this function is disabled.
  + The empty value means to keep the default (latest updated) value.

  -> This parameter is only supported by the `v2` version of the function.

* `versions` - (Optional, List) Specifies the versions management of the function.  
  The [versions](#function_versions) structure is documented below.

  -> The value of the parameter `func_code`, `code_url` or `func_filename` must be modified before each custom version
     add.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the function.

* `enable_lts_log` - (Optional, Bool) Specifies whether to enable the LTS log, defaults to **false**.

* `log_group_id` - (Optional, String) Specifies the LTS group ID for collecting logs.

* `log_group_name` - (Optional, String) Specifies the LTS group name for collecting logs.

* `log_stream_id` - (Optional, String) Specifies the LTS stream IID for collecting logs.

* `log_stream_name` - (Optional, String) Specifies the LTS stream name for collecting logs.

-> If the `enable_lts_log` parameter is set to **true**, and the `log_group_id`, `log_group_name`, `log_stream_id` and
   `log_stream_name` parameters are not specified, the FunctionGraph will automatically create a log group and log stream.

* `reserved_instances` - (Optional, List) Specifies the reserved instance policies of the function.  
  The [reserved_instances](#function_reserved_instances) structure is documented below.

* `concurrency_num` - (Optional, Int) Specifies the number of concurrent requests of the function.  
  The valid value is range from `1` to `1,000`, the default value is `1`.

  -> 1. This parameter is only supported by the `v2` version of the function.
     <br>2. This parameter is available only when the `runtime` parameter is set to **http** or **Custom Image**.

* `gpu_type` - (Optional, String) Specifies the GPU type of the function.  
  Currently, only **nvidia-t4** is supported.

* `gpu_memory` - (Optional, Int) Specifies the GPU memory size allocated to the function, in MByte (MB).  
  The valid value is range form `1,024` to `16,384`, the value must be a multiple of `1,024`.  
  If not specified, the GPU function is disabled.

  ~> Submit a service ticket to open this function (GPU configuration), for the way please refer to
  the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html).

  -> If the `gpu_memory` and `gpu_type` configured, the `runtime` must be set to **Custom** or **Custom Image**.

* `enable_dynamic_memory` - (Optional, Bool) Specifies whether the dynamic memory configuration is enabled.  
  Defaults to **false**.

* `is_stateful_function` - (Optional, Bool) Specifies whether the function is a stateful function.  
  Defaults to **false**.

* `network_controller` - (Optional, List) Specifies the network configuration of the function.  
  The [network_controller](#function_network_controller) structure is documented below.

* `peering_cidr` - (Optional, String) Specifies the VPC cidr blocks used in the function code to detect whether it
  conflicts with the VPC cidr blocks used by the service.  
  The cidr blocks are separated by semicolons and cannot exceed `5`.

* `enable_auth_in_header` - (Optional, Bool) Specifies whether the authentication in the request header is enabled.  
  Defaults to **false**.

* `enable_class_isolation` - (Optional, Bool) Specifies whether the class isolation is enabled for the JAVA runtime
  functions.  
  Defaults to **false**.

  ~> Enabes class isolation can support Kafka dumping and improve class loading efficiency, but it may also cause some
     compatibility issues.

* `ephemeral_storage` - (Optional, Int) Specifies the size of the function ephemeral storage.  
  The valid values are as follows:
  + **512**
  + **10240**

  Defaults to `512`. Only custom image or http runtime supported.

* `heartbeat_handler` - (Optional, String) Specifies the heartbeat handler of the function.  
  The rule is **xx.xx**, such as **com.huawei.demo.TriggerTests.heartBeat**, it must contain periods (.).
  The heartbeat function entry must be in the same file as the function execution entry.

* `restore_hook_handler` - (Optional, String) Specifies the restore hook handler of the function.

* `restore_hook_timeout` - (Optional, Int) Specifies the timeout of the function restore hook.  
  The function will be forcibly stopped if the time is end.
  The valid value is range from `1` to `300`, the unit is seconds (s).

  -> Only Java runtime supports the configurations of the heartbeat and restore hook.

* `lts_custom_tag` - (Optional, Map) Specifies the custom tags configuration that used to filter the LTS logs.

  -> This parameter is only supported by the `v2` version of the function.

* `user_data_encrypt_kms_key_id` - (Optional, String) Specifies the KMS key ID for encrypting the user data.

* `code_encrypt_kms_key_id` - (Optional, String, ForceNew) Specifies the KMS key ID for encrypting the function code.  
  Changing this will create a new resource.

<a name="function_func_mounts"></a>
The `func_mounts` block supports:

* `mount_type` - (Required, String) Specifies the mount type.
  + **sfs**
  + **sfsTurbo**
  + **ecs**

* `mount_resource` - (Required, String) Specifies the ID of the mounted resource (corresponding cloud service).

* `mount_share_path` - (Required, String) Specifies the remote mount path, e.g. **192.168.0.12:/data**.

* `local_mount_path` - (Required, String) Specifies the function access path.

<a name="function_custom_image"></a>
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

<a name="function_versions"></a>
The `versions` block supports:

* `name` - (Required, String) Specifies the version name.  
  The valid length is limited from `1` to `42` characters, only letters, digits, underscores (_), hyphens (-) and
  periods (.) are allowed. The name must start and end with a letter or digit.

* `description` - (Optional, String) Specifies the description of the version.

  -> The **latest** version does not support configuration through this parameter, the root parameter `description` is
  the correct configuration parameter.

* `aliases` - (Optional, List) Specifies the aliases management for specified version.  
  The [aliases](#function_versions_aliases) structure is documented below.

  -> 1. A version can configure at most **one** alias.
     <br>2. A function can have a maximum of `10` aliases.

<a name="function_versions_aliases"></a>
The `aliases` block supports:

* `name` - (Required, String) Specifies the name of the version alias.  
  The valid length is limited from `1` to `63` characters, only letters, digits, underscores (_) and hyphens (-) are
  allowed. The name must start with a letter and end with a letter or digit.

* `description` - (Optional, String) Specifies the description of the version alias.

* `additional_version_weights` - (Optional, String) Specifies the percentage grayscale configuration of the version
  alias, in JSON format.

* `additional_version_strategy` - (Optional, String) Specifies the rule grayscale configuration of the version
  alias, in JSON format.

~> Only one of `additional_version_weights` and `additional_version_strategy` can be configured.

<a name="function_reserved_instances"></a>
The `reserved_instances` block supports:

* `qualifier_type` - (Required, String) Specifies the qualifier type of reserved instance.  
  The valid values are as follows:
  + **version**
  + **alias**

  -> Reserved instances cannot be configured for both a function alias and the corresponding version.
     <br>For example, if the alias of the `latest` version is `1.0` and reserved instances have been configured for this
     version, no more instances can be configured for alias `1.0`.

* `qualifier_name` - (Required, String) Specifies the version name or alias name.

* `count` - (Required, Int) Specifies the number of reserved instance.  
  The valid value is range from `0` to `1,000`.  
  If this parameter is set to `0`, the reserved instance will not run.

* `idle_mode` - (Optional, Bool) Specifies whether to enable the idle mode.  
  Defaults to **false**.  
  If this parameter is enabled, reserved instances are initialized and the mode change needs some time to take effect.  
  You will still be billed at the price of reserved instances for non-idle mode in this period.

* `tactics_config` - (Optional, List) Specifies the auto scaling policies for reserved instance.  
  The [tactics_config](#function_reserved_instances_tactics_config) structure is documented below.

<a name="function_reserved_instances_tactics_config"></a>
The `tactics_config` block supports:

* `cron_configs` - (Optional, List) Specifies the list of scheduled policy configurations.  
  The [cron_configs](#function_reserved_instances_tactics_config_cron_configs) structure is documented below.

* `metric_configs` - (Optional, List) Specifies the list of metric policy configurations.  
  The [metric_configs](#function_reserved_instances_tactics_metric_configs) structure is documented below.

  ~> Submit a service ticket to open this function (metric policy), for the way please refer to
  the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html).

<a name="function_reserved_instances_tactics_config_cron_configs"></a>
The `cron_configs` block supports:

* `name` - (Required, String) Specifies the name of scheduled policy configuration.  
  The valid length is limited from `1` to `60` characters, only letters, digits, hyphens (-), and underscores (_) are allowed.
  The name must start with a letter and ending with a letter or digit.

* `cron` - (Required, String) Specifies the cron expression.  
  For the syntax, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-functiongraph/functiongraph_01_0908.html).

* `count` - (Required, Int) Specifies the number of reserved instance to which the policy belongs.  
  The valid value is range from `0` to `1,000`.

  -> The number of reserved instances must be greater than or equal to the number of reserved instances in the basic configuration.

* `start_time` - (Required, Int) Specifies the effective timestamp of policy. The unit is `s`, e.g. **1740560074**.

* `expired_time` - (Required, Int) Specifies the expiration timestamp of the policy. The unit is `s`, e.g. **1740560074**.

<a name="function_reserved_instances_tactics_metric_configs"></a>
The `metric_configs` block supports:

* `name` - (Required, String) Specifies the name of metric policy.  
  The valid length is limited from `1` to `60` characters, only letters, digits, hyphens (-), and underscores (_) are
  allowed. The name must start with a letter and ending with a letter or digit.

* `type` - (Required, String) Specifies the type of metric policy.  
  The valid value is as follows:
  + **Concurrency**: Reserved instance usage.

* `threshold` - (Required, Int) Specifies the metric policy threshold.  
  The valid value is range from `1` to `99`.

* `min` - (Required, Int) Specifies the minimun of traffic.  
  The valid value is range from `0` to `1,000`.

  -> The number of reserved instances must be greater than or equal to the number of reserved instances in the basic configuration.

<a name="function_network_controller"></a>
The `network_controller` block supports:

* `trigger_access_vpcs` - (Optional, List) Specifies the configuration of the VPCs that can trigger the function.  
  The [trigger_access_vpcs](#function_network_controller_trigger_access_vpcs) structure is documented below.

* `disable_public_network` - (Optional, Bool) Specifies whether to disable the public network access.

<a name="function_network_controller_trigger_access_vpcs"></a>
The `trigger_access_vpcs` block supports:

* `vpc_id` - (Optional, String) Specifies the ID of the VPC that can trigger the function.

* `vpc_name` - (Optional, String) Specifies the name of the VPC that can trigger the function.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, comsist of `urn` and current `version`, the format is `<urn>:<version>`.

* `func_mounts` - The list of function mount configurations.  
  The [func_mounts](#function_func_mounts_attr) structure is documented below.

* `urn` - The URN (Uniform Resource Name) of the function.

* `version` - The version of the function.

<a name="function_func_mounts_attr"></a>
The `func_mounts` block supports:

* `status` - The mount status.

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
API response. The missing attributes are: `func_code`, `encrypted_user_data`, `tags`.
It is generally recommended running `terraform plan` after importing a function.
You can then decide if changes should be applied to the function, or the resource definition should be updated to align
with the function. Also you can ignore changes as below.

```hcl
resource "huaweicloud_fgs_function" "test" {
  ...
  lifecycle {
    ignore_changes = [
      app, func_code, tags,
    ]
  }
}
```
