---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_advanced_query"
description: |-
  Use this data source to do an advanced query in the RMS resource aggregator.
---

# huaweicloud_rms_resource_aggregator_advanced_query

Use this data source to do an advanced query in the RMS resource aggregator.

## Example Usage

```hcl
variable "aggregator_id" {}
variable "exression" {}

data "huaweicloud_rms_resource_aggregator_advanced_query" "test" {
  aggregator_id = var.aggregator_id
  exression     = var.exression
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the aggregator ID.

* `expression` - (Required, String) Specifies the expression of the query.

  For example, **select name, id from aggregator_resources where provider = 'ecs' and type = 'cloudservers'**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The list of query results.

* `query_info` - The query info.

  The [query_info](#query_info) structure is documented below.

<a name="query_info"></a>
The `query_info` block supports:

* `select_fields` - The list of select fields.
