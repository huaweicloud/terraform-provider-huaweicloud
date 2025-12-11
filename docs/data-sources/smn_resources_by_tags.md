---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_resources_by_tags"
description: |-
  Use this data source to get the list of SMN resources by tags.
---

# huaweicloud_smn_resources_by_tags

Use this data source to get the list of SMN resources by tags.

## Example Usage

```hcl
variable "resource_type"{}

data "huaweicloud_smn_resources_by_tags" "test" {
  resource_type = var.resource_type
  action        = "filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type. Value options:
  + **smn_topic**: topic
  + **smn_sms**: SMS
  + **smn_application**: mobile push

* `action` - (Required, String) Specifies the operation to be performed. Value options:
  + **filter**: indicates pagination query
  + **count**: indicates that the total number of query results meeting the search criteria will be returned.

* `tags` - (Optional, List) Specifies the tags. A maximum of 10 keys can be queried at a time, and each key can contain
  a maximum of 10 values. The structure body must be included. The tag key cannot be left blank or be an empty string.
  Each tag key must be unique, and each tag value of a tag must also be unique. Resources identified by different keys
  are in AND relationship, and values in one tag are in OR relationship.
  The [tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies any included tags. Each tag contains a maximum of 10 keys, and each key contains
  a maximum of 10 values. The structure body cannot be missing, and the key cannot be left blank or set to an empty string.
  Each tag key must be unique, and each tag value of a tag must also be unique. Resources identified by different keys are
  in OR relationship, and values in one tag are in OR relationship.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the excluded tags. Each tag contains a maximum of 10 keys, and each key contains
  a maximum of 10 values. The structure body cannot be missing, and the key cannot be left blank or set to an empty string.
  Each tag key must be unique, and each tag value of a tag must also be unique. Resources not identified by different keys
  are in AND relationship, and values in one tag are in OR relationship.
  The [not_tags](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies any excluded tags. Each tag contains a maximum of 10 keys, and each key
  contains a maximum of 10 values. The structure body cannot be missing, and the key cannot be left blank or set to an
  empty string. Each tag key must be unique, and each tag value of a tag must also be unique. Resources not identified
  by different keys are in OR relationship, and values in one tag are in OR relationship.
  The [not_tags_any](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the key-value pair to be matched. The key can only be **resource_name**. The
  value will be exactly matched.
  The [matches](#matches_struct) structure is documented below.

* `without_any_tag` - (Optional, String) Specifies no tag is contained. If this parameter is set to **true**, all resources
  without tags are queried. In this case, the `tag`, `not_tags`, `tags_any`, and `not_tags_any` fields are ignored.
  Value options: **true**, **false**.

<a name="tags_struct"></a>
The `tags`, `not_tags`, `tags_any`, `not_tags_any` and `sys_tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.
  + A key can contain a maximum of 127 Unicode characters.
  + key must be specified.

* `values` - (Required, List) Specifies the values of the tag.
  + Each tag contains a maximum of 10 values.
  + Values of the same tag must be unique.
  + It can contain a maximum of 255 Unicode characters.
  + If this parameter is left blank, any value can be used.
  + All values of a tag key are in the OR relationship.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the tag key.
  + The key must be unique, and the value is used for matching.
  + The key field is a fixed dictionary value.
  + key cannot be left blank.

* `value` - (Required, String) Specifies the tag value.
  + It can contain a maximum of 255 Unicode characters.
  + This field cannot be left blank.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the list of resources.
  The [resources](#gresources_struct) structure is documented below.

* `total_count` - Indicates the total number of resources.

<a name="gresources_struct"></a>
The `resources` block supports:

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `resource_detail` - Indicates the resource details.
  The [resource_detail](#resources_resource_detail_struct) structure is documented below.

* `tags` - Indicates the list of tags.
  The [tags](#resources_tags_struct) structure is documented below.

<a name="resources_resource_detail_struct"></a>
The `resource_detail` block supports:

* `detail_id` - Indicates the details ID.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `topic_urn` - Indicates the unique identifier of the topic.

* `display_name` - Indicates the display name.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag value.
