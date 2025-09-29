---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_resources_filter"
description: |-
  Use this data source to get the list of SWR enterprise resources filtered by tags.
---

# huaweicloud_swr_enterprise_resources_filter

Use this data source to get the list of SWR enterprise resources filtered by tags.

## Example Usage

```hcl
variable "resource_type"{}

data "huaweicloud_swr_enterprise_resources_filter" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) The type of the resource used to filter the target resources.
  Valid value is **instances**.

* `tags` - (Optional, List) The resource tags used to filter the target resources.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) The key of the resource tag used to filter the target resources.

* `values` - (Required, List) The values corresponding to the current key used to filter the target resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total count of the resources.

* `resources` - The list of target resources that matched filter parameters.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The ID of the resource.

* `resource_detail` - The detailed information of the resource, in JSON format.

* `resource_name` - The name of the resource.

* `tags` - The key/value tag pairs to associate with the resource.

* `sys_tags` - The key/value system tag pairs to associate with the resource.
