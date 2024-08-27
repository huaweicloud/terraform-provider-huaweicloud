---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_environment"
description: ""
---

# huaweicloud_apig_environment

Manages an APIG environment resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "environment_name" {}
variable "description" {}

resource "huaweicloud_apig_environment" "test" {
  instance_id = var.instance_id
  name        = var.environment_name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dedicated instance is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the environment
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the environment name.  
  The valid length is limited from `3` to `64`, only letters, digits and underscores (_) are allowed.
  The name must start with a letter.

* `description` - (Optional, String) Specifies the environment description.  
  The value can contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.
  Chinese characters must be in **UTF-8** or **Unicode** format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the dedicated environment.

* `created_at` - The time when the environment was created.

## Import

Environments can be imported using their `name` and the ID of the related dedicated instance, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_environment.test <instance_id>/<name>
```
