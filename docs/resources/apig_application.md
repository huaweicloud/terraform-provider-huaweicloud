---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application"
description: ""
---

# huaweicloud_apig_application

Manages an APIG application resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "app_name" {}
variable "app_code" {}

resource "huaweicloud_apig_application" "test" {
  instance_id = var.instance_id
  name        = var.app_name
  description = "Created by script"

  app_codes = [var.app_code]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the application name.  
  The valid length is limited from `3` to `64`, only Chinese characters, English letters, digits and hyphens (-)
  are allowed.  
  The name must start with a Chinese character or English letter.

* `description` - (Optional, String) Specifies the application description.  
  The description contain a maximum of `255` characters and the angle brackets (< and >) are not allowed.

  -> The description does not support updating to an empty value.

* `app_codes` - (Optional, List) Specifies an array of one or more application codes that the application has.  
  Up to five application codes can be created.  
  The valid length of each application code is limited from can contain `64` to `180`.  
  The application code must start with a letter, digit, plus sign (+) or slash (/).  
  Only letters, digits and following special special characters are allowed: `!@#$%+-_/=`.

* `secret_action` - (Optional, String) Specifies the secret action to be done for the application.  
  The valid action is **RESET**.

  -> The `secret_action` is a one-time action.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The application ID.

* `registration_time` - the registration time.

* `updated_at` - The latest update time of the application.

* `app_key` - App key.

* `app_secret` - App secret.

## Import

Applications can be imported using their `id` and the ID of the related dedicated instance, separated by a slash, e.g.

```shell
$ terraform import huaweicloud_apig_application.test <instance_id>/<id>
```
