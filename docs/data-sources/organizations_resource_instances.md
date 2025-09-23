---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_resource_instances"
description: |-
  Use this data source to get the list of resource instances by resource type and tag.
---

# huaweicloud_organizations_resource_instances

Use this data source to get the list of resource instances by resource type and tag.

## Example Usage

```hcl
data "huaweicloud_organizations_resource_instances" "test"{
  resource_type = "organizations:accounts"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type. Value options: **organizations:policies**,
  **organizations:ous**, **organizations:accounts**, **organizations:roots**.

* `without_any_tag` - (Optional, Bool) Specifies whether only get the resources without tags.

* `tags` - (Optional, List) Specifies the list of tags to be queried. A maximum of 10 keys can be queried at a time, and
  each key can contain a maximum of 10 values. The tag key cannot be left blank or be an empty string. Each key must be
  unique, and each value for a key must be unique. Resources that contain all keys and one or multiple values listed in
  tags will be found and returned.

  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the fields to be queried.
  
  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `values` - (Required, List) Specifies the list of values of the tag.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the field name, and must be unique. Currently, only **resource_name** is supported.

* `value` - (Required, String) Specifies the field value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the list of resources.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `tags` - Indicates the list of resource tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the list of values of the tag
