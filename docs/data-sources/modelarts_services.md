---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_services"
description: ""
---

# huaweicloud_modelarts_services

Use this data source to get services of ModelArts.

## Example Usage

```hcl
data "huaweicloud_modelarts_services" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `service_id` - (Optional, String) Service ID.

* `name` - (Optional, String) Service name.

* `model_id` - (Optional, String) The model ID which the service used.

* `workspace_id` - (Optional, String) The workspace ID to which a service belongs.  
  The default value is 0, indicating the default workspace.

* `infer_type` - (Optional, String) Inference mode of the service.  
  Value options are as follows:
    + **real-time**: A real-time service. A model is deployed as a web service and provides real-time test UI and monitoring.
    + **batch**: A batch service, which can perform inference on batch data and automatically stops after data is processed.
    + **edge**: An edge service, which uses Intelligent EdgeFabric (IEF) to deploy a model as a web service on an edge
      node created on IEF.

* `status` - (Optional, String) Service status.  
  Value options are as follows:
    + **running**: The service is running properly.
    + **deploying**: The service is being deployed, including image creation and resource scheduling deployment.
    + **concerning**: An alarm has been generated, indicating that some backend instances malfunction.
    + **failed**: Deploying the service failed. For details about the failure cause, see the event and log tab pages.
    + **stopped**: The service has been stopped.
    + **finished**: Service running is completed. This status is available only for batch services.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `services` - The list of services.
  The [services](#ModelartServices_services) structure is documented below.

<a name="ModelartServices_services"></a>
The `services` block supports:

* `id` - Services ID.

* `name` - Services name.

* `workspace_id` - The workspace ID to which a service belongs.  
  The default value is 0, indicating the default workspace.

* `description` - The description of the service.  

* `status` - Service status.  
  Value options are as follows:
    + **running**: The service is running properly.
    + **deploying**: The service is being deployed, including image creation and resource scheduling deployment.
    + **concerning**: An alarm has been generated, indicating that some backend instances malfunction.
    + **failed**: Deploying the service failed. For details about the failure cause, see the event and log tab pages.
    + **stopped**: The service has been stopped.
    + **finished**: Service running is completed. This status is available only for batch services.

* `infer_type` - Inference mode of the service.  
  Value options are as follows:
    + **real-time**: A real-time service. A model is deployed as a web service and provides real-time test UI and monitoring.
    + **batch**: A batch service, which can perform inference on batch data and automatically stops after data is processed.
    + **edge**: An edge service, which uses Intelligent EdgeFabric (IEF) to deploy a model as a web service on an edge
      node created on IEF.

* `is_free` - Whether a free-of-charge flavor is used.

* `schedule` - Service scheduling configuration, which can be configured only for real-time services.  
By default, this parameter is not used. Services run for a long time.
  The [schedule](#ModelartServices_ServicesSchedule) structure is documented below.

* `additional_properties` - Additional service attribute, which facilitates service management.  
  The [additional_properties](#ModelartServices_ServicesAdditionalProperty) structure is documented below.

* `invocation_times` - Number of service calls.

* `failed_times` - Number of failed service calls.

* `is_shared` - Whether a service is subscribed.

* `shared_count` - Number of subscribed services.

* `is_opened_sample_collection` - Whether to enable data collection, which defaults to false.

* `owner` - User to which a service belongs.

<a name="ModelartServices_ServicesSchedule"></a>
The `schedule` block supports:

* `duration` - Value mapping a time unit.  
  For example, if the task stops after two hours, set `time_unit` to HOURS and `duration` to 2.

* `time_unit` - Scheduling time unit. Possible values are **DAYS**, **HOURS**, and **MINUTES**.

* `type` - Scheduling type. Only the value **stop** is supported.

<a name="ModelartServices_ServicesAdditionalProperty"></a>
The `additional_properties` block supports:

* `smn_notification` - SMN message notification configuration.
  The [smn_notification](#ModelartServices_AdditionalPropertySmnNotification) structure is documented below.

* `log_report_channels` - Advanced Log configuration.
  The [log_report_channels](#ModelartServices_AdditionalPropertyLogReportChannel) structure is documented below.

<a name="ModelartServices_AdditionalPropertySmnNotification"></a>
The `smn_notification` block supports:

* `topic_urn` - URN of an SMN topic.

* `events` - Event ID.  
  Value options are as follows::
    + **1**: failed.
    + **2**: stopped.
    + **3**: running.
    + **7**: alarm.
    + **9**: deleted.
    + **11**: pending.

<a name="ModelartServices_AdditionalPropertyLogReportChannel"></a>
The `log_report_channels` block supports:

* `type` - The type of log report channel. The valid value is **LTS**.  
  If this parameter is configured, the advanced log management service, Log Tank Service (LTS) will be used.
  If not, the ModelArts log system will be used, which provides simple log query and caches runtime logs
   for a maximum of seven days.
