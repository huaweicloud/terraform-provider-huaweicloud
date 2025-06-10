---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_tags"
description: |-
  Use this data source to query the tags of a specified Workspace desktop within HuaweiCloud.
---

# huaweicloud_workspace_desktop_tags

Use this data source to query the tags of a specified Workspace desktop within HuaweiCloud.

## Example Usage

```hcl
variable "desktop_id" {}

data "huaweicloud_workspace_desktop_tags" "test" {
  desktop_id = var.desktop_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region where the Workspace desktop is located.

* `desktop_id` - (Required, String) The ID of the desktop to query tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.  
  The [tag](#workspace_desktop_tag) structure is documented below.

<a name="workspace_desktop_tag"></a>
The `tag` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
