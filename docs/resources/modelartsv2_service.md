---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_service"
description: |-
  Manages a service resource of the ModelArts service within HuaweiCloud.
---

# huaweicloud_modelartsv2_service

Manages a service resource of the ModelArts service within HuaweiCloud.

## Example Usage

### Creates an online service with a group configuration

```hcl
variable "subnet_id" {}
variable "resource_pool_id" {}
variable "image_storage_path_in_swr" {}

resource "huaweicloud_modelartsv2_service" "test" {
  name         = var.service_name
  version      = "0.0.1"
  type         = "REAL_TIME"
  description  = "Created by terraform script"
  workspace_id = "0"
  deploy_type  = "SINGLE"

  group_configs {
    framework = "COMMON"
    name      = "group-1"
    pool_id   = var.resource_pool_id
    weight    = 100
    count     = 2

    unit_configs {
      image {
        source   = "SWR"
        swr_path = var.image_storage_path_in_swr
      }
      custom_spec {
        memory = 1024
        cpu    = 1
      }

      cmd      = "sleep 20"
      count    = 1
      recovery = "Instance"

      envs = {
        foo = "bar"
      }
    }
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  runtime_config = jsonencode({
    service_invoke = {
      port                       = 9876
      protocol                   = "HTTPS"
      auth_type                  = "TOKEN"
      direct_channel_auth_enable = false
    }
    service_limit = {
      request_size_limit = 20
      request_timeout    = 30
      ip_white_list      = []
      ip_black_list      = []

      rate_limit = {
        num  = 200
        unit = "SECONDS"
      }
    }
  })
  upgrade_config = jsonencode({
    type = "ROLLING"
    rolling_update = {
      max_surge       = "50%"
      max_unavailable = "50%"
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the service is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the service.

* `version` - (Required, String) Specifies the version of the service.  
  The maximum valid length is `8`, only digits and dots (.) are allowed.  
  
  -> Specifies the old version number can switch the historical version.

* `type` - (Required, String, NonUpdatable) Specifies the reasoning method of the service.  
  The valid values are as follows:
  + **BATCH**: Online service, deploys the model as a Web Service, and provides online test UI and monitoring
    capabilities. The service keeps running.
  + **REAL_TIME**: Batch service, can perform inference on batch data and automatically stop after completing data
    processing.
  + **EDGE**: Edge service, deploys the model as a Web Service on the edge node through the IEF.

* `group_configs` - (Required, List) Specifies the instance group configurations of the service.  
  The [group_configs](#v2_service_group_configs) structure is documented below.  
  When the value of parameter `type` is **BATCH** or **EDGE**, only one group can be configured.  
  When the value of parameter `type` is **REAL_TIME**, multiple service instances can be configured and weights can be
  assigned according to business needs.

* `runtime_config` - (Required, String) Specifies the configuration of the service runtime, in JSON format.

* `upgrade_config` - (Required, String) Specifies the upgrade configuration of the service, in JSON format.

* `description` - (Optional, String) Specifies the description of the service.
  The maximum valid length is `100`, and cannot contain these special characters (`!><=&'"`).

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID of the service.

* `deploy_type` - (Optional, String, NonUpdatable) Specifies the deploy type of the service.  
  The valid values are as follows:
  + **SINGLE**
  + **MULTI**
  + **DIST**

* `log_configs` - (Optional, List) Specifies the log configurations of the service.  
  The [log_configs](#v2_service_log_configs) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the service.
  A maximum of `20` tags can be configured.

<a name="v2_service_group_configs"></a>
The `group_configs` block supports:

* `name` - (Required, String) Specifies the name of the instance group.  
  The valid length is limited from `1` to `64`, only English letters, Chinese characters, digits, hyphens (-) and
  underscores (_) are allowed.

* `count` - (Required, Int) Specifies the number of service instances in the deployment scenario.

* `weight` - (Required, Int) Specifies the weight percentage of the instance group.

  -> The sum of all group weights must be `100`.

* `unit_configs` - (Required, List) Specifies the unit configurations of the instance group.  
  The [unit_configs](#v2_service_unit_configs) structure is documented below.  
  When the unit is used for **SINGLE** deploy type, the length of `unit_configs` is `1`.  
  When used for **SINGLE** and **DIST** deploy type, the number of units configuration depending on the framework.

* `pool_id` - (Optional, String) Specifies the ID of the dedicated resource pool for the instance group.

* `framework` - (Optional, String) Specifies the algorithm framework.  
  The valid values are as follows:
  + **COMMON**
  + **VLLM**
  + **MINDIE**

~> Updates to the `name`, `pool_id` and `framework` parameters in existing group configurations are not supported, but
   new group configurations and remove existing group configurations are not limited.

<a name="v2_service_unit_configs"></a>
The `unit_configs` block supports:

* `image` - (Required, List) Specifies the image configuration of the group unit.  
  The [image](#v2_service_unit_config_image) structure is documented below.

* `role` - (Optional, String) Specifies the role of the group unit.  
  The valid values are as follows:
  + **SCHEDULER**: Scheduling unit, valid in the **MINDIE** framework.
  + **MANAGER**: Management unit, valid in the **MINDIE** framework.
  + **WORKER**: Work unit, valid in the **MINDIE** framework.
  + **PREFILL**: Total unit, valid in the **VLLM** framework.
  + **DECODE**: Incremental unit, valid in the **VLLM** framework.
  + **COMMON**: Others.

* `custom_spec` - (Optional, List) Specifies the configuration of the custom resource specification.  
  The [custom_spec](#v2_service_unit_config_custom_spec) structure is documented below.

* `flavor` - (Optional, String) Specifies the instance flavor of the group unit.

* `models` - (Optional, List) Specifies the model configuration of the group unit.  
  The [models](#v2_service_unit_config_models) structure is documented below.

* `codes` - (Optional, List) Specifies the code configuration of the group unit.  
  The [codes](#v2_service_unit_config_codes) structure is documented below.

* `count` - (Optional, Int) Specifies the instance number of the group unit.

* `cmd` - (Optional, String) Specifies the startup commands of the group unit.

* `envs` - (Optional, Map) Specifies the environment variables of the group unit.

* `readiness_health` - (Optional, List) Specifies the configuration of the readiness health check.  
  The [readiness_health](#v2_service_unit_config_health_check) structure is documented below.

* `startup_health` - (Optional, List) Specifies the configuration of the startup health check.  
  The [startup_health](#v2_service_unit_config_health_check) structure is documented below.

* `liveness_health` - (Optional, List) Specifies the configuration of the liveness health check.  
  The [liveness_health](#v2_service_unit_config_health_check) structure is documented below.

* `port` - (Optional, Int) Specifies the port of the group unit.

* `recovery` - (Optional, String) Specifies the recovery strategy of the group unit.  
  The valid values are as follows:
  + **INSTANCE_GROUP**
  + **INSTANCE**

~> Updates to the `role` parameter in existing unit configurations are not supported, but new unit configurations and
   remove existing unit configurations are not limited.

<a name="v2_service_unit_config_image"></a>
The `image` block supports:

* `source` - (Required, String) Specifies the image type of the group unit.
  The valid values are as follows:
  + **SWR**
  + **IMAGE**

* `swr_path` - (Required, String) Specifies the SWR storage path of the group unit.

* `id` - (Optional, String) Specifies the image ID of the group unit.  
  Only available if the value of parameter `source` is **IMAGE**.

<a name="v2_service_unit_config_custom_spec"></a>
The `custom_spec` block supports:

* `gpu` - (Optional, Float) Specifies the GPU number of the custom specification.
  The input value must be greater than `0` and support two decimal places (the third decimal place will be rounded off).

* `memory` - (Optional, Int) Specifies the memory size of the custom specification.

* `cpu` - (Optional, Float) Specifies the CPU number of the custom specification.
  The input value must be greater than `0` and support two decimal places (the third decimal place will be rounded off).

* `ascend` - (Optional, Int) Specifies the number of Ascend chips.  
  This parameter cannot be configured together with `gpu`.

<a name="v2_service_unit_config_models"></a>
The `models` block supports:

* `source` - (Required, String) Specifies the source type of the model configuration.  
  The valid values are as follows:
  + **OBS**
  + **OBSFS**
  + **EFS**

* `mount_path` - (Required, String) Specifies the path to mount into the container.  
  The value must start with a slash (/) and the path content allows letters, digits, hyphens (-), underscores (_),
  backslash (\\) and dots (.).

* `address` - (Optional, String) Specifies the source address of the model configuration.  
  This parameter is mutually exclusive with `source_id` and only required if the value of `source` is not **EFS**.

* `source_id` - (Optional, String) Specifies the source ID of the model configuration.  
  This parameter is mutually exclusive with `mount_path` and the SFS Turbo ID is required if the value of `source` is
  **EFS**.

<a name="v2_service_unit_config_codes"></a>
The `codes` block supports:

* `source` - (Required, String) Specifies the source type of the code configuration.  
  The valid values are as follows:
  + **OBS**
  + **OBSFS**
  + **EFS**
  + **GIT**

* `mount_path` - (Required, String) Specifies the path to mount into the container.  
  The value must start with a slash (/) and the path content allows letters, digits, hyphens (-), underscores (_),
  backslash (\\) and dots (.).

* `address` - (Optional, String) Specifies the source address of the code configuration.  
  This parameter is mutually exclusive with `source_id` and only required if the value of `source` is not **EFS**.

* `source_id` - (Optional, String) Specifies the source ID of the code configuration.  
  This parameter is mutually exclusive with `mount_path` and the SFS Turbo ID is required if the value of `source` is
  **EFS**.

<a name="v2_service_unit_config_health_check"></a>
The `readiness_health`, `startup_health` and `liveness_health` blocks support:

* `initial_delay_seconds` - (Required, Int) Specifies the time to wait when performing the first probe.  
  The minimum timeout value is `1`.

* `timeout_seconds` - (Required, Int) Specifies the timeout for executing the probe.  
  The minimum timeout value is `1`.

* `period_seconds` - (Required, Int) Specifies the period time for performing health check.  
  The minimum timeout value is `1`.

* `failure_threshold` - (Required, Int) Specifies the minimum number of consecutive detection failures.  
  The minimum timeout value is `1`.

* `check_method` - (Required, String) Specifies the method of the health check.  
  The valid values are as follows:
  + **EXEC**
  + **HTTP**

* `command` - (Optional, String) Specifies the commands configuration of the health check.  
  Only available if the `check_method` is **EXEC**.

* `url` - (Optional, String) Specifies the request URL of the health check.  
  Only available if the `check_method` is **HTTP**.

<a name="v2_service_log_configs"></a>
The `log_configs` block supports:

* `type` - (Required, String) Specifies the type of LTS configuration.
  Currently, the valid value is **STDOUT**.

* `log_group_id` - (Optional, String) Specifies the ID of the LTS group.

* `log_stream_id` - (Optional, String) Specifies the ID of the LTS stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `group_configs` - The instance group configurations of the service.  
  The [group_configs](#v2_service_group_configs_attr) structure is documented below.

* `status` - The status of the service.

* `predict_url` - The access addresses of the service.  
  The [predict_url](#v2_service_predict_url_attr) structure is documented below.

<a name="v2_service_group_configs_attr"></a>
The `group_configs` block supports:

* `unit_configs` - The unit configurations of the instance group.  
  The [unit_configs](#v2_service_unit_configs_attr) structure is documented below.

* `id` - The ID of the instance group.

<a name="v2_service_unit_configs_attr"></a>
The `unit_configs` block supports:

* `id` - The ID of the unit configuration.

<a name="v2_service_predict_url_attr"></a>
The `predict_url` block supports:

* `type` - The type of service access.

* `urls` - The URLs of service access.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Service can be imported using resource `id`, e.g.

```bash
$ terraform import huaweicloud_modelartsv2_service.test <id>
```
