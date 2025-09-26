---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_roles"
description: ""
---

# huaweicloud_identity_roles

Use this data source to get a list of IAM roles with optional filtering by name or display name.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

The Roles in Terraform correspond to the IAM roles in HuaweiCloud. You can retrieve all **System-Defined Roles** and their details using this data source.

## Example Usage

### Retrieve Roles by Name
```hcl
data "huaweicloud_identity_roles" "roles" {
  name = "system_all_64"
}
```

### Retrieve Roles by Display Name
```hcl
data "huaweicloud_identity_roles" "roles" {
  display_name = "OBS ReadOnlyAccess"
}
```

### Retrieve All Roles (No Filtering)
```hcl
data "huaweicloud_identity_roles" "roles" {}
```

## Argument Reference

* `name` - (Optional, String) Specifies the internal name of the role. If provided, only roles matching this name will be returned.

* `display_name` - (Optional, String) Specifies the display name of the role as shown on the console. If provided, only roles matching this display name will be returned.

-> **NOTE:** At least one of `name` or `display_name` must be specified if you want to filter roles. If neither is specified, all roles will be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `roles` - A list of roles matching the specified criteria. Each role contains the following attributes:
  * `id` - The unique identifier of the role in UUID format.
  * `name` - The internal name of the role.
  * `display_name` - The display name of the role as shown on the console.
  * `description` - The description of the role.
  * `catalog` - The service catalog associated with the role.
  * `type` - The type or display mode of the role.
  * `policy` - The JSON content of the role's policy.
```