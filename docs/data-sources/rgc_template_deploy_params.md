---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_template_deploy_params"
description: |-
  Use this data source to get template deploy params in Resource Governance Center.
---

# huaweicloud_rgc_template_deploy_params

Use this data source to get template deploy params in Resource Governance Center.

## Example Usage

```hcl
variable template_name {}
variable template_version {}

data "huaweicloud_rgc_template_deploy_params" "test" {
  template_name    = var.template_name
  template_version = var.template_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `template_name` - (Required, String) The name of the template.

* `template_version` - (Required, String) The version of the template. It must follow the format `^V[1-9][0-9]{0,9}$`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `variables` - Information about the deployment template variables list.

The [variables](#variables) structure is documented below.

<a name="variables"></a>
The `variables` block supports:

* `name` - The name of the variable.

* `description` - The description of the variable.

* `type` - The type of the variable.

* `default` - The default value of the variable.

* `nullable` - Whether the variable can be null.

* `sensitive` - Whether the variable is a sensitive field.

* `validations` - Information about the validation rules for the variable list.

The [validations](#validations) structure is documented below.

<a name="validations"></a>
The `validations` block supports:

* `condition` - The validation expression.

* `error_message` - The error message if the validation fails.
