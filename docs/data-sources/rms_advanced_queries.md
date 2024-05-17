---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_advanced_queries"
description: |-
  Use this data source to get the list of RMS advanced queries.
---

# huaweicloud_rms_advanced_queries

Use this data source to get the list of RMS advanced queries.

## Example Usage

```hcl
variable "advanced_query_name" {}

data "huaweicloud_rms_advanced_queries" "test" {
  name = var.advanced_query_name
}
```

## Argument Reference

The following arguments are supported:

* `query_id` - (Optional, String) Specifies the advanced query ID.

* `name` - (Optional, String) Specifies the advanced query name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `queries` - The list of advanced queries.

  The [queries](#queries_struct) structure is documented below.

<a name="queries_struct"></a>
The `queries` block supports:

* `name` - The advanced query name.

* `id` - The advanced query ID.

* `description` - The advanced query description.

* `expression` - The advanced query expression.

* `created_at` - The creation time of the advanced query.

* `updated_at` - The latest update time of the advanced query.
