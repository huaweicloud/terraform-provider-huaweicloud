---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_application_configuration"
description: |-
  Manages the application confiration for the environment within HuaweiCloud.
---

# huaweicloud_servicestagev3_application_configuration

Manages the application confiration for the environment within HuaweiCloud.

## Example Usage

```hcl
variable "environment_id" {}
variable "application_id" {}
variable "environment_variables" {
  type = list(object({
    name  = string
    value = string
  }))
}

resource "huaweicloud_servicestagev3_application_configuration" "test" {
  environment_id = var.environment_id
  application_id = var.application_id

  configuration {
    dynamic "env" {
      for_each = var.environment_variables

      content {
        name  = env.value["name"]
        value = env.value["value"]
      }
    }
    assign_strategy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the environment and application are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `environment_id` - (Required, String, NonUpdatable) Specifies the ID of the environment to which the configuration
  applies.

* `application_id` - (Required, String, NonUpdatable) Specifies the ID of the application to which the configuration
  belongs.

* `configuration` - (Required, List) Specifies the configuration of the application.  
  The [configuration](#servicestage_v3_application_configuration) structure is documented below.

<a name="servicestage_v3_application_configuration"></a>
The `configuration` block supports:

* `env` - (Required, List) Specifies the list of the environment variables.  
  The [env](#servicestage_v3_application_configuration_env) structure is documented below.

* `assign_strategy` - (Optional, Bool) Specifies whether the effective strategy is the continuously effective.  
  The valid values are as follows:
  + **true**: First time effective. Application-level environment variables only take effect when the component is
    first created, and subsequent modifications of the application-level environment variables will not be synchronized
    with the environment variables in the component.
  + **false**: Continuously effective. Environment variables during component upgrades are updated according to the
    application-level environment variables.

  Defaults to **false**.

<a name="servicestage_v3_application_configuration_env"></a>
The `env` block supports:

* `name` - (Required, String) Specifies the name of the environment variable.  
  The valid length is limited from `1` to `64`, only Chinese characters, English letters, digits, hyphens (-),
  underscores (\_) and dots (.) are allowed.  
  The name must start with an English letter, hyphen (-) or underscore (\_).

  -> Variable names must be unique within the same application environment.

* `value` - (Required, String) Specifies the value of the environment variable.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consists of `environment_id` and `application_id`.

## Import

Application configuration can be imported using `environment_id` and `application_id`, e.g.

```bash
$ terraform import huaweicloud_servicestagev3_application_configuration.test <environment_id>/<application_id>
```
