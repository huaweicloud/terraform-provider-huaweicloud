---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_assignment_packages"
description: ""
---

# huaweicloud_rms_assignment_packages

Use this data source to get the list of RMS assignment packages.

## Example Usage

```hcl
variable "assignment_package_name" {}

data "huaweicloud_rms_assignment_packages" "test" {
  name   = var.assignment_package_name
  status = "CREATE_SUCCESSFUL"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the assignment package name. It contains 1 to 64 characters.

* `package_id` - (Optional, String) Specifies the assignment package ID.

* `status` - (Optional, String) Specifies the assignment package status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `packages` - The assignment package list.

  The [packages](#packages_struct) structure is documented below.

<a name="packages_struct"></a>
The `packages` block supports:

* `id` - The ID of an assignment package.

* `name` - The assignment package name.

* `deployment_id` - The deployment ID.

* `stack_name` - The name of a resource stack.

* `stack_id` - The unique ID of a resource stack.

* `status` - The deployment status of an assignment package.

* `error_message` - The error message when you failed to deploy or delete an assignment package.

* `created_by` - The creator of the assignment package.

* `created_at` - The time when the assignment package was created.

* `updated_at` - The time when the assignment package was updated.

* `vars_structure` - The parameters of the assignment package.

  The [vars_structure](#packages_vars_structure_struct) structure is documented below.

<a name="packages_vars_structure_struct"></a>
The `vars_structure` block supports:

* `var_key` - The name of a parameter.

* `var_value` - The value of a parameter.
