---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_roles"
description: |-
  Use this data source to query the list of IAM roles within HuaweiCloud.
---

# huaweicloud_identity_roles

Use this data source to query the list of IAM roles within HuaweiCloud.

-> You *must* have IAM read privileges to use this data source.

## Example Usage

### Query All System Roles

```hcl
data "huaweicloud_identity_roles" "all" {
}
```

### Query Roles by Display Name

```hcl
variable "display_name" {}

data "huaweicloud_identity_roles" "by_display_name" {
  display_name = var.display_name
}
```

### Query Roles by Catalog

```hcl
data "huaweicloud_identity_roles" "by_catalog" {
  catalog = "ELB"
}
```

### Query Custom Policies

```hcl
variable "domain_id" {}

data "huaweicloud_identity_roles" "custom" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional, String) Specifies the display name of the role to be queried.  
  This parameter can be used to filter roles by permission name:
  + **Permission name**: If you set this parameter to **ECS FullAccess**, the information about this permission
    will be returned.
  + **Filter condition**: If you set this parameter to **Administrator**, all administrator permissions that match
    the condition will be returned.

* `name` - (Optional, String) Specifies the name of the role to be queried.  
  The system internal name of the permission.  
  For example, the name of **CCS User** permission is **ccs_user**.  
  It is recommended to use the `display_name` parameter.

* `catalog` - (Optional, String) Specifies the catalog of the role to be queried.  
  The service catalog of the permission.

* `type` - (Optional, String) Specifies the display mode of the role to be queried.  
  The valid values are as follows:
  + **domain**: Returns roles with type **AA** and **AX**.
  + **project**: Returns roles with type **AA** and **XA**.
  + **all**: Returns roles with type **AA**, **AX** and **XA**.

  -> The type values in the response have the following meanings:<br>
     **AX** indicates displayed at the domain level;<br>
     **XA** indicates displayed at the project level;<br>
     **AA** indicates displayed at both domain and project levels;<br>
     **XX** indicates not displayed at either domain or project level.

* `permission_type` - (Optional, String) Specifies the permission type of the role to be queried.  
  The valid values are as follows:
  + **policy**: Returns system policies.
  + **role**: Returns system roles.

  -> This parameter takes effect when the `domain_id` parameter is empty.

* `domain_id` - (Optional, String) Specifies the domain ID to be queried.  
  If this parameter is specified, only custom policies of the account are returned.  
  If this parameter is not specified, all system permissions (including system policies and system roles) are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `roles` - The list of the roles.
  The [roles](#identity_roles_attr) structure is documented below.

<a name="identity_roles_attr"></a>
The `roles` block supports:

* `id` - The ID of the role.

* `name` - The name of the role.

* `display_name` - The display name of the role.

* `description` - The description of the role.

* `description_cn` - The description of the role in Chinese.

* `catalog` - The catalog of the role.

* `type` - The display mode of the role.

* `flag` - The flag of the role.

* `domain_id` - The domain ID of the role. This field is only returned for custom policies.

* `policy` - The content of the role, in JSON format.

* `created_at` - The creation time of the role, in RFC339 format.

* `updated_at` - The last update time of the role, in RFC339 format.

* `links` - The links of the role.
  The [links](#identity_roles_links_attr) structure is documented below.

<a name="identity_roles_links_attr"></a>
The `links` block supports:

* `self` - The self link of the role.

* `previous` - The previous link of the role.

* `next` - The next link of the role.
