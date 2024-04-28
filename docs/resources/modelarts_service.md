---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_service"
description: ""
---

# huaweicloud_modelarts_service

Manages a ModelArts service resource within HuaweiCloud.  

## Example Usage

### Create a real-time service

```hcl
variable "model_id" {}

resource "huaweicloud_modelarts_service" "test" {
  name        = "demo"
  infer_type  = "real-time"
  description = "This is a demo"

  config {
    specification  = "modelarts.vm.gpu.p4u8.container"
    instance_count = 1
    weight         = 100
    model_id       = var.model_id
    envs = {
      "a" : "1",
      "b" : "2"
    }
  }

  additional_properties {
    smn_notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = [3]
    }
    log_report_channels {
      type = "LTS"
    }
  }
}
```

### Create a real-time service and configuring it to automatically stop

```hcl
variable "model_id" {}

resource "huaweicloud_modelarts_service" "test" {
  name        = "demo"
  infer_type  = "real-time"
  description = "This is a demo"

  config {
    specification  = "modelarts.vm.gpu.p4u8.container"
    instance_count = 1
    weight         = 100
    model_id       = var.model_id
    envs = {
      "a" : "1",
      "b" : "2"
    }
  }

  additional_properties {
    smn_notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = [3]
    }
    log_report_channels {
      type = "LTS"
    }
  }

  schedule {
    type      = "stop"
    duration  = 1
    time_unit = "HOURS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Service name, which consists of 1 to 64 characters.  
  Only letters, digits, hyphens (-), and underscores (_) are allowed.

  Changing this parameter will create a new resource.

* `infer_type` - (Required, String, ForceNew) Inference mode.  
  Value options are as follows:
    + **real-time**: A real-time service. A model is deployed as a web service and provides real-time test UI and monitoring.
    + **batch**: A batch service, which can perform inference on batch data and automatically stops after data is processed.
    + **edge**: An edge service, which uses Intelligent EdgeFabric (IEF) to deploy a model as a web service on an edge
                node created on IEF.

  Changing this parameter will create a new resource.

* `config` - (Required, List) Model running configurations.  
  If `infer_type` is **batch** or **edge**, you can configure only one model.
  If `infer_type` is **real-time**, you can configure multiple models and assign weights based on service requirements.
  However, the versions of multiple models must be unique.
  The [Config](#ModelartsService_Config) structure is documented below.

* `workspace_id` - (Optional, String, ForceNew) ID of the workspace to which a service belongs.  
  The default value is **0**, indicating the default workspace.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the service.  

* `pool_name` - (Optional, String, ForceNew) The ID of the new dedicated resource pool.  
  When using dedicated resource pool to deploy services, ensure that the cluster status is normal.
  If both `pool_name` and `config.pool_name` are configured, `pool_name` in real-time config is preferred.

  Changing this parameter will create a new resource.

* `vpc_id` - (Optional, String, ForceNew) The VPC ID to which a real-time service instance is deployed.  
  By default, this parameter is left blank. In this case, ModelArts allocates a dedicated VPC to each user,
  and users are isolated from each other.
  To access other service components in the VPC of the service instance,
  set this parameter to the ID of the corresponding VPC.
  Once a VPC is configured, it cannot be modified. If both `vpc_id` and `pool_name` are configured,
  only the dedicated resource pool takes effect.

  Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) The subnet ID.  
  By default, this parameter is left blank.
  This parameter is mandatory if `vpc_id` is configured.
  Enter the network ID displayed in the subnet details on the VPC management console.
  A subnet provides dedicated network resources that are isolated from other networks.

  Changing this parameter will create a new resource.

* `security_group_id` - (Optional, String, ForceNew) The security group ID.  
  By default, this parameter is left blank.
  This parameter is mandatory if `vpc_id` is configured.
  A security group is a virtual firewall that provides secure network access control policies for service instances.
  A security group must contain at least one inbound rule to permit the requests whose protocol is TCP,
  source address is 0.0.0.0/0, and port number is 8080.

  Changing this parameter will create a new resource.

* `schedule` - (Optional, List) Service scheduling configuration, which can be configured only for real-time services.
  If this parameter is configured, the service will be stopped automatically.
  By default, the service runs for a long time.
  The [Schedule](#ModelartsService_Schedule) structure is documented below.

* `additional_properties` - (Optional, List) Additional properties.
  The [AdditionalProperty](#ModelartsService_AdditionalProperty) structure is documented below.

* `change_status_to` - (Optional, String) Which status you want to change the service to.
  The valid value can be **running** or **stopped**.  
  If this parameter is not configured, the service status is not changed.

<a name="ModelartsService_Config"></a>
The `Config` block supports:

* `custom_spec` - (Optional, List) Custom resource specifications.  
  The [CustomSpec](#ModelartsService_CustomSpec) structure is documented below.

* `envs` - (Optional, Map) Environment variable key-value pair required for running a model.  

* `specification` - (Optional, String) Resource flavors.  
  The valid values are **modelarts.vm.cpu.2u**, **modelarts.vm.gpu.p4** (must be requested for),
  **modelsarts.vm.ai1.a310** (must be requested for),
  and **custom** (available only when the service is deployed in a dedicated resource pool) in the current version.
  To request for a flavor, submit a service ticket and obtain permissions from ModelArts O&M engineers.
  If this parameter is set to custom, the custom_spec parameter must be specified.
  Value options are as follows:
    + **modelarts.vm.cpu.free**: [Time-limited free] CPU: 1 vCPUs | 4 GiB.
    + **modelarts.vm.cpu.2u**: CPU: 2 vCPUs | 8 GiB.
    + **modelarts.vm.gpu.p4**: CPU: 1 vCPUs | 4 GiB GPU：P4 (must be requested for).
    + **modelarts.vm.gpu.p4u8.container**: CPU: 8 vCPUs | 32 GiB GPU：P4.
    + **modelarts.vm.gpu.t4u8.container**: CPU: 8 vCPUs | 32 GiB GPU：T4.
    + **custom**: available only when the service is deployed in a dedicated resource pool,
        and the `custom_spec` parameter must be specified.

* `weight` - (Optional, Int) Weight of traffic allocated to a model.  
  This parameter is mandatory only when `infer_type` is set to **real-time**.
  The sum of all weights must be equal to 100. If multiple model versions are configured with different
  traffic weights in a real-time service, ModelArts will continuously access the prediction API of the
  service and forward prediction requests to the model instances of the corresponding versions based on the weights.

* `model_id` - (Optional, String) Model ID, which can be obtained by calling the API for obtaining a model list.  

* `src_path` - (Optional, String) OBS path to the input data of a batch job.  
  Mandatory for batch services.

* `req_uri` - (Optional, String) Inference API called in a batch task, which is the RESTful API exposed in the model image.
  Mandatory for batch services.
  You must select an API URL from the config.json file of the model for inference.
  If a built-in inference image of ModelArts is used, the API is displayed as /.

* `mapping_type` - (Optional, String) Mapping type of the input data. Mandatory for batch services.  
  The value can be file or csv. file indicates that each inference request corresponds to a file
  in the input data directory.
  If this parameter is set to file, req_uri of the model can have only one input parameter and the type
   of this parameter is file.
  If this parameter is set to csv, each inference request corresponds to a row of data in the CSV file.
  When csv is used, the file in the input data directory can only be suffixed with .csv,
  and the mapping_rule parameter must be configured to map the index of each parameter in
  the inference request body to the CSV file.

* `pool_name` - (Optional, String) The ID of the new dedicated resource pool.  
  When using dedicated resource pool to deploy services, ensure that the cluster status is normal.
  If both `pool_name` and `config.pool_name` are configured, `pool_name` in real-time config is preferred.

* `nodes` - (Optional, List) Edge node ID array. Mandatory for edge services.  
  The node ID is the edge node ID on IEF, which can be obtained after the edge node is created on IEF.

* `mapping_rule` - (Optional, Map) Mapping between input parameters and CSV data. Optional for batch services.  
  This parameter is mandatory only when mapping_type is set to csv.
  The mapping rule is similar to the definition of the input parameters in the config.json file.
  You only need to configure the index parameters under each parameter of the string, number, integer,
  or boolean type, and specify the value of this parameter to the values of the index parameters
  in the CSV file to send an inference request. Use commas (,) to separate multiple pieces of CSV data.
  The values of the index parameters start from 0. If the value of the index parameter is -1, ignore this parameter.
  For details, see the sample of creating a batch service.

* `src_type` - (Optional, String) Data source type, which can be ManifestFile. Mandatory for batch services.  
  By default, this parameter is left blank, indicating that only files in the src_path directory are read.
  If this parameter is set to ManifestFile, src_path must be set to a specific manifest path.
  Multiple data paths can be specified in the manifest file. For details, see the manifest inference specifications.

* `dest_path` - (Optional, String) OBS path to the output data of a batch job. Mandatory for batch services.  

* `instance_count` - (Optional, Int) Number of instances deployed for a model.  
  The maximum number of instances is 5. To use more instances, submit a service ticket.

* `additional_properties` - (Optional, Map) Additional attributes for model deployment, facilitating service instance management.

<a name="ModelartsService_CustomSpec"></a>
The `CustomSpec` block supports:

* `memory` - (Optional, Int) Memory in MB, which must be an integer.

* `cpu` - (Optional, Float) Number of CPU cores, which can be a decimal. The value cannot be smaller than 0.01.

* `gpu_p4` - (Optional, Float) Number of GPU cores, which can be a decimal.  
  The value cannot be smaller than 0, which allows up to two decimal places.

* `ascend_a310` - (Optional, Int) Number of Ascend chips. Either this parameter or `gpu_p4` is configured.

<a name="ModelartsService_Schedule"></a>
The `Schedule` block supports:

* `duration` - (Required, Int) Value mapping a time unit.  
  For example, if the task stops after two hours, set time_unit to HOURS and duration to 2.

* `time_unit` - (Required, String) Scheduling time unit. Possible values are DAYS, HOURS, and MINUTES.

* `type` - (Required, String) Scheduling type. Only the value **stop** is supported.

<a name="ModelartsService_AdditionalProperty"></a>
The `AdditionalProperty` block supports:

* `smn_notification` - (Optional, List) SMN message notification configuration.
  The [SmnNotification](#ModelartsService_SmnNotification) structure is documented below.

* `log_report_channels` - (Optional, List) Advanced Log configuration.
  The [LogReportChannel](#ModelartsService_LogReportChannel) structure is documented below.

<a name="ModelartsService_SmnNotification"></a>
The `SmnNotification` block supports:

* `topic_urn` - (Optional, String) URN of an SMN topic.

* `events` - (Optional, List) Event ID.  
  Value options are as follows:
    + **1**: failed.
    + **2**: stopped.
    + **3**: running.
    + **7**: alarm.
    + **9**: deleted.
    + **11**: pending.

<a name="ModelartsService_LogReportChannel"></a>
The `LogReportChannel` block supports:

* `type` - (Optional, String) The type of log report channel. The valid value is **LTS**.  
  If this parameter is configured, the advanced log management service, Log Tank Service (LTS) will be used.
  If not, the ModelArts log system will be used, which provides simple log query and caches runtime logs
   for a maximum of seven days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `owner` - User to which a service belongs

* `status` - Service status.  
  Value options are as follows:
    + **running**: The service is running properly.
    + **deploying**: The service is being deployed, including image creation and resource scheduling deployment.
    + **concerning**: An alarm has been generated, indicating that some backend instances malfunction.
    + **failed**: Deploying the service failed. For details about the failure cause, see the event and log tab pages.
    + **stopped**: The service has been stopped.
    + **finished**: Service running is completed. This status is available only for batch services.

* `access_address` - Access address of an inference request.  
  This parameter is available when `infer_type` is set to **real-time**.

* `bind_access_address` - Request address of a custom domain name.  
  This parameter is available after a domain name is bound.  

* `invocation_times` - Total number of service calls.  

* `failed_times` - Number of failed service calls.

* `is_shared` - Whether a service is subscribed.

* `shared_count` - Number of subscribed services.

* `debug_url` - Online debugging address of a real-time service.  
  This parameter is available only when the model supports online debugging and there is only one instance.

* `is_free` - Whether a free-of-charge flavor is used.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The modelarts service can be imported using `id` e.g.

```bash
$ terraform import huaweicloud_modelarts_service.test 60495dd7-d56b-43c7-8f98-03833833f8e0
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `change_status_to`.
It is generally recommended running `terraform plan` after
importing a dataset. You can then decide if changes should be applied to the dataset, or the resource definition
should be updated to align with the dataset. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_service" "test" {
    ...

  lifecycle {
    ignore_changes = [
      change_status_to,
    ]
  }
}
```
