---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_search_criteria"
description: |-
  Use this data source to get the list of LTS search criteria.
---

# huaweicloud_lts_search_criteria

Use this data source to get the list of LTS search criteria.

## Example Usage

```hcl
variable log_group_id {}

data "huaweicloud_lts_search_criteria" "test" {
  log_group_id = var.log_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `log_group_id` - (Required, String) Specifies the ID of the log group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `search_criteria` - All search criteria that match the filter parameters.

  The [search_criteria](#search_criteria_struct) structure is documented below.

<a name="search_criteria_struct"></a>
The `search_criteria` block supports:

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the log stream.

* `criteria` - The list of the search criteria under specified log stream.

  The [criteria](#search_criteria_criteria_struct) structure is documented below.

<a name="search_criteria_criteria_struct"></a>
The `criteria` block supports:

* `id` - The ID of the search criterion.

* `name` - The name of the search criterion.

* `criteria` - The content of search criterion.

* `type` - The name of the search criterion.
