---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_asv2_policies"
description: |-
  Use this data source to get the list of the AS group policics and AS bandwidth policies.
---

# huaweicloud_asv2_policies

Use this data source to get the list of the AS group policics and AS bandwidth policies.

## Example Usage

```hcl
data "huaweicloud_asv2_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scaling_policy_id` - (Optional, String) Specifies the AS policy ID.

* `scaling_policy_type` - (Optional, String) Specifies the AS policy type.
  The valid values are as follows:
  + **ALARM**: Alarm policy.
  + **SCHEDULED**: Scheduled policy.
  + **RECURRENCE**: Periodic policy.

* `scaling_policy_name` - (Optional, String) Specifies the AS policy name.
  Fuzzy search is supported.

* `scaling_resource_id` - (Optional, String) Specifies the ID of the resource associate with the AS policy.

* `scaling_resource_type` - (Optional, String) Specifies the resource type associate with the AS policy.
  The valid values are as follows:
  + **SCALING_GROUP**: AS group.
  + **BANDWIDTH**: Bandwidth.

* `alarm_id` - (Optional, String) Specifies the alarm rule ID associate with the AS policy.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  Only support **all_granted_eps**.
  This field is only valid for enterprise users.

* `sort_by` - (Optional, String) Specifies the sorting method of the AS policies.
  The valid values are as follows:
  + **POLICY_NAME**: AS policies are sorted by name.
  + **TRIGGER_CONDITION**: AS policies are sorted by trigger condition.
  + **CREATE_TIME**: AS policies are sorted based on the creation time.

* `order` - (Optional, String) Specifies the sorting order of the AS policies.
  The valid values are as follows:
  + **ASC**: ascending order.
  + **DESC**: descending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scaling_policies` - The list of the AS policies.

  The [scaling_policies](#scaling_policies_struct) structure is documented below.

<a name="scaling_policies_struct"></a>
The `scaling_policies` block supports:

* `scaling_policy_id` - The AS policy ID.

* `scaling_policy_name` - The AS policy name.

* `policy_status` - The AS policy status.

* `description` - The AS policy description.

* `scaling_policy_type` - The AS policy type.

* `scaling_resource_id` - The ID of the resource associate with the AS policy.

* `scaling_resource_type` - The  resource type associate with the AS policy.

* `alarm_id` - The alarm rule ID associate with the AS policy.

* `scaling_policy_action` - The AS policy execute actions.

  The [scaling_policy_action](#scaling_policies_scaling_policy_action_struct) structure is documented below.

* `scheduled_policy` - The schedule and periodic policy contents.

  The [scheduled_policy](#scaling_policies_scheduled_policy_struct) structure is documented below.

* `meta_data` - The bandwidth policy additional information.

  The [meta_data](#scaling_policies_meta_data_struct) structure is documented below.

* `cool_down_time` - The cooldown period, in seconds.

* `create_time` - The creation time of the AS policy. in UTC format.

<a name="scaling_policies_scaling_policy_action_struct"></a>
The `scaling_policy_action` block supports:

* `operation` - The operation to be performed.

* `size` - The operation size.

* `percentage` - The percentage of instances to be operated.

* `limits` - The operation restrictions.

<a name="scaling_policies_scheduled_policy_struct"></a>
The `scheduled_policy` block supports:

* `launch_time` - The time when the scaling action is triggered.

* `recurrence_type` - The periodic triggering type.

* `recurrence_value` - The day when a periodic scaling action is triggered.

* `start_time` - The start time of the scaling action triggered periodically.

* `end_time` - The end time of the scaling action triggered periodically.

<a name="scaling_policies_meta_data_struct"></a>
The `meta_data` block supports:

* `metadata_bandwidth_share_type` - The bandwidth sharing type in the bandwidth scaling policy.

* `metadata_eip_id` - The EIP ID for the bandwidth in the bandwidth scaling policy.

* `metadata_eip_address` - The EIP IP address for the bandwidth in the bandwidth scaling policy.
