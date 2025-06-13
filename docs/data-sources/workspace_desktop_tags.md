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

* `region` - (Optional, String) Specifies the region where the Workspace desktop is located.  
  If omitted, the provider-level region will be used.

* `desktop_id` - (Required, String) Specifies the ID of the desktop to query tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.  
  The [tag](#workspace_tags_tag) structure is documented below.

<a name="workspace_tags_tag"></a>
The `tag` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
