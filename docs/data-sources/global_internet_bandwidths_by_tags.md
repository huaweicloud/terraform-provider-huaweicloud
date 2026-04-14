---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_internet_bandwidths_by_tags"
description: |-
  Use this data source to get a list of global internet bandwidth resources by tags.
---

# huaweicloud_global_internet_bandwidths_by_tags

Use this data source to get a list of global internet bandwidth resources filtered by tags.

## Example Usage

### Query all global internet bandwidth resources

```hcl 
data "huaweicloud_global_internet_bandwidths_by_tags" "all" { }
```

### Query global internet bandwidth resources with specific tags

```hcl 
data "huaweicloud_global_internet_bandwidths_by_tags" "filtered" { 
  tags { 
    key = "environment" 
    value = "production" 
  }
  tags {
  key = "project" 
  value = "my-project" 
  } 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tags` - (Optional, List) Specifies the list of tags to filter resources.
  The data source will return resources that match **all** the specified tags (AND logic).

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `value` - (Optional, String) Specifies the tag value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `request_id` - The request ID returned by the API.

* `resources` - The list of global internet bandwidth resources matching the filter criteria.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The ID of the global internet bandwidth resource.

* `resource_name` - The name of the global internet bandwidth resource.

* `resource_detail` - Detailed information about the resource in JSON string format.

* `tags` - The list of tags associated with the resource.

  The [resource_tags](#resource_tags_struct) structure is documented below.

<a name="resource_tags_struct"></a>
The `tags` block (within `resources`) supports:

* `key` - The tag key.

* `value` - The tag value.
