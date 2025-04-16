---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_resource_instances"
description: |-
  Use this data source to get the list of resources by tag.
---

# huaweicloud_tms_resource_instances

Use this data source to get the list of resources by tag.

## Example Usage

```hcl
data "huaweicloud_tms_resource_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `resource_types` - (Required, List) Specifies the list of resource types. It is case-sensitive.

* `tags` - (Required, List) Specifies the tags.
  The [tags](#tags_struct) structure is documented below.

* `project_id` - (Optional, String) Specifies the project ID. It is mandatory when `resource_types` contains region-specific
  service.

* `without_any_tag` - (Optional, String) Specifies whether query untagged resources.
  + **true**: only untagged resources are queried.
  + **false**: only tagged resources are queried.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `values` - (Required, List) Specifies the tag values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the list of resources.
  The [resources](#resources_struct) structure is documented below.

* `errors` - Indicates the list of errors.
  The [errors](#errors_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `project_id` - Indicates the project ID.

* `project_name` - Indicates the project name.

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `resource_type` - Indicates the resource type.

* `tags` - Indicates the resource tags.
  The [tags](#resources_tags_struct) structure is documented below.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `value` - Indicates the tag value.

<a name="errors_struct"></a>
The `errors` block supports:

* `error_code` - Indicates the error code.

* `error_msg` - Indicates the error message.

* `project_id` - Indicates the project ID.

* `resource_type` - Indicates the resource type.
