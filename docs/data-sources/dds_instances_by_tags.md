---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instances_by_tags"
description: |-
  Use this data source to get the list of instances.
---

# huaweicloud_dds_instances_by_tags

Use this data source to get the list of instances.

## Example Usage

```hcl
variable "action" {}

data "huaweicloud_dds_instances_by_tags" "test" {
  action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the resource type.
  The valid values are **filter** and **count**.
  If `action` set to **filter**, indicates query the instances based on tags filtering conditions.
  If `action` set to **count**, indicates only query the total number of intances.

* `matches` - (Optional, List) Specifies the fields to be queried.
  The [matches](#query_matches) structure is documented below.

* `tags` - (Optional, List) Specifies the list of the tags to be queried.
  The [tags](#query_tags) structure is documented below.

<a name="query_matches"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key to be matched.
  The valid values are as follows:
  + **instance_name**: Indicates matching queries based on instance name.
  + **instance_id**: Indicates matching queries based on instance ID.

* `value` - (Required, String) Specifies the value of the matching field.

<a name="query_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `values` - (Required, List) Specifies the list of values of the tag.
  The `values` can be empty.
  If the values are an empty list, it indicates that any value is queried. The relationship between values is OR.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of instances.
  The [instances](#instances_struct) structure is documented below.

* `total_count` - The total number of the instances.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

* `tags` - Indicates the tag list associated with the resource.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
