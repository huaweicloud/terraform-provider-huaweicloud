---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_search_logs"
description: |-
  Use this data source to search SecMaster logs within HuaweiCloud.
---

# huaweicloud_secmaster_search_logs

Use this data source to search SecMaster logs within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_id" {}
variable "pipe_id" {}

data "huaweicloud_secmaster_search_logs" "test" {
  workspace_id = var.workspace_id
  dataspace_id = var.dataspace_id
  pipe_id      = var.pipe_id
  query        = "*"
  from         = 1781506114075
  to           = 1781507014075
  sort         = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `dataspace_id` - (Required, String) Specifies the dataspace ID.

* `pipe_id` - (Required, String) Specifies the data pipe ID.

* `query` - (Required, String) Specifies the query statement.

* `from` - (Required, Int) Specifies the start time of the query (in milliseconds).

* `to` - (Required, Int) Specifies the end time of the query (in milliseconds).

* `sort` - (Optional, String) Specifies the sort order of the results by time.
  The valid values are **asc** (ascending) and **desc** (descending). Defaults to **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `analysis_results` - The analysis results.

  The [analysis_results](#analysis_results_struct) structure is documented below.

* `results` - The list of search results.

  The [results](#results_struct) structure is documented below.

<a name="analysis_results_struct"></a>
The `analysis_results` block supports:

* `size` - The number of analysis results returned.

* `total` - The total number of analysis results.

* `schema` - The list of analysis field definitions.

  The [schema](#analysis_field_struct) structure is documented below.

* `datarows` - The analysis result data.

<a name="analysis_field_struct"></a>
The `schema` block supports:

* `name` - The field name.

* `type` - The field type. The valid values are **boolean**, **byte**, **short**, **integer**, **long**, **float**,
  **half_float**, **scaled_float**, **double**, **keyword**, **text**, **date**, **ip**, **binary**, **object**, **nested**.

* `alias` - The field alias.

<a name="results_struct"></a>
The `results` block supports:

* `data_source` - The original log content.

* `timestamp` - The data reception time.
