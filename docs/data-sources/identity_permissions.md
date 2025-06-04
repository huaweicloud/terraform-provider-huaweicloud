---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_permissions"
description: ""
---

# huaweicloud_identity_permissions

Use this data source to get the available IAM [Permissions](https://support.huaweicloud.com/intl/en-us/productdesc-iam/iam_01_0023.html#section5),
including roles and policies.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

### Full Permissions of ECS Service

```hcl
data "huaweicloud_identity_permissions" "ecs_full" {
  name = "ECS FullAccess"
}
```

### All System Defined Permissions of ELB Service

```hcl
data "huaweicloud_identity_permissions" "elb_all" {
  catalog = "ELB"
}
```

### All Administrator Permissions

```hcl
data "huaweicloud_identity_permissions" "all_adm" {
  name = "Administrator"
}
```

### All Custom Policies

```hcl
data "huaweicloud_identity_permissions" "custom" {
  type = "custom"
}
```

### All Project Level System Defined Permissions

```hcl
data "huaweicloud_identity_permissions" "project_all" {
  scope_type = "project"
}
```

## Argument Reference

* `type` - (Optional, String) Specifies the type of the permission. The default value is **system**, and valid values are
  as follows:
  + **system**: The system-defined permissions (including system-defined policies and roles).
    We can get all **System-defined Permissions** from [HuaweiCloud](https://support.huaweicloud.com/intl/en-us/usermanual-permissions/iam_01_0001.html).
  + **system-policy**: The system-defined policies.
  + **system-role**: The system-defined roles.
  + **custom**: The custom policies.

* `scope_type` - (Optional, String) Specifies the mode of the permission. Valid values are
  as follows:
  + **domain**: All permissions of the AA and AX levels.
  + **project**: All permissions of the AA and XA levels
  + **all**: Permissions of the AA, AX, and XA permissions.

  Note:
  + **AX**: Account level.
  + **XA**: Project level.
  + **AA**: Both the account level and project level.
  + **XX**: Neither the account level nor project level.

* `catalog` - (Optional, String) Specifies the service catalog of the permission.

* `name` - (Optional, String) Specifies the permission name or filter condition.
  + Permission name: For example, if you set this parameter to **ECS FullAccess**, information about the permission will
    be returned.
  + Filter condition: For example, if you set this parameter to **Administrator**, all administrator permissions that
    meet the conditions will be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permissions` - An array of available permissions. The structure is documented below.

The `permissions` block supports:

* `id` - The permission ID.
* `name` - The permission name.
* `description` - The description of the permission.
* `description_cn` - The description of the permission in Chinese.
* `catalog` - The service catalog of the permission.
* `policy` - The content of the permission.
