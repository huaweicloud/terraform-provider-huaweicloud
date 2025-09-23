---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_logstash_parsers"
description: |-
  Use this data source to query the SecMaster collector logstash parsers within HuaweiCloud.
---

# huaweicloud_secmaster_collector_logstash_parsers

Use this data source to query the SecMaster collector logstash parsers within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_collector_logstash_parsers" "example" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the collector logstash parsers.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to query collector logstash parsers.

* `query_type` - (Optional, String) Specifies the query type. The value can be **QUICK** or **GENERAL**.

* `title` - (Optional, String) Specifies the title of the parser to query.

* `description` - (Optional, String) Specifies the description of the parser to query.

* `sort_key` - (Optional, String) Specifies the field to sort the results by.

* `sort_dir` - (Optional, String) Specifies the sorting direction. The value can be **asc** or **desc**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the list of collector logstash parsers.
  The [records](#secmaster_collector_logstash_parser) structure is documented below.

<a name="secmaster_collector_logstash_parser"></a>
The `records` block supports:

* `parser_id` - The ID of the parser.

* `title` - The title of the parser.

* `description` - The description of the parser.

* `channel_refer_count` - The number of channels that reference this parser.
