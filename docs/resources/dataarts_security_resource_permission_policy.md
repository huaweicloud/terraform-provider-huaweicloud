---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_resource_permission_policy"
description: |-
  Manages a DataArts Studio Security resource permission policy resource within HuaweiCloud.
---


# huaweicloud_dataarts_security_resource_permission_policy

Manages a DataArts Studio Security resource permission policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "policy_name" {}
variable "resources" {
  type = list(object({
    resource_id   = string
    resource_name = string
    resource_type = string
  }))
}

variable "members" {
  type = list(object({
    member_id   = string
    member_name = string
    member_type = string
  }))
}

resource "huaweicloud_dataarts_security_resource_permission_policy" "test" {
  workspace_id = var.workspace_id
  name         = var.policy_name

  dynamic "resources" {
    for_each = var.resources

    content {
      resource_id   = resources.value.resource_id
      resource_name = resources.value.resource_name
      resource_type = resources.value.resource_type
    }
  }

  dynamic "members" {
    for_each = var.members

    content {
      member_id   = members.value.member_id
      member_name = members.value.member_name
      member_type = members.value.member_type
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the resource permission
  policy belongs.

* `name` - (Required, String) Specifies the name of the resource permission policy.  
  The name can contain `2` to `64` characters.  
  Only letters, Chinese characters, digits, and underscores (_) are allowed, and must start with
  a letter or Chinese character.  

* `resources` - (Required, List) Specifies the list of resources.  
  The [resources](#resource_permission_policy_resources) structure is documented below.

* `members` - (Required, List) Specifies the list of members.  
  The [members](#resource_permission_policy_members) structure is documented below.

<a name="resource_permission_policy_resources"></a>
The `resources` block supports:

* `resource_id` - (Required, String) Specifies the ID of the resource.

* `resource_name` - (Required, String) Specifies the name of the resource.

* `resource_type` - (Required, String) Specifies the type of the resource.  
  The valid values are as follows:
  + **DATA_CONNECTION**
  + **AGENCY**

<a name="resource_permission_policy_members"></a>
The `members` block supports:

* `member_id` - (Required, String) Specifies the ID of the member.

* `member_name` - (Required, String) Specifies the name of the member.

* `member_type` - (Required, String) Specifies the type of the member.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**
  + **WORKSPACE_ROLE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the resource permission policy, in RFC3339 format.

* `create_user` - The creator of the resource permission policy.

* `update_time` - The latest update time of the resource permission policy, in RFC3339 format.

## Import

The resource can be imported using `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_security_resource_permission_policy.test <workspace_id>/<id>
```
