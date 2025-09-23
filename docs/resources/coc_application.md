---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_application"
description: |-
  Manage a COC application resource within HuaweiCloud.
---

# huaweicloud_coc_application

Manage a COC application resource within HuaweiCloud.

## Example Usage

```hcl
variable "application_name" {}

resource "huaweicloud_coc_application" "test" {
  name = var.application_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the application name.

* `parent_id` - (Optional, String, NonUpdatable) Specifies the parent ID.

* `description` - (Optional, String) Specifies the description.

* `is_collection` - (Optional, Bool) Specifies whether to add to collection. The default value is `false`.

  -> This parameter only takes effect when `parent_id` is empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `code` - Indicates the application code.

* `path` - Indicates the application level path.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the modification time.

## Import

The COC application can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_coc_application.test <id>
```
