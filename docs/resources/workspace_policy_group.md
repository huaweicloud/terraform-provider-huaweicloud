---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_policy_group"
description: ""
---

# huaweicloud_workspace_policy_group

Manages policy group resource within HuaweiCloud.

## Example Usage

### Create a policy group and allow two IP address access

```hcl
variable "policy_group_name" {}
variable "workspace_user_id" {}
variable "workspace_user_name" {}

resource "huaweicloud_workspace_policy_group" "test" {
  name     = var.policy_group_name
  priority = 1

  targets {
    type = "USER"
    id   = var.workspace_user_id
    name = var.workspace_user_name
  }
  policy {
    access_control {
      ip_access_control = "112.20.53.2|255.255.240.0;112.20.53.3|255.255.240.0"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the policy group is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the policy group name.

* `priority` - (Optional, Int) Specifies the priority of the policy group.

* `description` - (Optional, String) Specifies the description of the policy group.

* `targets` - (Optional, List) Specifies the configuration of the access targets.
  The [targets](#policy_group_targets) structure is documented below.

* `policy` - (Optional, List) Specifies the configuration of the access policy.
  The [policy](#policy_group_policy) structure is documented below.

<a name="policy_group_targets"></a>
The `targets` block supports:

* `type` - (Required, String) Specifies the target type.
  The valid values are as follows:
  + **INSTANCE**: Desktop.
  + **USER**: User.
  + **USERGROUP**: User group.
  + **CLIENTIP**: Terminal IP address.
  + **OU**: Organization unit.
  + **ALL**: All desktops.

* `id` - (Required, String) Specifies the target ID.  
  If the `targets` type is **INSTANCE**, the ID means the SID of the desktop.  
  If the `targets` type is **USER**, the ID means the user ID.  
  If the `targets` type is **USERGROUP**, the ID means the user group ID.  
  If the `targets` type is **CLIENTIP**, the ID means the terminal IP address.  
  If the `targets` type is **OU**, the ID means the OUID.  
  If the `targets` type is **ALL**, the ID fixed with string **default-apply-all-targets**.

* `name` - (Required, String) Specifies the target name.  
  If the `targets` type is **INSTANCE**, the ID means the desktop name.  
  If the `targets` type is **USER**, the ID means the user name.  
  If the `targets` type is **USERGROUP**, the ID means the user group name.  
  If the `targets` type is **CLIENTIP**, the ID means the terminal IP address.  
  If the `targets` type is **OU**, the ID means the OU name.  
  If the `targets` type is **ALL**, the ID fixed with string **All-Targets**.

<a name="policy_group_policy"></a>
The `policy` block supports:

* `access_control` - (Required, List) Specifies the configuration of the access policy control.
  The [access_control](#policy_group_access_control) structure is documented below.

<a name="policy_group_access_control"></a>
The `access_control` block supports:

* `ip_access_control` - (Required, String) Specifies the IP access control.  
  It consists of multiple groups of IP addresses and network masks, separated by ';',
  and spliced together by '|' between IP addresses and network masks, e.g. `IP|mask;IP|mask;IP|mask`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The policy group ID in UUID format.

* `updated_at` - The latest update time of the policy group.

## Import

Policy groups can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_policy_group.test <id>
```
