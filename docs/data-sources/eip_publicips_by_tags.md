---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eip_publicips_by_tags"
description: |-
  Use this data source to get the list of EIP public IPs filtered by tags.
---

# huaweicloud_eip_publicips_by_tags

Use this data source to get the list of EIP public IPs filtered by tags.

## Example Usage

```hcl
variable "action" {}

data "huaweicloud_eip_publicips_by_tags" "test" {
  action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the operation type. The valid values are **filter** and **count**.
  + If `action` is **filter**, resources are queried with pagination.
  + If `action` is **count**, only the total number of resources matching the conditions is returned.

* `tags` - (Optional, List) Specifies the tag filter conditions. It can contain up to `10` keys.
  + The maximum number of values under each key is `10`, and the structure cannot be missing. The key cannot be empty
    or an empty string.
  + Keys cannot be duplicated, and values within the same key cannot be duplicated.

  The [tags](#query_tags) structure is documented below.

* `matches` - (Optional, List) Specifies the search fields.
  + The `key` is the field to be matched, currently only **resource_name** is supported.
  + The `value` is the matching value. This field is a fixed dictionary value.

  The [matches](#query_matches) structure is documented below.

<a name="query_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `values` - (Required, List) Specifies the list of tag values.

<a name="query_matches"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the field to match. The valid value is **resource_name**.

* `value` - (Required, String) Specifies the match value. Currently limited to EIP addresses only.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of resources.

  The [resources](#resources_struct) structure is documented below.

* `total_count` - The total number of resources.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_detail` - The resource detail object.

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `tags` - The tag list associated with the resource.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
