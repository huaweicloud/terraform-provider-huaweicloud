---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_tags_filter"
description: |-
  Use this data source to filter resources by tags.
---

# huaweicloud_dns_tags_filter

Use this data source to filter resources by tags.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_dns_tags_filter" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.

* `matches` - (Optional, List) Specifies the fields to be queried.
  The [matches](#block--matches) structure is documented below.

* `tags` - (Optional, List) Specifies the list of the tags to be queried.
  The [tag_values](#block--tag_values) structure is documented below.

* `tags_any` - (Optional, List) Specifies the list of the tags to be queried.
  The [tag_values](#block--tag_values) structure is documented below.

* `not_tags` - (Optional, List) Specifies the list of the tags to be queried.
  The [tag_values](#block--tag_values) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the list of the tags to be queried.
  The [tag_values](#block--tag_values) structure is documented below.

<a name="block--matches"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key to be matched.

* `value` - (Optional, String) Specifies the value of the matching field.

<a name="block--tag_values"></a>
The `tag_values` block supports:

* `key` - (Optional, String) Specifies the key of tag.

* `values` - (Optional, List) Specifies the list of values of the tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates all dedicated resources that match the filter parameters.
  The [resources](#attrblock--resources) structure is documented below.

<a name="attrblock--resources"></a>
The `resources` block supports:

* `resource_id` - Indicates the ID of the resource.

* `resource_name` - Indicates the name of the resource.

* `tags` - Indicates the tag list associated with the resource.
  The [tags](#attrblock--resources--tags) structure is documented below.

<a name="attrblock--resources--tags"></a>
The `tags` block supports:

* `key` - Indicates the key of the resource tag.

* `value` - Indicates the value of the resource tag.
