---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_tags_filter"
description: |-
  Use this data source to get the list of the Workspace desktops by tags within HuaweiCloud.
---

# huaweicloud_workspace_desktop_tags_filter

Use this data source to get the list of the Workspace desktops by tags within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_desktop_tags_filter" "test" {}
```

### Filter desktops without any tags

```hcl
data "huaweicloud_workspace_desktop_tags_filter" "test" {
  without_any_tag = true
}
```

### Filter desktops with all specified tags

```hcl
variable "tag_key" {}
variable "tag_value" {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  tags {
    key    = var.tag_key
    values = [ var.tag_value ]
  }
}
```

### Filter desktops with any of specified tags

```hcl
variable "tag_key" {}
variable "tag_value" {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  tags_any {
    key    = var.tag_key
    values = [ var.tag_value ]
  }
}
```

### Filter desktops without all specified tags

```hcl
variable "tag_key" {}
variable "tag_value" {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  not_tags {
    key    = var.tag_key
    values = [ var.tag_value ]
  }
}
```

### Filter desktops without any of specified tags

```hcl
variable "tag_key" {}
variable "tag_value" {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  not_tags_any {
    key    = var.tag_key
    values = [ var.tag_value ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region where the desktop tags are located.

* `without_any_tag` - (Optional, Bool) Specifies whether to query resources without any tag.

* `tags` - (Optional, List) The list of tags to filter desktops. Resources must  
  contain all specified tags.
  The [tag](#workspace_desktop_tag) structure is documented below.

* `tags_any` - (Optional, List) The list of tags to filter desktops. Resources must contain  
  at least one of specified tags.
  The [tag](#workspace_desktop_tag) structure is documented below.

* `not_tags` - (Optional, List) The list of tags to filter desktops. Resources must  
  not contain specified tags.
  The [tag](#workspace_desktop_tag) structure is documented below.

* `not_tags_any` - (Optional, List) The list of tags to filter desktops. Resources must  
  not contain any of specified tags.
  The [tag](#workspace_desktop_tag) structure is documented below.

* `matches` - (Optional, List) The list of matching rules to filter desktops.  
  The [match](#workspace_desktop_match) structure is documented below.

<a name="workspace_desktop_tag"></a>
The `tag` block supports:

* `key` - (Required, String) The key of tag.

* `values` - (Required, List) The value of tag.

<a name="workspace_desktop_match"></a>
The `match` block supports:

* `key` - (Required, String) The name of desktop property.

* `value` - (Required, String) The value of desktop property.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `desktops` - List of filtered desktops.  
  The [desktop object](#desktop_object) structure is documented below.

<a name="desktop_object"></a>
The `desktops` block supports:

* `resource_id` - The ID of the desktop.

* `resource_name` - The name of the desktop.

* `resource_detail` - The detail of the desktop.

* `tags` - The list of tags attached to the desktop.  
  The [desktop tag object](#desktop_tag_object) structure is documented below.

<a name="desktop_tag_object"></a>
The `tags` block in desktop object supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
