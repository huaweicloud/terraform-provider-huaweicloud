---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_application"
description: |-
  Manages an application resource within HuaweiCloud.
---

# huaweicloud_cae_application

Manages an application resource within HuaweiCloud.

## Example Usage

```hcl
variable "environment_id" {}
variable "application_name" {}

resource "huaweicloud_cae_application" "test" {
  environment_id = var.environment_id
  name           = var.application_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the environment to which the application
  belongs.  
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the application.  
  The valid length is limited from `2` to `63`, only lowercase letters, digits and hyphens (-) are allowed.  
  The name must start with a lowercase letter and end with a lowercase letter or a digit.  
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.

## Import

The application can be imported using `environment_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_cae_application.test <environment_id>/<id>
```
