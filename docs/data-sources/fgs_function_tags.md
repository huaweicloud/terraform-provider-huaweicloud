---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_tags"
description: |-
  Use this data source to query function tags within HuaweiCloud.
---

# huaweicloud_fgs_function_tags

Use this data source to query function tags within HuaweiCloud.

## Example Usage

### Query function tags by function ID

```hcl
variable "function_id" {}

data "huaweicloud_fgs_function_tags" "test" {
  function_id = var.function_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the function is located and tags to be queried.  
  If omitted, the provider-level region will be used.

* `function_id` - (Required, String) Specifies the ID of the function to which the tags belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of the function tags.  
  The [tags](#fgs_function_tags_attr) structure is documented below.

<a name="fgs_function_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
