---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_application"
description: |-
  Manages an application resource within HuaweiCloud.
---

# huaweicloud_servicestagev3_application

Manages an application resource within HuaweiCloud.

## Example Usage

```hcl
variable "app_name" {}
variable "enterprise_project_id" {}

resource "huaweicloud_servicestagev3_application" "test" {
  name                  = var.app_name
  description           = "Created by terraform script"
  enterprise_project_id = var.enterprise_project_id

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the application.  
  The valid length is limited from `2` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter and end with a letter or a digit.

* `description` - (Optional, String) Specifies the description of the application.  
  The value can contain a maximum of `128` characters.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the application belongs.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the application that used to filter resource.

<a name="servicestage_v3_app_labels"></a>
The `labels` block supports:

* `key` - (Required, String) Specifies the label key.

* `value` - (Required, String) Specifies the label value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `creator` - The creator name of the application.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.

## Import

Applications can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_servicestagev3_application.test <id>
```
