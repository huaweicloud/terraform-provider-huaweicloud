---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_tags"
description: |-
  Use this data source to query the desktop tags under a specified region within HuaweiCloud.
---

# huaweicloud_workspace_tags

Use this data source to query the desktop tags under a specified region within HuaweiCloud.

## Example Usage

### Querying tag list under all desktops

```hcl
data "huaweicloud_workspace_tags" "test" {}
```

### Querying all desktop tags that are equal to the key name

```hcl
variable "key_name" {}

data "huaweicloud_workspace_tags" "test" {
  key = var.key_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the desktop tags.

* `key` - (Optional, String) The key of the tag to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of workspace tags.  
  The [tag](#workspace_desktop_tag) structure is documented below.

<a name="workspace_desktop_tag"></a>
The `tag` block supports:

* `key` - The key of the tag.

* `values` - The values of the tag.
