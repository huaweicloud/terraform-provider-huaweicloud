---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_resource_tags"
description: |-
  Use this data source to get the list of tags by resource type.
---

# huaweicloud_tms_resource_tags

Use this data source to get the list of tags by resource type.

## Example Usage

```hcl
data "huaweicloud_tms_resource_tags" "test" {
  resource_types = "test_resource_type"
  project_id     = "test_project_id"
}
```

## Argument Reference

The following arguments are supported:

* `resource_types` - (Required, String) Specifies the resource type.

* `project_id` - (Required, String) Specifies the project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - Indicates the list of tags.
The [tags](#tags_struct) structure is documented below.

* `errors` - Indicates the list of errors.
The [errors](#errors_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `values` - Indicates the value list of the tag.

<a name="errors_struct"></a>
The `errors` block supports:

* `error_code` - Indicates the error code.

* `error_msg` - Indicates the error message.

* `project_id` - Indicates the project ID.

* `resource_type` - Indicates the resource type.
