---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_run_query"
description: |-
  Use this data source to run a Resource Query Language (ResourceQL) query within HuaweiCloud Resource Center.
---

# huaweicloud_rc_run_query

Use this data source to run a Resource Query Language (ResourceQL) query within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_run_query" "test" {
  expression = "SELECT name, type, region_id FROM resources WHERE type = 'ecs.cloudservers'"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `expression` - (Required, String) Specifies the ResourceQL expression.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `query_info` - The query information.

  The [query_info](#query_info_struct) structure is documented below.

* `results` - The query results, each element is the JSON string of the queried resource.

<a name="query_info_struct"></a>
The `query_info` block supports:

* `select_fields` - The select fields of the query.
