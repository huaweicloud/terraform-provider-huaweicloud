---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_notification_templates"
description: |-
  Use this data source to get the list of LTS notification templates.
---

# huaweicloud_lts_notification_templates

Use this data source to get the list of LTS notification templates.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_lts_notification_templates" "test" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_id` - (Required, String) Specified the account ID of IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - All notification templates that match the filter parameters.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `name` - The name of the notification template.

* `source` - The source of the notification template.

* `locale` - The language of the notification template.

* `templates` - The list of notification template bodies.

  The [templates](#templates_templates_struct) structure is documented below.

* `description` - The description of the notification template.

* `created_at` - The creation time of the log group, in RFC3339 format.

* `updated_at` - The latest update time of the log group, in RFC3339 format.

<a name="templates_templates_struct"></a>
The `templates` block supports:

* `sub_type` - The type of the template body.

* `content` - The content of the template body.
