---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_application"
description: ""
---

# huaweicloud_servicestage_application

Manages an application resource within HuaweiCloud ServiceStage.

## Example Usage

### Create an application and an environment variable

```hcl
variable "env_id" {}
variable "app_name" {}
variable "vpc_id" {}

resource "huaweicloud_servicestage_application" "test" {
  name = var.app_name

  environment {
    id = var.env_id

    variable {
      name  = "owner"
      value = "terraform"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ServiceStage application.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the application name.
  The name can contain `2` to `64` characters, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter and end with a letter or digit.

* `description` - (Optional, String) Specifies the application description.
  The description can contain a maximum of `128` characters.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the application
  belongs. Changing this will create a new resource.

* `environment` - (Optional, List) Specifies the configurations of the environment variables.
  The [object](#servicestage_app_environments) structure is documented below.

<a name="servicestage_app_environments"></a>
The `environment` block supports:

* `id` - (Required, String) Specifies the environment ID to which the variables belongs.

* `variable` - (Required, List) Specifies the list of environment variables.
  The [object](#servicestage_app_variables) structure is documented below.

<a name="servicestage_app_variables"></a>
The `variable` block supports:

* `name` - (Required, String) Specifies the variable name. The name can contain `1` to `64` characters, only letters,
  digits, underscores (_), hyphens (-) and dots (.) are allowed. The name cannot start with a digit or dot.

* `value` - (Required, String) Specifies the variable value. The value can contain a maximum of `2,048` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The application ID in UUID format.

* `component_ids` - The list of component IDs associated under the application.

## Import

Applications can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_servicestage_application.test eeea08e7-c838-4794-926c-abc12f3e10e8
```
