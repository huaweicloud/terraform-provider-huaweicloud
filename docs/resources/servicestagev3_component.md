---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_component"
description: |-
  Manages a component resource within HuaweiCloud.
---

# huaweicloud_servicestagev3_component

Manages a component resource within HuaweiCloud.

## Example Usage

```hcl
variable "application_id" {}
variable "environment_id" {}
variable "component_name" {}
variable "ims_docker_image_url" {}
variable "associated_cce_cluster_id" {}
variable "associated_cse_engine_id" {}

resource "huaweicloud_servicestagev3_component" "test" {
  application_id = var.application_id
  environment_id = var.environment_id
  name           = var.component_name

  runtime_stack {
    deploy_mode = "container"
    name        = "Docker"
    type        = "Docker"
  }

  source = jsonencode({
    "auth": "iam",
    "kind": "image",
    "storage": "swr",
    "url": var.ims_docker_image_url
  })

  version = "1.0.1"
  replica = 2

  refer_resources {
    id         = var.associated_cce_cluster_id
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
    })
  }
  refer_resources {
    id   = var.associated_cse_engine_id
    type = "cse"
  }

  tags = {
    foo = "bar"
  }

  description    = "Created by terraform script"
  limit_cpu      = 0.25
  limit_memory   = 0.5
  request_cpu    = 0.25
  request_memory = 0.5

  envs {
    name  = "env_name"
    value = "env_value"
  }

  storages {
    type       = "HostPath"
    name       = "%[2]s"
    parameters = jsonencode({
      "default_mode": 0,
      "path": "/tmp"
    })
    mounts {
      path      = "/category"
      sub_path  = "sub"
      read_only = false
    }
  }

  command = jsonencode({
    "args": ["-a"],
    "command": ["ls"]
  })

  post_start {
    command = ["test"]
    type    = "command"
  }

  pre_stop {
    command = ["test"]
    type    = "command"
  }

  mesher {
    port = 60
  }

  timezone = "Asia/Shanghai"

  logs {
    log_path         = "/tmp"
    rotate           = "Hourly"
    host_path        = "/tmp"
    host_extend_path = "PodName"
  }

  custom_metric {
    path       = "/tmp"
    port       = 600
    dimensions = "cpu_usage,mem_usage"
  }

  affinity {
    condition = "required"
    kind      = "node"
    match_expressions {
      key       = "affinity1"
      value     = "foo"
      operation = "In"
    }
    weight = 100
  }
  affinity {
    condition = "preferred"
    kind      = "node"
    match_expressions {
      key       = "affinity2"
      value     = "bar"
      operation = "NotIn"
    }
    weight = 1
  }

  anti_affinity {
    condition = "required"
    kind      = "pod"
    match_expressions {
      key       = "anit-affinity1"
      operation = "Exists"
    }
    weight = 100
  }
  anti_affinity {
    condition = "preferred"
    kind      = "pod"
    match_expressions {
      key       = "anti-affinity2"
      operation = "DoesNotExist"
    }
    weight = 1
  }

  liveness_probe {
    type    = "tcp"
    delay   = 30
    timeout = 30
    port    = 800
  }

  readiness_probe {
    type    = "http"
    delay   = 30
    timeout = 30
    scheme  = "HTTPS"
    host    = "127.0.0.1"
    port    = 8000
    path    = "/v1/test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the component is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the application ID to which the component belongs.  
  Changing this will create a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the environment ID where the component is deployed.  
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the component.  
  The valid length is limited from `2` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter and end with a letter or a digit.  
  Changing this will create a new resource.

* `runtime_stack` - (Required, List, ForceNew) Specifies the configuration of the runtime stack.  
  The [runtime_stack](#servicestage_v3_component_runtime_stack) structure is documented below.  
  Changing this will create a new resource.

* `source` - (Required, String) Specifies the source configuration of the component, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__en-us_topic_0220056058_ref28944532).

* `version` - (Required, String) Specifies the version of the component.  
  The format is **{number}.{number}.{number}** or **{number}.{number}.{number}.{number}**, e.g. **1.0.1**.

* `refer_resources` - (Required, List) Specifies the configuration of the reference resources.  
  The [refer_resources](#servicestage_v3_component_refer_resources) structure is documented below.

* `config_mode` - (Optional, String) Specifies the configuration mode of the component.
  The valid values are as follows:
  + **yaml**

* `workload_content` - (Optional, String) Specifies the workload content of the component.

* `description` - (Optional, String) Specifies the description of the component.  
  The value can contain a maximum of `128` characters.

  -> The value of the `description` cannot be set to empty value by updating.

* `build` - (Optional, String) Specifies the build configuration of the component, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__en-us_topic_0220056060_table7559740).

* `replica` - (Optional, Int, ForceNew) Specifies the replica number of the component.  
  Changing this will create a new resource.

* `limit_cpu` - (Optional, Float) Specifies the maximum number of the CPU limit.  
  The unit is **Core**.

* `limit_memory` - (Optional, Float) Specifies the maximum number of the memory limit.  
  The unit is **GiB**.

* `request_cpu` - (Optional, Float) Specifies the number of the CPU request resources.  
  The unit is **Core**.

* `request_memory` - (Optional, Float) Specifies the number of the memory request resources.  
  The unit is **GiB**.

* `envs` - (Optional, List) Specifies the configuration of the environment variables.  
  The [envs](#servicestage_v3_component_envs) structure is documented below.

* `storages` - (Optional, List) Specifies the storage configuration.  
  The [storages](#servicestage_v3_component_storages) structure is documented below.

* `deploy_strategy` - (Optional, List) Specifies the configuration of the deploy strategy.  
  The [deploy_strategy](#servicestage_v3_component_deploy_strategy) structure is documented below.

* `update_strategy` - (Optional, String) Specifies the configuration of the update strategy, in JSON format.

* `command` - (Optional, String) Specifies the start commands of the component, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table856311795212).

* `post_start` - (Optional, List) Specifies the post start configuration.  
  The [post_start](#servicestage_v3_component_lifecycle) structure is documented below.

* `pre_stop` - (Optional, List) Specifies the pre stop configuration.  
  The [pre_stop](#servicestage_v3_component_lifecycle) structure is documented below.

* `mesher` - (Optional, List) Specifies the configuration of the access mesher.  
  The [mesher](#servicestage_v3_component_mesher) structure is documented below.

* `timezone` - (Optional, String) Specifies the time zone in which the component runs, e.g. **Asia/Shanghai**.

* `jvm_opts` - (Optional, String) Specifies the JVM parameters of the component. e.g. **-Xms256m -Xmx1024m**.  
  If there are multiple parameters, separate them by spaces.  
  If this parameter is left blank, the default value is used.

* `tomcat_opts` - (Optional, String) Specifies the configuration of the tomcat server, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table2836191954317).

* `logs` - (Optional, List) Specifies the configuration of the logs collection.  
  The [logs](#servicestage_v3_component_logs) structure is documented below.

* `custom_metric` - (Optional, List) Specifies the configuration of the monitor metric.  
  The [custom_metric](#servicestage_v3_component_custom_metric) structure is documented below.

* `affinity` - (Optional, List) Specifies the affinity configuration of the component.  
  The [affinity](#servicestage_v3_component_affinity) structure is documented below.

* `anti_affinity` - (Optional, List) Specifies the anti-affinity configuration of the component.  
  The [anti_affinity](#servicestage_v3_component_affinity) structure is documented below.

* `liveness_probe` - (Optional, List) Specifies the liveness probe configuration of the component.  
  The [liveness_probe](#servicestage_v3_component_probe) structure is documented below.

* `readiness_probe` - (Optional, List) Specifies the readiness probe configuration of the component.  
  The [readiness_probe](#servicestage_v3_component_probe) structure is documented below.

* `external_accesses` - (Optional, List) Specifies the configuration of the external accesses.  
  The [external_accesses](#servicestage_v3_component_external_accesses) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the component.

<a name="servicestage_v3_component_runtime_stack"></a>
The `runtime_stack` block supports:

* `name` - (Required, String, ForceNew) Specifies the stack name.  
  Changing this will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the stack type.  
  Changing this will create a new resource.

* `deploy_mode` - (Required, String, ForceNew) Specifies the deploy mode of the stack.  
  Changing this will create a new resource.

* `version` - (Optional, String, ForceNew) Specifies the stack version.  
  Changing this will create a new resource.

<a name="servicestage_v3_component_refer_resources"></a>
The `refer_resources` block supports:

* `id` - (Required, String) Specifies the resource ID.

* `type` - (Required, String) Specifies the resource type.

* `parameters` - (Optional, String) Specifies the resource parameters, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table838321632514).

<a name="servicestage_v3_component_envs"></a>
The `envs` block supports:

* `name` - (Required, String) Specifies the name of the environment variable.

* `value` - (Optional, String) Specifies the value of the environment variable.

<a name="servicestage_v3_component_storages"></a>
The `storages` block supports:

* `type` - (Required, String) Specifies the type of the data storage.
  + **HostPath**: Host path for local disk mounting.
  + **EmptyDir**: Temporary directory for local disk mounting.
  + **ConfigMap**: Configuration item for local disk mounting.
  + **Secret**: Secrets for local disk mounting.
  + **PersistentVolumeClaim**: Cloud storage mounting.

* `name` - (Required, String) Specifies the name of the disk where the data is stored.  
  Only lowercase letters, digits, and hyphens (-) are allowed and must start and end with a lowercase letter or digit.

* `parameters` - (Required, String) Specifies the information corresponding to the specific types of data storage,
  in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table16441247172510).

* `mounts` - (Required, List) Specifies the configuration of the disk mounts.  
  The [mounts](#servicestage_v3_component_storage_mounts) structure is documented below.

<a name="servicestage_v3_component_storage_mounts"></a>
The `mounts` block supports:

* `path` - (Required, String) Specifies the mount path.

* `sub_path` - (Required, String) Specifies the sub mount path.

* `read_only` - (Required, Bool) Specifies whether the disk mount is read-only.

<a name="servicestage_v3_component_deploy_strategy"></a>
The `deploy_strategy` block supports:

* `type` - (Required, String) Specifies the deploy type.
  + **OneBatchRelease**: Single-batch upgrade.
  + **RollingRelease**: Rolling deployment and upgrade.
  + **GrayRelease**: Dark launch upgrade.

* `rolling_release` - (Optional, String) Specifies the rolling release parameters, in JSON format.  
  Required if the `type` is **RollingRelease**.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table4696103920).

* `gray_release` - (Optional, String) Specifies the gray release parameters, in JSON format.  
  Required if the `type` is **GrayRelease**.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table888818707).

<a name="servicestage_v3_component_lifecycle"></a>
The `post_start` and `pre_stop` blocks support:

* `type` - (Required, String) Specifies the processing method.
  + **http**
  + **command**

* `scheme` - (Optional, String) Specifies the HTTP request type.
  + **HTTP**
  + **HTTPS**

  This parameter is only available when the `type` is set to `http`.

* `host` - (Optional, String) Specifies the host (IP) of the lifecycle configuration.  
  If this parameter is left blank, the pod IP address is used.  
  This parameter is only available when the `type` is set to `http`.

* `port` - (Optional, String) Specifies the port number of the lifecycle configuration.
  This parameter is only available when the `type` is set to `http`.

* `path` - (Optional, String) Specifies the request path of the lifecycle configuration.
  This parameter is only available when the `type` is set to `http`.

* `command` - (Optional, List) Specifies the command list of the lifecycle configuration.
  This parameter is only available when the `type` is set to `command`.

<a name="servicestage_v3_component_mesher"></a>
The `mesher` block supports:

* `port` - (Required, Int) Specifies the process listening port.

<a name="servicestage_v3_component_logs"></a>
The `logs` block supports:

* `log_path` - (Required, String) Specifies the log path of the container, e.g. **/tmp**.

* `rotate` - (Required, String) Specifies the interval for dumping logs.
  + **Hourly**
  + **Daily**
  + **Weekly**

* `host_path` - (Required, String) Specifies the mounted host path, e.g. **/tmp**.

* `host_extend_path` - (Required, String) Specifies the extension path of the host.
  + **None**: the extended path is not used.
  + **PodUID**: extend the host path based on the pod ID.
  + **PodName**: extend the host path based on the pod name.
  + **PodUID/ContainerName**: extend the host path based on the pod ID and container name.
  + **PodName/ContainerName**: extend the host path based on the pod name and container name.

<a name="servicestage_v3_component_custom_metric"></a>
The `custom_metric` block supports:

* `path` - (Required, String) Specifies the collection path, such as **./metrics**.

* `port` - (Required, Int) Specifies the collection port, such as **9090**.

* `dimensions` - (Required, String) Specifies the monitoring dimension, such as **cpu_usage**, **mem_usage** or
  **cpu_usage,mem_usage** (separated by a comma).

<a name="servicestage_v3_component_affinity"></a>
The `affinity` and `anti_affinity` blocks support:

* `condition` - (Required, String) Specifies the condition type of the (anti) affinity rule.

* `kind` - (Required, String) Specifies the kind of the (anti) affinity rule.

* `match_expressions` - (Required, List) Specifies the list of the match rules for (anti) affinity.  
  The [match_expressions](#servicestage_v3_component_affinity_match_expressions) structure is documented below.

* `weight` - (Optional, Int) Specifies the weight of the (anti) affinity rule.  
  The valid value is range from `1` to `100`.

<a name="servicestage_v3_component_affinity_match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Required, String) Specifies the key of the match rule.

* `operation` - (Required, String) Specifies the operation of the match rule.

* `value` - (Required, String) Specifies the value of the match rule.

<a name="servicestage_v3_component_probe"></a>
The `liveness_probe` and `readiness_probe` blocks support:

* `type` - (Required, String) Specifies the type of the probe.
  + **http**
  + **tcp**
  + **command**

* `delay` - (Required, Int) Specifies the delay time of the probe.

* `timeout` - (Required, Int) Specifies the timeout of the probe.

* `scheme` - (Optional, String) Specifies the scheme type of the probe.
  + **HTTP**
  + **HTTPS**

  This parameter is only available when the `type` is set to `http`.

* `host` - (Optional, String) Specifies the host of the probe.  
  Defaults to pod ID, also custom IP address can be specified.  
  This parameter is only available when the `type` is set to `http`.

* `port` - (Optional, Int) Specifies the port of the probe.  
  This parameter is only available when the `type` is set to `tcp` or `http`.

* `path` - (Optional, String) Specifies the path of the probe.  
  This parameter is only available when the `type` is set to `http`.

* `command` - (Optional, List) Specifies the command list of the probe.  
  This parameter is only available when the `type` is set to `command`.

<a name="servicestage_v3_component_external_accesses"></a>
The `external_accesses` block supports:

* `protocol` - (Required, String) Specifies the protocol of the external access.

* `address` - (Optional, String) Specifies the address of the external access.

* `forward_port` - (Optional, Int) Specifies the forward port of the external access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `status` - The status of the component.
  + **RUNNING**
  + **PENDING**

* `created_at` - The creation time of the component, in RFC3339 format.

* `updated_at` - The latest update time of the component, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 5 minutes.

## Import

Components can be imported using `application_id` and `id` separated by a slash e.g.

```bash
$ terraform import huaweicloud_servicestagev3_component.test <application_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to attributes missing from the API
response, security or some other reason.
The missing attribute is `workload_content`, `tags`.
It is generally recommended running `terraform plan` after importing resource.
You can decide if changes should be applied to resource, or the definition should be updated to align with the resource.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_servicestagev3_component" "test" {
  ...

  lifecycle {
    ignore_changes = [
      workload_content, tags,
    ]
  }
}
```
