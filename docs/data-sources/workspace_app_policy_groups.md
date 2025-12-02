---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_policy_groups"
description: |-
  Use this data source to get policy group list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_policy_groups

Use this data source to get policy group list of the Workspace APP within HuaweiCloud.

## Example Usage

### Query all policy groups

```hcl
data "huaweicloud_workspace_app_policy_groups" "test" {}
```

### Query policy groups by type

```hcl
data "huaweicloud_workspace_app_policy_groups" "test" {
  policy_group_type = 4
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the policy groups are located.  
  If omitted, the provider-level region will be used.

* `policy_group_name` - (Optional, String) Specifies the name of the policy group.  
  Fuzzy search is supported.

* `policy_group_type` - (Optional, Int) Specifies the type of the policy group.  
  The valid values are as follows:
  + **0**: VM class.
  + **4**: Custom policy template.

  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policy_groups` - The list of policy groups that match the filter parameters.  
  The [policy_groups](#app_policy_groups) structure is documented below.

<a name="app_policy_groups"></a>
The `policy_groups` block supports:

* `id` - The ID of the policy group.

* `name` - The name of the policy group.

* `priority` - The priority of the policy group.

* `description` - The description of the policy group.

* `targets` - The list of target objects.  
  The [targets](#app_policy_groups_targets) structure is documented below.

* `policies` - The policies of the policy group, in JSON format.

* `created_at` - The creation time of the policy group, in RFC3339 format.

* `updated_at` - The latest update time of the policy group, in RFC3339 format.

<a name="app_policy_groups_targets"></a>
The `targets` block supports:

* `id` - The ID of the target object.

* `name` - The name of the target object.

* `type` - The type of the target object.  
  + **USER**
  + **USERGROUP**
  + **APPGROUP**
  + **OU**
  + **ALL**
