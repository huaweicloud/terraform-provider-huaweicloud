---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_resource_tags"
description: |-
  Use this data source to query all resource tags within HuaweiCloud.
---

# huaweicloud_fgs_resource_tags

Use this data source to query all resource tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_resource_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resource tags are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of resource tags.

  The [tags](#fgs_resource_tags_attr) structure is documented below.

<a name="fgs_resource_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The values of the tag.
