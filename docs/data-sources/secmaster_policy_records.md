---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_policy_records"
description: |-
  Use this data source to get the list of policy records within HuaweiCloud SecMaster.
---

# huaweicloud_secmaster_policy_records

Use this data source to get the list of policy records within HuaweiCloud SecMaster.

## Example Usage

```hcl
variable "workspace_id" {}
variable "policy_id" {}

data "huaweicloud_secmaster_policy_records" "test" {
  workspace_id = var.workspace_id
  policy_id    = var.policy_id

  sort {
    sort_by = "create_time"
    order   = "desc"
  }

  group_by {
    group_by_fields = ["workspace_id"]
    group_by_hit {
      source = "defense_policy_object"
      dest   = "defense_policy_list"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the policy belongs.

* `policy_id` - (Required, String) Specifies the ID of the policy to query records for.

* `condition` - (Optional, List) Specifies the query conditions.

  The [condition](#policys_condition_struct) structure is documented below.

* `sort` - (Optional, List) Specifies the sort conditions.

  The [sort](#policys_sort_struct) structure is documented below.

* `group_by` - (Optional, List) Specifies the aggregation conditions.

  The [group_by](#policys_group_by_struct) structure is documented below.

<a name="policys_condition_struct"></a>
The `condition` block supports:

* `conditions` - (Optional, List) Specifies the list of query conditions.

  The [conditions](#policys_conditions_struct) structure is documented below.

* `logics` - (Optional, List) Specifies the list of condition names.

<a name="policys_conditions_struct"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the condition name.

* `data` - (Optional, List) Specifies the condition values.

<a name="policys_sort_struct"></a>
The `sort` block supports:

* `sort_by` - (Optional, String) Specifies the sort field.

* `order` - (Optional, String) Specifies the sort direction.

<a name="policys_group_by_struct"></a>
The `group_by` block supports:

* `group_by_fields` - (Optional, List) Specifies the aggregation fields.

* `group_by_hit` - (Optional, List) Specifies the aggregation result mapping.

  The [group_by_hit](#policys_group_by_hit_struct) structure is documented below.

<a name="policys_group_by_hit_struct"></a>
The `group_by_hit` block supports:

* `source` - (Optional, String) Specifies the source field.

* `dest` - (Optional, String) Specifies the destination field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of policy records.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - The ID of the policy record.

* `policy_id` - The ID of the policy.

* `policy_task_id` - The ID of the policy task.

* `policy_task_name` - The name of the policy task.

* `policy_category` - The category of the policy. Valid values: **BLOCK**, **ALLOW**.

* `policy_direction` - The direction of the policy. Valid values: **INGRESS**, **EGRESS**.

* `policy_automation` - The automation type of the policy. Valid values: **AUTOMATION**, **MANUAL**.

* `block_target` - The block target of the policy.

* `description` - The description of the policy record.

* `is_deleted` - Whether the policy record is deleted.

* `trigger_flag` - Whether the policy is triggered.

* `create_time` - The creation time of the policy record.

* `update_time` - The update time of the policy record.

* `dataclass_id` - The data class ID of the policy record.

* `domain_id` - The domain ID.

* `domain_name` - The domain name.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `region_name` - The region name.

* `workflow_instance` - The workflow instance information.
  The [workflow_instance](#workflow_instance_struct) structure is documented below.

* `block_age` - The block aging information.
  The [block_age](#block_age_struct) structure is documented below.

* `policyrecord_type` - The policy record type information.
  The [policyrecord_type](#policyrecord_type_struct) structure is documented below.

* `environment` - The environment information.
  The [environment](#environment_struct) structure is documented below.

* `defense_policy_list` - The defense policy list.
  The [defense_policy_list](#defense_policy_list_struct) structure is documented below.

<a name="workflow_instance_struct"></a>
The `workflow_instance` block supports:

* `workflow_instance_id` - The ID of the workflow instance.

* `workflow_instance_name` - The name of the workflow instance.

* `workflow_instance_status` - The status of the workflow instance.

<a name="block_age_struct"></a>
The `block_age` block supports:

* `is_block_ageing` - Whether the block is aging.

* `block_ageing` - The block aging time.

<a name="policyrecord_type_struct"></a>
The `policyrecord_type` block supports:

* `policy_type` - The type of the policy record.

* `id` - The ID of the policy record type.

* `category` - The category of the policy record type.

<a name="environment_struct"></a>
The `environment` block supports:

* `domain_id` - The domain ID.

* `domain_name` - The domain name.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `region_name` - The region name.

* `vendor_type` - The vendor type.

<a name="defense_policy_list_struct"></a>
The `defense_policy_list` block supports:

* `defense_id` - The ID of the defense.

* `defense_type` - The type of the defense.

* `defense_policy_id` - The ID of the defense policy.

* `defense_policy_name` - The name of the defense policy.

* `defense_connection_id` - The ID of the defense connection.

* `defense_connection_name` - The name of the defense connection.

* `defense_connection_region_id` - The region ID of the defense connection.

* `defense_connection_region_name` - The region name of the defense connection.

* `defense_block_status` - The block status of the defense.

* `defense_block_action` - The block action of the defense.

* `defense_failure_description` - The failure description of the defense.

* `defense_creator_id` - The creator ID of the defense.

* `defense_creator_name` - The creator name of the defense.

* `target_project_id` - The target project ID.

* `target_project_name` - The target project name.

* `target_enterprise_id` - The target enterprise ID.

* `target_enterprise_name` - The target enterprise name.

* `description` - The description of the defense policy.
