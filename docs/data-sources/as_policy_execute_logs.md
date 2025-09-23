---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_policy_execute_logs"
description: |
  Use this data source to get a list of AS policy execution logs within HuaweiCloud.
---

# huaweicloud_as_policy_execute_logs

Use this data source to get a list of AS policy execution logs within HuaweiCloud.

-> Currently, only the latest up to `20` logs can be queried.

## Example Usage

```hcl
variable "scaling_policy_id" {}

data "huaweicloud_as_policy_execute_logs" "test" {
  scaling_policy_id = var.scaling_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the policy execution logs.
  If omitted, the provider-level region will be used.

* `scaling_policy_id` - (Required, String) Specifies the scaling policy ID.

* `log_id` - (Optional, String) Specifies the policy execution log ID.

* `scaling_resource_id` - (Optional, String) Specifies the scaling resource ID.

* `scaling_resource_type` - (Optional, String) Specifies the scaling resource type.
  The value can be **SCALING_GROUP** or **BANDWIDTH**.

* `execute_type` - (Optional, String) Specifies the policy execution type.  
  The valid values are as follows:
  + **SCHEDULED**: automatically triggered scheduled policy.
  + **RECURRENCE**: automatically triggered recurrence policy.
  + **ALARM**: automatically triggered alarm policy.
  + **MANUAL**: manually triggered policy.

* `start_time` - (Optional, String) Specifies the start time of the policy execution used for query. The time format is
  **yyyy-MM-ddThh:mm:ssZ**.  
  The query result is all data with a policy execution time greater than or equal to this value.

* `end_time` - (Optional, String) Specifies the end time of the policy execution used for query. The time format is
  **yyyy-MM-ddThh:mm:ssZ**.  
  The query result shows all data with policy execution time less than this value.

* `status` - (Optional, String) Specifies the policy execution status. The value can be **SUCCESS**, **FAIL**
  or **EXECUTING**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `execute_logs` - All policy execution logs that match the filter parameters.  
  The [execute_logs](#as_execute_logs_attr) structure is documented below.

<a name="as_execute_logs_attr"></a>
The `execute_logs` block supports:

* `id` - The policy execution log ID.

* `status` - The policy execution status.

* `failed_reason` - The reason of policy execution failure.

* `execute_type` - The policy execution type.

* `execute_time` - The policy execution time, the time format is **yyyy-MM-ddThh:mm:ssZ**.

* `scaling_policy_id` - The scaling policy ID.

* `scaling_resource_id` - The scaling resource ID.

* `scaling_resource_type` - The scaling resource type.

* `type` - The policy execution task type. The value can be **REMOVE**, **ADD** or **SET**.

* `old_value` - The scaling original value.
  + When `scaling_resource_type` is **SCALING_GROUP**, this field represents the number of instances.
  + When `scaling_resource_type` is **BANDWIDTH**, this field represents the bandwidth size, in Mbit/s.

* `desire_value` - The scaling target value.
  + When `scaling_resource_type` is **SCALING_GROUP**, this field represents the number of instances.
  + When `scaling_resource_type` is **BANDWIDTH**, this field represents the bandwidth size, in Mbit/s.

* `limit_value` - The operational limitations. When `scaling_resource_type` is **BANDWIDTH** and `type` is not **SET**,
  this field takes effect in Mbit/s.
  + When `type` is **ADD**, this field represents the maximum bandwidth limit.
  + When `type` is **REMOVE**, this field represents the minimum bandwidth limit.

* `job_records` - The concrete tasks included in executing actions.  
  The [job_records](#as_job_records_attr) structure is documented below.

* `metadata` - The additional information.

<a name="as_job_records_attr"></a>
The `job_records` block supports:

* `job_name` - The job name.

* `job_status` - The job execution status. The value can be **SUCCESS** or **FAIL**.

* `record_type` - The record type. The value can be **API** or **MEG**.

* `record_time` - The record time, the time format is **YYYY-MM-DDThh:mmZ**.

* `request` - The request information, the field is valid while `record_type` is **API**.

* `response` - The response information, the field is valid while `record_type` is **API**.

* `code` - The response code, the field is valid while `record_type` is **API**.

* `message` - The message content, the field is valid while `record_type` is **MEG**.
