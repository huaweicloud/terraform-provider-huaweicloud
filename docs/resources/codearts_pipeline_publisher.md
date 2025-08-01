---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_publisher"
description: |-
  Manages a CodeArts pipeline publisher resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_publisher

Manages a CodeArts pipeline publisher resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "en_name" {}
variable "description" {}
variable "website" {}
variable "support_url" {}
variable "source_url" {}

resource "huaweicloud_codearts_pipeline_publisher" "test" {
  name        = var.name
  en_name     = var.en_name
  description = var.description
  website     = var.website
  support_url = var.support_url
  source_url  = var.source_url
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `en_name` - (Required, String, NonUpdatable) Specifies the publisher English name.

* `name` - (Required, String, NonUpdatable) Specifies the publisher name.

* `support_url` - (Required, String, NonUpdatable) Specifies the support URL.

* `description` - (Optional, String, NonUpdatable) Specifies the description.

* `logo_url` - (Optional, String, NonUpdatable) Specifies the logo URL.

* `source_url` - (Optional, String, NonUpdatable) Specifies the source URL.

* `website` - (Optional, String, NonUpdatable) Specifies the website URL.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `auth_status` - Indicates the authorization status.

* `last_update_time` - Indicates the update time.

* `last_update_user_id` - Indicates the updater ID.

* `last_update_user_name` - Indicates the updater name.

* `user_id` - Indicates the user ID.

## Import

The publisher can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_publisher.test <id>
```
