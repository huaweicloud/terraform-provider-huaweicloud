---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateway_tags"
description: |-
  Use this data source to get the list of all NAT gateways tags under the specified project.
---

# huaweicloud_nat_gateway_tags

Use this data source to get the list of all NAT gateways tags under the specified project.

## Example Usage

```hcl
data "huaweicloud_nat_gateway_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of the tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The list of the tag values.
