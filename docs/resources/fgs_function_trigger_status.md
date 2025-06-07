---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_trigger_status_action"
description: |-
  Use this resource to update the status of the FunctionGraph trigger within HuaweiCloud.
---

# huaweicloud_fgs_function_trigger_status_action

Use this resource to update the status of the FunctionGraph trigger within HuaweiCloud.

-> 1. One function trigger can only manage one status action resource.
   <br>2. This resource is only a one-time action resource for updating function trigger status. Deleting this resource
   will not revert the trigger status, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "function_urn" {}
variable "function_trigger_id" {}
variable "event_data_in_trigger_definition" {}

resource "huaweicloud_fgs_function_trigger_status_action" "test" {
  function_urn      = var.function_urn
  trigger_type_code = "KAFKA"
  trigger_id        = var.function_trigger_id
  trigger_status    = "ACTIVE"
  event_data        = var.event_data_in_trigger_definition
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the function trigger is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `function_urn` - (Required, String, NonUpdatable) Specifies the URN of the function.

* `trigger_type_code` - (Required, String, NonUpdatable) Specifies the trigger type code.  
  The valid values are as follows:
  + **TIMER**
  + **CTS**
  + **DDS**
  + **DMS**
  + **DIS**
  + **LTS**
  + **OBS**
  + **SMN**
  + **KAFKA**
  + **RABBITMQ**
  + **DEDICATEDGATEWAY**
  + **OPENSOURCEKAFKA**
  + **APIC**
  + **GeminiDB Mongo**
  + **EVENTGRID**
  + **IOTDA**

* `trigger_id` - (Required, String, NonUpdatable) Specifies the trigger ID.  
  Changing this will create a new resource.

* `trigger_status` - (Required, String, NonUpdatable) Specifies the status of the trigger to be changed.  
  Valid values are **ACTIVE** and **DISABLED**.  
  Changing this will create a new resource.

* `event_data` - (Optional, String, NonUpdatable) Specifies the trigger event data configuration, in JSON format.

~> It is not recommended to reference the event_data returned by the function trigger resource.
   <br>Instead, it is recommended to reference the same variable source.
   <br>If not, the terraform will report this error: The provider produced inconsistent final plan.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also function trigger ID).

## Import

This resource cannot be imported.
