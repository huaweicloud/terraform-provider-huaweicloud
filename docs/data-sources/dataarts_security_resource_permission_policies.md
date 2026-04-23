---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_resource_permission_policies"
description: |-
  Use this data source to get the list of DataArts Security resource permission policies within HuaweiCloud.
---

# huaweicloud_dataarts_security_resource_permission_policies

  Use this data source to get the list of DataArts Security resource permission policies within HuaweiCloud.

## Example Usage

### Query all resource permission policies

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_resource_permission_policies" "test" {
  workspace_id = var.workspace_id
}
```

### Query resource permission policies by name

```hcl
variable "workspace_id" {}
variable "policy_name" {}

data "huaweicloud_dataarts_security_resource_permission_policies" "test" {
  workspace_id = var.workspace_id
  policy_name  = var.policy_name
}
  ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource permission policies.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the resource permission policies belong.

* `policy_name` - (Optional, String) Specifies the name of the resource permission policy to be queried.  
  Fuzzy matching is supported.

* `resource_name` - (Optional, String) Specifies the name of the authorized resource to be queried.  
  Fuzzy matching is supported.

* `member_name` - (Optional, String) Specifies the name of the authorized member to be queried.  
  Fuzzy matching is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of resource permission policies that matched filter parameters.  
  The [policies](#dataarts_security_resource_permission_policies_attr) structure is documented below.

<a name="dataarts_security_resource_permission_policies_attr"></a>
The `policies` block supports:

* `policy_id` - The ID of the resource permission policy.

* `policy_name` - The name of the resource permission policy.

* `resources` - The resource list of the resource permission policy.  
  The [resources](#dataarts_security_resource_permission_policies_resources) structure is documented below.

* `members` - The member list of the resource permission policy.  
  The [members](#dataarts_security_resource_permission_policies_members) structure is documented below.

* `create_time` - The creation time of the resource permission policy, in RFC3339 format.

* `create_user` - The creator of the resource permission policy.

* `update_time` - The latest update time of the resource permission policy, in RFC3339 format.

<a name="dataarts_security_resource_permission_policies_resources"></a>
The `resources` block supports:

* `resource_id` - The ID of the resource.

* `resource_name` - The name of the resource.

* `resource_type` - The type of the resource.
  + **DATA_CONNECTION**
  + **AGENCY**

  <a name="dataarts_security_resource_permission_policies_members"></a>
  The `members` block supports:

* `member_id` - The ID of the member.

* `member_name` - The name of the member.

* `member_type` - The type of the member.
  + **USER**
  + **USER_GROUP**
  + **WORKSPACE_ROLE**
