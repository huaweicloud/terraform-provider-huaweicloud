---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_organizational_assignment_package"
description: ""
---

# huaweicloud_rms_organizational_assignment_package

Manages RMS organizational assignment package resource within HuaweiCloud.

## Example Usage

```hcl
variable "organization_id" {}

data "huaweicloud_rms_assignment_package_templates" "test" {}

resource "huaweicloud_rms_organizational_assignment_package" "test" {
  organization_id = var.organization_id
  name            = "test"
  template_key    = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key

  dynamic "vars_structure" {
    for_each = data.huaweicloud_rms_assignment_package_templates.test.templates.0.parameters
    content {
      var_key   = vars_structure.value["name"]
      var_value = vars_structure.value["default_value"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required, String, ForceNew) Specifies the ID of the organization.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the assignment package name. It contains `1` to `64` characters.

* `excluded_accounts` - (Optional, List) Specifies the excluded accounts for conformance package deployment.

* `template_key` - (Optional, String, ForceNew) Specifies the name of a predefined conformance package. It contains `1`
  to `128` characters.

  Changing this parameter will create a new resource.

* `template_body` - (Optional, String, ForceNew) Specifies the content of a custom assignment package. It contains `1`
  to `51200` characters.

  Changing this parameter will create a new resource.

* `template_uri` - (Optional, String, ForceNew) Specifies the OBS address of a conformance package. It contains `1` to
  `1024` characters.

  Changing this parameter will create a new resource.

* `vars_structure` - (Optional, List) Specifies the parameters of a conformance package.

The [vars_structure](#OrgAssignmentPackage_VarStructure) structure is documented below.

<a name="OrgAssignmentPackage_VarStructure"></a>
The `vars_structure` block supports:

* `var_key` - (Optional, String) Specifies the name of a parameter. It contains `1` to `64` characters.

* `var_value` - (Optional, String) Specifies the value of a parameter. It's a json string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `owner_id` - Indicates the creator of an organization conformance package.

* `org_conformance_pack_urn` - Indicates the unique identifier of resources in an organization conformance package.

* `created_at` - Indicates the creation time of an organization conformance package.

* `updated_at` - Indicates the latest update time of an organization conformance package.

## Import

The RMS organizational assignment package can be imported using the `organization_id` and `id` separated by a slash,
e.g.

```bash
$ terraform import huaweicloud_rms_originational_assignment_package.test <organization_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `template_key`,`template_body`,
`template_uri`. It is generally recommended running `terraform plan` after importing a RMS organizational assignment
package. You can then decide if changes should be applied to the RMS organizational assignment package, or the resource
definition should be updated to align with the RMS organizational assignment package. Also you can ignore changes as
below.

```hcl
resource "huaweicloud_rms_organizational_assignment_package" "test" {
    ...

  lifecycle {
    ignore_changes = [
      template_key, template_body, template_uri,
    ]
  }
}
```
