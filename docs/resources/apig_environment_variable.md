---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_environment_variable"
description: |-
  Manages an APIG environment variable resource within HuaweiCloud.
---

# huaweicloud_apig_environment_variable

Manages an APIG environment variable resource within HuaweiCloud.

-> A maximum of `50` variable can be created on the same environment.

## Example Usage

```hcl
variable "instance_id" {}
variable "environment_id" {}
variable "group_id" {}
variable "variable_name" {}
variable "variable_value" {}

resource "huaweicloud_apig_environment_variable" "test" {
  instance_id = var.instance_id
  env_id      = var.environment_id
  group_id    = var.group_id
  name        = var.variable_name
  value       = var.variable_value
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the environment
  variable belongs. Changing this creates a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the group to which the environment variable belongs.
  Changing this creates a new resource.

* `env_id` - (Required, String, ForceNew) Specifies the ID of the environment to which the environment variable belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the environment variable.
  Changing this creates a new resource.
  The valid length is limited from `3` to `32` characters.
  Only letters, digits, hyphens (-), and underscores (_) are allowed, and must start with a letter.
  In the definition of an API, the`name` (case-sensitive) indicates a variable, for example, `#Name#`.
  It is replaced by the actual value when the API is published in an environment. The variable name must be unique.

* `value` - (Required, String) Specifies the value of the environment variable.
  The valid length is limited from `1` to `255` characters. Only letters, digits and special characters (_-/.:) are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using `instance_id`, `group_id` and `name`, separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_apig_environment_variable.test <instance_id>/<group_id>/<name>
```
