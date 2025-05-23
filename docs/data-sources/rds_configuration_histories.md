---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_configuration_histories"
description: |-
  Use this data source to query the parameter change history of a RDS instance.
---

# huaweicloud_rds_configuration_histories

Use this data source to query the parameter change history of a RDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_configuration_histories" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `start_time` - (Optional, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format. 
  The default value is seven days before the current time.

* `end_time` - (Optional, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format. 
  The default value is the current time.

* `param_name` - (Optional, String) Specifies the parameter name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `histories` - Indicates the list of parameter change history objects.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `parameter_name` - Indicates the name of the parameter.

* `old_value` - Indicates the parameter value before the change.

* `new_value` - Indicates the parameter value after the change.

* `update_result` - Indicates the result of the change operation.
  Value can be as follow:
  + **SUCCEDD**: Indicates the parameter change succeeded.
  + **FAILED**: Indicates the parameter change failed.

* `applied` - Indicates whether the new value will be applied to the instance.

* `update_time` - Indicates the start time in the **yyyy-mm-ddThh:mm:ssZ** format. The default value 
  is seven days before the current time. T is the separator between the calendar and the hourly 
  notation of time. Z indicates the time zone offset. For example, in the Beijing time zone, the 
  time zone offset is shown as +0800.

* `apply_time` - Indicates the end time in the **yyyy-mm-ddThh:mm:ssZ** format. The default value 
  is the current time. T is the separator between the calendar and the hourly notation of time. 
  Z indicates the time zone offset. For example, in the Beijing time zone, the time zone offset is 
  shown as +0800.
