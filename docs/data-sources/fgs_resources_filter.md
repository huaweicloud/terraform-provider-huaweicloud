---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_resources_filter"
description: |-
  Use this data source to query basic information about resources within HuaweiCloud.
---

# huaweicloud_fgs_resources_filter

Use this data source to query basic information about resources within HuaweiCloud.

## Example Usage

### Filter resources by function name using filter matches

```hcl
variable "function_name" {}

data "huaweicloud_fgs_resources_filter" "test" {
  resource_type = "functions"

  matches {
    key   = "resource_name"
    value = var.function_name
  }
}
```

### Filter resources by tags

```hcl
data "huaweicloud_fgs_resources_filter" "test" {
  resource_type = "functions"

  tags {
    key    = "owner"
    values = ["Administrator", "User"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the target resources are located.  
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the type of the resource used to filter the target resources.  
  Only **functions** is supported currently.

* `matches` - (Optional, List) Specifies the key-value pairs used to filter the target resources.  
  The [matches](#fgs_resources_filter_matches) structure is documented below.

* `tags` - (Optional, List) Specifies the resource tags used to filter the target resources.  
  The [tags](#fgs_resources_filter_tags) structure is documented below.

<a name="fgs_resources_filter_matches"></a>
The `matches` block supports:

* `key` - (Required, String) The match key used to filter the target resources.  
  Only **resource_name** is supported currently.

* `value` - (Required, String) The match value used to filter the target resources.  
  The value is fuzzy match if `key` is **resource_name**.

<a name="fgs_resources_filter_tags"></a>
The `tags` block supports:

* `key` - (Required, String) The key of the resource tag used to filter the target resources.

* `values` - (Required, List) The values corresponding to the current key used to filter the target resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of target resources that matched filter parameters.  
  The [resources](#fgs_filtered_resources_attr) structure is documented below.

<a name="fgs_filtered_resources_attr"></a>
The `resources` block supports:

* `id` - The ID of the resource.

* `name` - The name of the resource.

* `detail` - The detailed information of the resource, in JSON format.

* `tags` - The tags of the resource.  
  The [tags](#fgs_filtered_resource_tags) structure is documented below.

<a name="fgs_filtered_resource_tags"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
