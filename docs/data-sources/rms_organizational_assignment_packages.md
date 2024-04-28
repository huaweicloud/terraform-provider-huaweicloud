---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_organizational_assignment_packages"
description: ""
---

# huaweicloud_rms_organizational_assignment_packages

Use this data source to get the list of RMS organizational assignment packages.

## Example Usage

```hcl
variable "assignment_package_name" {}

data "huaweicloud_organizations_organization" "test" {}

data "huaweicloud_rms_organizational_assignment_packages" "test" {
  organization_id = data.huaweicloud_organizations_organization.test.id
  name            = var.assignment_package_name
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required, String) Specifies the ID of the organization.

* `name` - (Optional, String) Specifies the organizational assignment package name.

* `package_id` - (Optional, String) Specifies the organizational assignment package ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `packages` - The list of organizational assignment packages.

  The [packages](#packages_struct) structure is documented below.

<a name="packages_struct"></a>
The `packages` block supports:

* `name` - The organizational assignment package name.

* `id` - The organizational assignment package ID.

* `organization_id` - The ID of the organization.

* `owner_id` - The creator of the organizational assignment package.

* `org_conformance_pack_urn` - The unique identifier of the organizational assignment package.

* `vars_structure` - The parameters of the organizational assignment package.

  The [vars_structure](#packages_vars_structure_struct) structure is documented below.

* `created_at` - The creation time of the organizational assignment package.

* `updated_at` - The update time of the organizational assignment package.

<a name="packages_vars_structure_struct"></a>
The `vars_structure` block supports:

* `var_key` - The name of a parameter.

* `var_value` - The value of a parameter.
