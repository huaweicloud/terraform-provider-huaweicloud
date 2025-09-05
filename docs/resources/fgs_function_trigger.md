---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_trigger"
description: ""
---

# huaweicloud_fgs_function_trigger

Manages the function trigger resource within HuaweiCloud.

## Example Usage

### Create the Timing Triggers with rate and cron schedule types

```hcl
variable "function_urn" {}
variable "trigger_name" {}

// Timing trigger (with rate schedule type)
resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = var.function_urn
  type         = "TIMER"
  event_data   = jsonencode({
    "name": format("%s_rate", var.trigger_name),
    "schedule_type": "Rate",
    "user_event": "Created by terraform script",
    "schedule": "3m"
  })
}

// Timing trigger (with cron schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_cron" {
  function_urn = var.function_urn
  type         = "TIMER"
  event_data   = jsonencode({
    "name": format("%s_cron", var.trigger_name),
    "schedule_type": "Cron",
    "user_event": "Created by terraform script",
    "schedule": "@every 1h30m"
  })
}
```

### Create an OBS Event Trigger with EventGrid

```hcl
variable "function_urn" {}
variable "trigger_name_suffix" {}
variable "agency_name" {} # Allow FunctionGraph to access EG
variable "bucket_name" {}

resource "huaweicloud_fgs_function_trigger" "obs_trigger" {
  function_urn                   = var.function_urn
  type                           = "EVENTGRID"
  cascade_delete_eg_subscription = true
  status                         = "ACTIVE"
  event_data                     = jsonencode({
    "channel_id"   = "your_event_channel_id"
    "channel_name" = "default"
    "source_name"  = "HC.OBS.DWR"
    "trigger_name" = var.trigger_name_suffix
    "agency"       = var.agency_name
    "bucket"       = var.bucket_name
    "event_types"  = ["OBS:DWR:ObjectCreated:PUT", "OBS:DWR:ObjectCreated:POST"]
    "Key_encode"   = true
  })

  lifecycle {
    ignore_changes = [
      event_data, # Parts of event_data content are not returned by FunctionGraph Service and will cause changes
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the function trigger is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `function_urn` - (Required, String, ForceNew) Specifies the function URN to which the function trigger belongs.  
  Changing this will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the function trigger.  
  The valid values are **TIMER**, **APIG**, **CTS**, **DDS**, **DEDICATEDGATEWAY**, etc.
  Changing this will create a new resource.

  -> For more available values, please refer to the [documentation table 3](https://support.huaweicloud.com/intl/en-us/api-functiongraph/functiongraph_06_0122.html#section2).

* `event_data` - (Required, String) Specifies the detailed configuration of the function trigger event, in JSON
  format.  
  For various types of trigger parameter configurations, please refer to the
  [documentation](https://support.huaweicloud.com/intl/en-us/api-functiongraph/functiongraph_06_0122.html#functiongraph_06_0122__request_TriggerEventDataRequestBody).

  -> Please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-functiongraph/functiongraph_06_0124.html#functiongraph_06_0124__request_UpdateriggerEventData)
     for updateable fields.

* `status` - (Optional, String) Specifies the status of the function trigger.  
  The valid values are **ACTIVE** and **DISABLED**.  
  About `DDS` and `Kafka` triggers, the default value is **DISABLED**, for the other triggers, the default value is
  **ACTIVE**.

  -> Currently, only some triggers support setting the **DISABLED** value, such as `TIMER`, `DDS`, `DMS`, `KAFKA` and
     `LTS`. For more details, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-functiongraph/functiongraph_06_0122.html).

* `cascade_delete_eg_subscription` - (Optional, Bool) Specifies whether to cascade delete the related EG event
  subscription of the function trigger.  
  This parameter is mainly used for `EVENTGRID` type triggers to automatically clean up related EG event subscription
  when the trigger is deleted.  
  Defaults to **false**.

  -> This parameter cannot be imported. The effective value when deleted is the value of the parameter after the last
     change.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - resource ID in UUID format.

* `created_at` - The creation time of the function trigger.

* `updated_at` - The latest update time of the function trigger.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

Function trigger can be imported using the `function_urn`, `type` and `id`, separated by the slashes (/), e.g.

```bash
$ terraform import huaweicloud_fgs_function_trigger.test <function_urn>/<type>/<id>
```
