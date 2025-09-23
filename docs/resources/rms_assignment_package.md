---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_assignment_package"
description: ""
---

# huaweicloud_rms_assignment_package

Manages a RMS assignment package resource within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_rms_assignment_package_templates" "test" {}

resource "huaweicloud_rms_assignment_package" "test" {
  name         = "test"
  template_key = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key

  dynamic "vars_structure" {
    for_each = data.huaweicloud_rms_assignment_package_templates.test.templates.0.parameters
    content {
      var_key = vars_structure.value["name"]
      var_value = vars_structure.value["default_value"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the assignment package name. It contains 1 to 64 characters.

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name. The agency needs to authorize RFS to invoke
  the Config APIs for creating, updating, deleting an assignment. It contains 1 to 64 characters.

  Changing this parameter will create a new resource.

* `template_key` - (Optional, String, ForceNew) Specifies the name of a built-in assignment package template. It
  contains 1 to 128 characters.

  Changing this parameter will create a new resource.

* `template_body` - (Optional, String, ForceNew) Specifies the content of a custom assignment package. It contains 1 to
  51200 characters.

  Changing this parameter will create a new resource.

* `template_uri` - (Optional, String, ForceNew) Specifies the URL address of the OBS bucket where an assignment package
  template was stored. It contains 1 to 1024 characters.

  Changing this parameter will create a new resource.

  -> **NOTE:** Exactly one of `template_key`, `template_body`, `template_uri` should be specified.

* `vars_structure` - (Optional, List) Specifies the parameters of an assignment package.

  The [vars_structure](#AssignmentPackage_VarStructure) structure is documented below.

<a name="AssignmentPackage_VarStructure"></a>
The `vars_structure` block supports:

* `var_key` - (Optional, String) Specifies the name of a parameter. It contains 1 to 64 characters.

* `var_value` - (Optional, String) Specifies the value of a parameter. It's a json string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `stack_id` - Indicates the unique ID of a resource stack.

* `stack_name` - Indicates the name of a resource stack.

* `deployment_id` - Indicates the deployment ID.

* `status` - Indicates the deployment status of an assignment package.

## Import

The RMS assignment package can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rms_assignment_package.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `agency_name`, `template_key`,
`template_body`, `template_uri`. It is generally recommended running `terraform plan` after importing a RMS assignment
package. You can then decide if changes should be applied to the RMS assignment package, or the resource definition
should be updated to align with the RMS assignment package. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rms_assignment_package" "test" {
    ...

  lifecycle {
    ignore_changes = [
      agency_name, template_key, template_body, template_uri,
    ]
  }
}
```
