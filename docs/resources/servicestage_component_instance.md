---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_component_instance"
description: ""
---

# huaweicloud_servicestage_component_instance

This resource is used to deploy a component under specified application within HuaweiCloud ServiceStage service.

## Example Usage

### Deploy a component in the container with specified SWR image

```hcl
variable "app_id" {}
variable "component_id" {}
variable "env_id" {}
variable "instance_name" {}
variable "flavor_id" {}
variable "component_name" {}
variable "swr_image_url" {}
variable "cce_cluster_id" {}
variable "cse_engine_id" {}

resource "huaweicloud_servicestage_component_instance" "default" {
  application_id = var.app_id
  component_id   = var.component_id
  environment_id = var.env_id

  name      = var.instance_name
  version   = "1.0.0"
  replica   = 1
  flavor_id = var.flavor_id

  artifact {
    name      = var.component_name
    type      = "image"
    storage   = "swr"
    url       = var.swr_image_url
    auth_type = "iam"
  }

  refer_resource {
    type = "cce"
    id   = var.cce_cluster_id

    parameters = {
      type      = "VirtualMachine"
      namespace = "default"
    }
  }

  refer_resource {
    type = "cse"
    id   = var.cse_engine_id
  }

  configuration {
    env_variable {
      name  = "TZ"
      value = "Asia/Shanghai"
    }

    log_collection_policy {
      host_path = "/tmp"

      container_mounting {
        path         = "/attached/01"
        aging_period = "Hourly"
      }
    }
  }
}
```

### Deploy a component in the ECS instance with specified jar package

```hcl
variable "app_id" {}
variable "component_id" {}
variable "env_id" {}
variable "instance_name" {}
variable "flavor_id" {}
variable "component_name" {}
variable "jar_url" {}
variable "obs_bucket_name" {}
variable "obs_bucket_endpoint" {}
variable "obs_object_key" {}
variable "ecs_instance_id" {}

resource "huaweicloud_servicestage_component_instance" "test" {
  application_id = var.app_id
  component_id   = var.component_id
  environment_id = var.env_id

  name      = var.instance_name
  version   = "1.0.0"
  replica   = 1
  flavor_id = var.flavor_id

  artifact {
    name      = var.component_name
    auth_type = "iam"
    type      = "package"
    storage   = "obs"
    url       = var.jar_url

    properties {
      bucket   = var.obs_bucket_name
      endpoint = var.obs_bucket_endpoint
      key      = var.obs_object_key
    }
  }

  refer_resource {
    type = "ecs"
    id   = "Default"

    parameters = {
      hosts = "[\"${var.ecs_instance_id}\"]"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create (deploy) the ServiceStage (component) instance.
  If omitted, the provider-level region will be used. Changing this will create a new instance.

* `application_id` - (Required, String, ForceNew) Specifies the application ID to which the instance belongs.
  Changing this will create a new instance.

* `component_id` - (Required, String, ForceNew) Specifies the component ID to build (deploy).
  Changing this will create a new instance.

* `environment_id` - (Required, String, ForceNew) Specifies the environment ID in which the component to build (deployed).
  Changing this will create a new instance.

* `name` - (Required, String, ForceNew) Specifies the instance name.
  The name can contain `2` to `63` characters, only lowercase letters, digits and hyphens (-) are allowed.
  The name must start with a lowercase letter and end with a lowercase letter or digit.
  Changing this will create a new instance.

* `version` - (Required, String) Specifies the application component version that meets version semantics.
  For example: `1.0.0`.

* `replica` - (Required, Int, ForceNew) Specifies the number of instance replicas.
  Changing this will create a new instance.

* `flavor_id` - (Required, String) Specifies the resource specifications, which can be obtained by using data source or
  the customize resource specifications.
  The format of customize resource specifications is **CUSTOM-xxG:xxC-xxC:xxGi-xxGi**.
  The meaning of each part is:
    + **xxG**: storage capacity allocated to a component instance (reserved field). You can set it to a fixed number.
    + **xxC-xxC**: the maximum and minimum number of CPU cores allocated to a component instance.
    + **xxGi-xxGi**: the maximum and minimum memory allocated to a component instance.

  For example, **CUSTOM-10G:0.5C-0.25C:1.6Gi-0.8Gi** indicates the maximum number of CPU cores allocated to a
  component instance is 0.5, the minimum number of CPU cores is 0.25, the maximum memory is 1.6 Gi, and the minimum
  memory is 0.8 Gi.

* `refer_resource` - (Required, List) Specifies the deployed resources.
  The [object](#servicestage_refer_resource) structure is documented below.

* `artifact` - (Optional, List) Specifies the component artifact settings.
  The key indicates the component name. In the Docker container scenario, the key indicates the container name.
  The [object](#servicestage_artifact) structure is documented below.

-> If the source parameters of the component resource specify the software package source, this parameter is optional,
  and the software package source of the component is inherited by default. Otherwise, this parameter is required.

* `description` - (Optional, String) Specifies the description of the instance.
  The description can contain a maximum of `128` characters.

* `configuration` - (Optional, List) Specifies the configuration parameters, such as environment variables,
  deployment configurations, and O&M monitoring.
  The [object](#servicestage_configuration) structure is documented below.

* `external_access` - (Optional, List) Specifies the configuration of the external network access.
  The [object](#servicestage_external_access) structure is documented below.

<a name="servicestage_refer_resource"></a>
The `refer_resource` block supports:

* `type` - (Required, String) Specifies the resource type.
  The basic resources include:
  + **cce**: Cloud Container Engine (CCE)
  + **cci**: Cloud Container Instance (CCI)
  + **ecs**: Elastic Cloud Server (ECS).
  + **as**: Auto Scaling (AS)

  The Optional resources include:
  + **rds**: Relational Database Service (RDS)
  + **dcs**: Distributed Cache Service (DCS),
  + **elb**: Elastic Load Balance (ELB)

  For other resource types, please refer to the environment documentation.

* `id` - (Required, String) Specifies the resource ID.
  If the `type` is set to **ecs**, the value of this parameter must be **Default**.

* `alias` - (Optional, String) Specifies the application alias, which is provided only in DCS scenario.
  The valid values are: **distributed_session**, **distributed_cache** and **distributed_session, distributed_cache**.
  Defaults to **distributed_session, distributed_cache**.

* `parameters` - (Optional, Map) Specifies the reference resource parameter.
  + When `type` is set to **cce**, this parameter is mandatory, and need to specify the namespace of the cluster where
  the component is to be deployed, such as **{"namespace": "default"}**.
  + When `type` is set to **ecs**, this parameter is mandatory, and need to specify the hosts where the component is to
  be deployed, such as **{"hosts":"[\"04d9f887-9860-4029-91d1-7d3102903a69\", \"04d9f887-9860-4029-91d1-7d3102903a70\"]"}**.

<a name="servicestage_artifact"></a>
The `artifact` block supports:

* `name` - (Required, String) Specifies the component name.
  But for **Docker container scenario**, this name is the container name.

* `type` - (Required, String) Specifies the source type.
  The valid values are **package** (VM-based deployment) and **image** (container-based deployment).

* `storage` - (Required, String) Specifies the storage mode. The valid values are **swr** and **obs**.

* `url` - (Required, String) Specifies the software package or image address.
  For a component deployed on a VM, this parameter is the software package address.
  For a component deployed based on a container, this parameter is the image address or component name:v${index}.
  The latter indicates that the component source code or the image automatically built using the software package
  will be used.

* `auth_type` - (Optional, String) Specifies the authentication mode.
  The valid values are **iam** and **none**. Defaults to **iam**.

* `version` - (Optional, String) Specifies the version number.

* `properties` - (Optional, List) Specifies the properties of the OBS object.
  This parameter is available only `storage` is **obs**.
  The [object](#servicestage_properties) structure is documented below.

<a name="servicestage_properties"></a>
The `properties` block supports:

* `bucket` - (Optional, String) Specifies the OBS bucket name.

* `endpoint` - (Optional, String) Specifies the OBS bucket endpoint.

* `key` - (Optional, String) Specifies the key name of the OBS object.

<a name="servicestage_configuration"></a>
The `configuration` block supports:

* `env_variable` - (Optional, List) Specifies the environment variables.
  The [object](#servicestage_env_variables) structure is documented below.

* `storage` - (Optional, List) Specifies the data storage configuration.
  The [object](#servicestage_storages) structure is documented below.

* `strategy` - (Optional, List) Specifies the upgrade policy.
  The [object](#servicestage_strategy) structure is documented below.

* `lifecycle` - (Optional, List) Specifies the lifecycle.
  The [object](#servicestage_lifecycle) structure is documented below.

* `log_collection_policy` - (Optional, List) Specifies the policies of the log collection.
  The [object](#servicestage_log_collection_policies) structure is documented below.

* `scheduler` - (Optional, List) Specifies the scheduling policy.
  The key indicates the component name. In the Docker container scenario, key indicates the container name.
  If the source parameters of a component specify the software package source, this parameter is optional, and the
  software package source of the component is inherited by default. Otherwise, this parameter is required.
  The [object](#servicestage_scheduler) structure is documented below.

* `probe` - (Optional, List) Specifies the variable value.
  The [object](#servicestage_probe) structure is documented below.

<a name="servicestage_env_variables"></a>
The `env_variable` block supports:

* `name` - (Required, String) Specifies the variable name.
  The name can contain `1` to `64` characters, only letters, digits, hyphens (-), underscores (_) and dots (.) are
  allowed. The name cannot start with a digit.

* `value` - (Required, String) Specifies the variable value.

<a name="servicestage_storages"></a>
The `storage` block supports:

* `type` - (Required, String) Specifies the variable name.
  The valid values are as follows:
  + **HostPath**: host path mounting.
  + **EmptyDir**: temporary directory mounting.
  + **ConfigMap**: configuration item mounting.
  + **Secret**: secret volume mounting.
  + **PersistentVolumeClaim**: cloud storage mounting.

* `parameter` - (Required, List) Specifies the storage parameters.
  The [object](#servicestage_storage_parameters) structure is documented below.

* `mount` - (Required, List) Specifies the directory mounted to the container.
  The [object](#servicestage_storage_mounts) structure is documented below.

<a name="servicestage_storage_parameters"></a>
The `parameter` block supports:

* `path` - (Optional, String) Specifies the host path. Required if the storage `type` is **HostPath**.

* `name` - (Optional, String) Specifies the configuration item.

* `claim_name` - (Optional, String) Specifies the PVC name.

* `secret_name` - (Optional, String) Specifies the Secret name. Required if the storage `type` is **Secret**.

<a name="servicestage_storage_mounts"></a>
The `mount` block supports:

* `path` - (Required, String) Specifies the mounted disk path.

* `readonly` - (Required, Bool) Specifies the mounted disk permission is read-only or read-write.
  + **true**: read-only.
  + **false**: read-write.

* `subpath` - (Optional, String) Specifies the subpath of the mounted disk.
  This parameter is applicable to `http` type.

<a name="servicestage_strategy"></a>
The `strategy` block supports:

* `upgrade` - (Optional, String) Specifies the upgrade policy.
  The valid values are **Recreate** or **RollingUpdate**. The default value is **RollingUpdate**.
  The **Recreate** indicates in-place upgrade while the **RollingUpdate** indicates rolling upgrade.

<a name="servicestage_lifecycle"></a>
The `lifecycle` block supports:

* `entrypoint` - (Optional, List) Specifies the startup commands.
  The [object](#servicestage_entrypoint) structure is documented below.

* `post_start` - (Optional, List) Specifies the post-start processing.
  The [object](#servicestage_lifecycle_process) structure is documented below.

* `pre_stop` - (Optional, List) Specifies the pre-stop processing.
  The [object](#servicestage_lifecycle_process) structure is documented below.

<a name="servicestage_log_collection_policies"></a>
The `log_collection_policy` block supports:

* `container_mounting` - (Required, List) Specifies the configurations of the container mounting.
  The [object](#servicestage_container_mounting) structure is documented below.

* `host_path` - (Optional, String) Specifies the The host path that will be mounted to the specified container path.

<a name="servicestage_container_mounting"></a>
The `container_mounting` block supports:

* `path` - (Required, String) Specifies the path of the container mounting.

* `host_extend_path` - (Optional, String) Specifies the extended host path.
  This parameter can be configured only when `host_path` is configured.
  The valid values are as follows:
  + **PodUID**
  + **PodName**
  + **PodUID/ContainerName**
  + **PodName/ContainerName**

* `aging_period` - (Optional, String) Specifies the aging period.
  The valid values are **Hourly**, **Daily** and **Weekly**. The default value is **Hourly**.

<a name="servicestage_entrypoint"></a>
The `entrypoint` block supports:

* `commands` - (Required, List) Specifies the commands.

* `args` - (Required, List) Specifies the running parameters.

<a name="servicestage_lifecycle_process"></a>
The `post_start` and `pre_stop` block supports:

* `type` - (Required, List) Specifies the process type. The valid values are **command** and **http**.

* `parameters` - (Required, List) Specifies the start post-processing or stop pre-processing parameters.
  The [object](#servicestage_process_param) structure is documented below.

<a name="servicestage_process_param"></a>
The `parameters` block supports:

* `commands` - (Optional, List) Specifies the commands, such as **["sleep", "1"]**.
  This parameter is required if process type is **command**, and it is applicable to **command** type.

* `host` - (Optional, String) Specifies the custom IP address. The default address is pod IP address.
  This parameter is required if process type is **http**, and it is applicable to **http** type.

* `port` - (Optional, Int) Specifies the port number.
  This parameter is required if process type is **http**, and it is applicable to **http** type.

* `path` - (Optional, String) Specifies the request URL.
  This parameter is required if process type is **http**, and it is applicable to **http** type.

<a name="servicestage_scheduler"></a>
The `scheduler` block supports:

* `affinity` - (Optional, List) Specifies the commands.
  The [object](#servicestage_affinity) structure is documented below.

* `anti_affinity` - (Optional, List) Specifies the commands.
  The [object](#servicestage_affinity) structure is documented below.

<a name="servicestage_affinity"></a>
The `affinity` and `anti_affinity` block supports:

* `availability_zones` - (Optional, List) Specifies the AZ list.

* `private_ips` - (Optional, List) Specifies the node private IP address list.

* `instance_names` - (Optional, List) Specifies the list of component instance names.

<a name="servicestage_probe"></a>
The `probe` block supports:

* `liveness` - (Optional, List) Specifies the component liveness probe.
  The [object](#servicestage_probe_detail) structure is documented below.

* `readiness` - (Optional, List) Specifies the component service probe.
  The [object](#servicestage_probe_detail) structure is documented below.

<a name="servicestage_probe_detail"></a>
The `liveness` and `readiness` block supports:

* `type` - (Required, String) Specifies the probe type. The valid values are as follows:
  + **command**: command execution check.
  + **http**: HTTP request check.
  + **tcp**: TCP port check.

* `command_param` - (Optional, List) Specifies the commands. Required if `type` is **command**.
  The [object](#servicestage_command_param) structure is documented below.

* `http_param` - (Optional, List) Specifies the commands. Required if `type` is **http**.
  The [object](#servicestage_http_param) structure is documented below.

* `tcp_param` - (Optional, List) Specifies the commands. Required if `type` is **tcp**.
  The [object](#servicestage_tcp_param) structure is documented below.

* `delay` - (Optional, Int) Specifies the interval between the startup and detection.

* `timeout` - (Optional, Int) Specifies the detection timeout interval.

<a name="servicestage_command_param"></a>
The `command_param` block supports:

* `commands` - (Required, List) Specifies the command list.

<a name="servicestage_http_param"></a>
The `http_param` block supports:

* `scheme` - (Required, String) Specifies the protocol scheme. The valid values are **HTTP** and **HTTPS**.

* `port` - (Required, Int) Specifies the port number.

* `path` - (Required, String) Specifies the request path.

* `host` - (Optional, String) Specifies the custom IP address. The default address is pod IP address.

<a name="servicestage_tcp_param"></a>
The `tcp_param` block supports:

* `port` - (Optional, Int) Specifies the port number.

<a name="servicestage_external_access"></a>
The `external_access` block supports:

* `protocol` - (Optional, String) Specifies the protocol. The valid values are **HTTP** and **HTTPS**.

* `address` - (Optional, String) Specifies the access address. For example: `www.example.com`.

* `port` - (Optional, Int) Specifies the listening port of the application component process.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The instance ID in UUID format.

* `status` - The instance status, which supports:
  + **FAILED**
  + **RUNNING**
  + **DOWN**
  + **STOPPED**
  + **UNKNOWN**
  + **PARTIALLY_FAILED**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Instances can be imported using their related `application_id`, `component_id` and `id`, separated by a slash (/), e.g.

```bash
terraform import huaweicloud_servicestage_component_instance.test 4e65a759-e7b1-4e9e-8277-857f8e261f3c/4e65a759-e7b1-4e9e-8277-857f8e261f3c/c0a13d88-d4e3-11ec-93a9-0255ac101d30
```
