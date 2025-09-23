---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_operation_logs"
description: |-
  Use this data source to get a list of user operation logs.
---

# huaweicloud_dbss_operation_logs

Use this data source to get a list of user operation logs.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dbss_operation_logs" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the audit instance ID to which the user operation logs belong.

* `user_name` - (Optional, String) Specifies the name of the operation user.

* `operate_name` - (Optional, String) Specifies the name of the operation object.

* `result` - (Optional, String) Specifies the execution result of user operation.
  The value can be **success** or **fail**.

* `start_time` - (Optional, String) Specifies the start time of the user operation.
  The time format is UTC. e.g. **2024-09-01 09:00:10**.

* `end_time` - (Optional, String) Specifies the end time of the user operation.
  The time format is UTC. e.g. **2024-09-01 09:15:20**.

* `time_range` - (Optional, String) Specifies the time segment.
  The valid values are as follows:
  + **HALF_HOUR**
  + **HOUR**
  + **THREE_HOUR**
  + **TWELVE_HOUR**
  + **DAY**
  + **WEEK**
  + **MONTH**

-> 1.The parameter `start_time` and `end_time` must be used together.
   <br>2. If parameter `time_range`, `start_time` and `end_time` are set at the same time,
   only the parameter `time_range` will take effect.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of the user operation logs.

  The [logs](#logs_struct) structure is documented below.

<a name="logs_struct"></a>
The `logs` block supports:

* `id` - The ID of the user operation log.

* `name` - The name of the operation object.

* `description` - The description of the user operation.

* `result` - The execution result of user operation.

* `action` - The type of the user operation.
  The valid values are as follows:
  + **create**
  + **update**
  + **delete**
  + **operate**

* `function` - The function type of the operation record.

* `user` - The name of the operation user.

* `time` - The time of the operation record is generated, in UTC format.
