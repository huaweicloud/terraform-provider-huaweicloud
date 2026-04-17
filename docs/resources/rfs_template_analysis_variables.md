---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_template_analysis_variables"
description: |-
  Manages a resource to analyze template variables within HuaweiCloud.
---

# huaweicloud_rfs_template_analysis_variables

Manages a resource to analyze template variables within HuaweiCloud.

-> This resource is a one-time action resource used to analyze template variables. Deleting this resource will
  not cancel the analysis operation, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "template_body" {}

resource "huaweicloud_rfs_template_analysis_variables" "test" {
  template_body = var.template_body
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `template_body` - (Optional, String) Specifies the HCL template that describes the target state of resources.

* `template_uri` - (Optional, String) Specifies the OBS address of the HCL template that describes the target state of
  resources. Please ensure that the OBS address is in the same region as the RFS service.
  The corresponding file should be a pure tf file or zip compressed package.
  A pure .tf file must end with .tf or .tf.json and comply with the HCL syntax.
  Currently, only the .zip package is supported. The file name extension must be .zip. The decompressed file cannot
  contain the .tfvars file and must be encoded in UTF8 format (the .tf.json file cannot contain the BOM header).
  The .zip package supports a maximum of `100` subfiles.

-> Either `template_body` or `template_uri` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `variables` - The list of variables parsed from the template.

  The [variables](#variables_struct) structure is documented below.

<a name="variables_struct"></a>
The `variables` block supports:

* `name` - The name of the variable.

* `type` - The type of the variable.

* `description` - The description of the variable.

* `default` - The default value of the variable. The type of the return value is the same as that defined in the
  `type` field.

* `sensitive` - Whether the variable is a sensitive field.
  If `sensitive` is not defined in the variable, **false** is returned by default.

* `nullable` - Whether the variable can be set to null.
  If `nullable` is not defined in the variable, **true** is returned by default.

* `validations` - The validation rules of the variable.

  The [validations](#validations_struct) structure is documented below.

<a name="validations_struct"></a>
The `validations` block supports:

* `condition` - The validation expression.

* `error_message` - The error message when validation fails.
