---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_policy_groups"
description: |-
  Use this data source to get the list of Workspace policy groups within HuaweiCloud.
---

# huaweicloud_workspace_policy_groups

Use this data source to get the list of Workspace policy groups within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_policy_groups" "test" {}
```

### Filter policy groups by priority

```hcl
variable "policy_priority" {}

data "huaweicloud_workspace_policy_groups" "test" {
  priority = var.policy_priority
}
```

### Filter policy groups by name

```hcl
variable "policy_group_name" {}

data "huaweicloud_workspace_policy_groups" "test" {
  policy_group_name = var.policy_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.  
  If omitted, the provider-level region will be used.

* `policy_group_id` - (Optional, String) Specifies the ID of the policy group.

* `policy_group_name` - (Optional, String) Specifies the name of the policy group.  
  The name support fuzzy match.

* `priority` - (Optional, Int) Specifies the priority of the policy group.  
  Defaults to **0**.

* `description` - (Optional, String) Specifies the description of the policy group.  
  The description support fuzzy match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policy_groups` - The list of policy groups that match the filter parameters.  
  The [policy_groups](#workspace_policy_groups_attr) structure is documented below.

<a name="workspace_policy_groups_attr"></a>
The `policy_groups` block supports:

* `policy_group_id` - The ID of the policy group.

* `policy_group_name` - The name of the policy group.

* `priority` - The priority of the policy group.

* `update_time` - The update time of the policy group, in RFC3339 format.

* `description` - The description of the policy group.

* `policies` - The list of policy configurations.  
  The [policies](#workspace_policy_groups_policies) structure is documented below.

* `targets` - The list of target configurations.  
  The [targets](#workspace_policy_groups_targets) structure is documented below.

<a name="workspace_policy_groups_policies"></a>
The `policies` block supports:

* `peripherals` - The peripheral device policies, in JSON format.

* `audio` - The audio policies, in JSON format.

* `client` - The client policies, in JSON format.

* `display` - The display policies, in JSON format.

* `file_and_clipboard` - The file and clipboard policies, in JSON format.

* `session` - The session policies, in JSON format.

* `virtual_channel` - The virtual channel policies, in JSON format.

* `watermark` - The watermark policies, in JSON format.

* `keyboard_mouse` - The keyboard and mouse policies, in JSON format.

* `seamless` - The general audio and video bypass policies, in JSON format.

* `personalized_data_mgmt` - The personalized data management policies, in JSON format.

* `custom` - The custom policies, in JSON format.

* `record_audit` - The screen recording audit policies, in JSON format.

<a name="workspace_policy_groups_targets"></a>
The `targets` block supports:

* `target_id` - The ID of the target.

* `target_type` - The type of the target.

* `target_name` - The name of the target.
