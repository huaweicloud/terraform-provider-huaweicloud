---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_policy_group"
description: |-
  Manages a Workspace APP policy group resource within HuaweiCloud.
---

# huaweicloud_workspace_app_policy_group

Manages a Workspace APP policy group resource within HuaweiCloud.

## Example Usage

```hcl
variable "app_policy_group_name" {}
variable "app_group_id" {}
variable "app_group_name" {}

resource "huaweicloud_workspace_app_policy_group" "test" {
  name     = var.app_policy_group_name
  priority = 1

  targets {
    id   = var.app_group_id
    name = var.app_group_name
    type = "APPGROUP"
  }

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval" : 5,
      "session_persistence_time" : 180,
      "forbid_screen_capture" : false
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the policy group.  
  The name valid length is limited from `1` to `55`, only English characters, digits, underscores (_) are allowed.  
  The name must be unique.

* `priority` - (Optional, Int) Specifies the priority of the policy group.  
  The valid value is range from `1` to the total number of policy groups.  
  The smaller value means the higher priority.

* `description` - (Optional, String) Specifies the description of the policy group.  
  The maximum length is limited to `255` characters.

* `targets` - (Optional, List) Specifies the list of target objects.
  The [targets](#app_policy_group_targets) structure is documented below.
  
* `policies` - (Optional, String) Specifies the policies of the policy group, in JSON format.  
  For policy configuration items, please refer to the [documentation](https://support.huaweicloud.com/api-workspace/CreatePolicyGroup.html#CreatePolicyGroup__request_Policies).

<a name="app_policy_group_targets"></a>
The `targets` block supports:

* `id` - (Required, String) Specifies the object ID.  
  If the `type` is set to **USER**, the ID means the user ID.  
  If the `type` is set to **USERGROUP**, the ID means the user group ID.  
  If the `type` is set to **APPGROUP**, the ID means the APP group ID.  
  If the `type` is set to **OU**, the ID means the OU ID.  
  If the `type` is set to **ALL**, the ID fixed with string **default-apply-all-targets**.

* `name` - (Required, String) Specifies the object name.  
  If the `type` is set to **USER**, the name means the user name.  
  If the `type` is set to **USERGROUP**, the name means the user group name.  
  If the `type` is set to **APPGROUP**, the name means the APP group name.  
  If the `type` is set to **OU**, the name means the OU name.  
  If the `type` is set to **ALL**, the name fixed with string **All-Targets**.

* `type` - (Required, String) Specifies the object type.  
  The valid values are as follows:
  + **USER**
  + **USERGROUP**
  + **APPGROUP**
  + **OU**
  + **ALL**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The APP policy group can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_policy_group.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `priority`, `policies`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_app_policy_group" "test" {
  ...

  lifecycle {
    ignore_changes = [
      priority, policies,
    ]
  }
}
```
