---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_agency_permission_policy"
description: |-
  Manages an agency permission policy resource within HuaweiCloud.
---

# huaweicloud_gaussdb_agency_permission_policy

Manages an agency permission policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "bind_role_names" {
  type = list(string)
}
variable "unbind_role_names" {
  type = list(string)
}

resource "huaweicloud_gaussdb_agency_permission_policy" "test" {
  agency_name       = "RDSAccessProjectResource"
  bind_role_names   = var.bind_role_names
  unbind_role_names = var.unbind_role_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `agency_name` - (Required, String, NonUpdatable) Specifies the agency name.
  Currently, the value only can be **RDSAccessProjectResource**.

* `bind_role_names` - (Required, List) Specifies the permission policies to be bound from an agency.

* `unbind_role_names` - (Required, List) Specifies the permission policies to be unbound to an agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the `agency_name`.

* `is_existed` - Whether the agency exists.

* `roles` - Indicates policy information bound to the agency.
  The [roles](#roles_struct) structure is documented below.

<a name="roles_struct"></a>
The `roles` block supports:

* `name` - Indicates the policy name.

* `description` - Indicates the policy description.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_agency_permission_policy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `bind_role_names`, `unbind_role_names`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_agency_permission_policy" "test" {
  ...

  lifecycle {
    ignore_changes = [
      bind_role_names, unbind_role_names,
    ]
  }
}
```
