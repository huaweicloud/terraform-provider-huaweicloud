---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_tags_filter"
description: |-
  Use this data source to filter Workspace desktops by tags within HuaweiCloud.
---

# huaweicloud_workspace_desktop_tags_filter

Use this data source to filter Workspace desktops by tags within HuaweiCloud.

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
variable tag_key = {}
variable tag_value = {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  tags {
    key = var.tag_key
    values = [ var.tag_value ]
  }
}
```

### Filter desktops with any of specified tags

```hcl
variable tag_key = {}
variable tag_value = {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  tags_any {
    key = var.tag_key
    values = [ var.tag_value ]
  }
}
```

### Filter desktops without all specified tags

```hcl
variable tag_key = {}
variable tag_value = {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  not_tags {
    key = var.tag_key
    values = [ var.tag_value ]
  }
}
```

### Filter desktops without any of specified tags

```hcl
variable tag_key = {}
variable tag_value = {}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  not_tags_any {
    key = var.tag_key
    values = [ var.tag_value ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the desktops are located.  
  If omitted, the provider-level region will be used.

* `without_any_tag` - (Optional, Bool) Specifies whether to query resources without any tag.

* `tags` - (Optional, List) Specifies the list of tags to filter desktops. Resources must
contain all specified tags.  
  The [tag object](#tag_object) structure is documented below.

* `tags_any` - (Optional, List) Specifies the list of tags to filter desktops. Resources must contain
at least one of specified tags.  
  The [tag object](#tag_object) structure is documented below.

* `not_tags` - (Optional, List) Specifies the list of tags to filter desktops. Resources must
not contain specified tags.  
  The [tag object](#tag_object) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the list of tags to filter desktops. Resources must
not contain any of specified tags.  
  The [tag object](#tag_object) structure is documented below.

* `matches` - (Optional, List) Specifies the list of matching rules to filter desktops.  
  The [match object](#match_object) structure is documented below.

<a name="tag_object"></a>
The `tag` block supports:

* `key` - (Required, String) Specifies the tag key.

* `values` - (Required, List) Specifies the list of tag values.

<a name="match_object"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the matching field.

* `value` - (Required, String) Specifies the matching value.

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
