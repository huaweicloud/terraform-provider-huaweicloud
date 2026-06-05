---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_ai_api_key"
description: |-
  Manages an AI API key in application resource within HuaweiCloud.
---

# huaweicloud_apig_application_ai_api_key

Manages an AI API key in application resource within HuaweiCloud.

-> The `ai_api_key_enabled` feature must be enabled before the AI API key resource create.

## Example Usage

### Auto generate the AI API key value

```hcl
variable "instance_id" {}
variable "application_id" {}

resource "huaweicloud_apig_application_ai_api_key" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
}
```

### Manually configure the AI API key value

```hcl
variable "instance_id" {}
variable "application_id" {}

resource "huaweicloud_apig_application_ai_api_key" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
  value          = "bMQBWYxyln0J5QsePabHtDzN3h2n22HdzEcZbPfQNgs0tBytH5oCqcNFz7RKYcoU7lz44oo9H9uc7XeD3O5rjA"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the AI API key is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the application
  and AI API key belong.

* `application_id` - (Required, String, NonUpdatable) Specifies the ID of the application to which the AI API key
  belongs.

* `alias` - (Optional, String, NonUpdatable) Specifies the alias of the AI API key.  
  The value can contain `1` to `100` characters, including uppercase and lowercase letters, digits, underscores (_),
  and hyphens (-).

* `value` - (Optional, String, NonUpdatable) Specifies the value of the AI API key.  
  The value can contain `8` to `128` characters, including uppercase and lowercase letters, digits, and the
  following special characters: `+_-/=`.  
  If omitted, a random value will be generated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the AI API key.

* `created_at` - The creation time of the AI API key, in RFC3339 format.

## Import

AI API keys can be imported using related `instance_id`, `application_id` and their `id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_apig_application_ai_api_key.test <instance_id>/<application_id>/<id>
```
