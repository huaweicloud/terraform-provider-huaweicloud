---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_advanced_query"
description: |-
  Use this data source to do an RMS advanced query.
---

# huaweicloud_rms_advanced_query

Use this data source to do an RMS advanced query.

## Example Usage

```hcl
variable "exression" {}

data "huaweicloud_rms_advanced_query" "test" {
  exression = var.exression
}
```

## Argument Reference

The following arguments are supported:

* `expression` - (Required, String) Specifies the expression of the query.

  For example, **select name, id from tracked_resources where provider = 'ecs' and type = 'cloudservers'**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The list of query results.

* `query_info` - The query info.

  The [query_info](#query_info) structure is documented below.

<a name="query_info"></a>
The `query_info` block supports:

* `select_fields` - The list of select fields.
