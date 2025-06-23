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
variable "filter_tag_key" {}
variable "filter_tag_values" {
  type = list(string)
}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  tags {
    key    = var.filter_tag_key
    values = var.filter_tag_values
  }
}
```

### Filter desktops with any of specified tags

```hcl
variable "filter_tag_key" {}
variable "filter_tag_values" {
  type = list(string)
}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  tags_any {
    key    = var.filter_tag_key
    values = var.filter_tag_values
  }
}
```

### Filter desktops without all specified tags

```hcl
variable "filter_tag_key" {}
variable "filter_tag_values" {
  type = list(string)
}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  not_tags {
    key    = var.filter_tag_key
    values = var.filter_tag_values
  }
}
```

### Filter desktops without any of specified tags

```hcl
variable "filter_tag_key" {}
variable "filter_tag_values" {
  type = list(string)
}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  not_tags_any {
    key    = var.filter_tag_key
    values = var.filter_tag_values
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the desktop tags are located.

* `without_any_tag` - (Optional, Bool) Specifies whether to query resources without any tag.

* `tags` - (Optional, List) Specifies the list of tags to filter desktops.  
  Resources must contain all specified tags.
  The [tags](#workspace_desktop_filter_tags) structure is documented below.

* `tags_any` - (Optional, List) Specifies the list of tags to filter desktops.  
  Resources must contain at least one of specified tags.
  The [tags_any](#workspace_desktop_filter_tags) structure is documented below.

* `not_tags` - (Optional, List) Specifies the list of tags to filter desktops.  
  Resources must not contain specified tags.
  The [not_tags](#workspace_desktop_filter_tags) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the list of tags to filter desktops.  
  Resources must not contain any of specified tags.
  The [not_tags_any](#workspace_desktop_filter_tags) structure is documented below.

* `matches` - (Optional, List) Specifies the list of matching rules to filter desktops.  
  The [matches](#workspace_desktop_filter_match) structure is documented below.

<a name="workspace_desktop_filter_tags"></a>
The `tag` block supports:

* `key` - (Required, String) Specifies the key of tag.

* `values` - (Required, List) Specifies the list of tag values that matched corresponding key.

<a name="workspace_desktop_filter_match"></a>
The `match` block supports:

* `key` - (Required, String) Specifies the name of desktop property.

* `value` - (Required, String) Specifies the value of desktop property.  
  When the key is resource_name, it is a fuzzy search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `desktops` - The list of desktops that match the filter parameters.  
  The [desktops](#workspace_desktop_filter_desktops_attr) structure is documented below.

<a name="workspace_desktop_filter_desktops_attr"></a>
The `desktops` block supports:

* `resource_id` - The ID of the desktop.

* `resource_name` - The name of the desktop.

* `resource_detail` - The detail of the desktop.

* `tags` - The list of tags attached to the desktop.  
  The [tags](#workspace_desktop_filter_tags_attr) structure is documented below.

<a name="workspace_desktop_filter_tags_attr"></a>
The `tags` block in desktop object supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
