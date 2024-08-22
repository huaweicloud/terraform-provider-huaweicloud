---
subcategory: "API Gateway (Shared APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_api_gateway_environment"
description: ""
---

# huaweicloud_api_gateway_environment

Manages a shared APIG environment resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_api_gateway_environment" "test_env" {
  name        = "test"
  description = "test env"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the shared APIG environment is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the environment name.
  The valid length is limited from `3` to `64`, only letters, digits and underscores (_) are allowed.
  The name must start with a letter.

* `description` - (Optional, String) Specifies the environment description.
  The value can contain a maximum of `255` characters.
  Chinese characters must be in **UTF-8** or **Unicode** format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The environment ID.

* `created_at` - The time when the shared APIG environment was created.

## Import

APIG environments can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_api_gateway_environment.test_env 774438a28a574ac8a496325d1bf51807
```
