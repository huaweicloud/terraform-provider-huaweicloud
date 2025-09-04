---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_logs"
description: |-
  Use this data source to query logs under specified log stream within HuaweiCloud.
---

# huaweicloud_lts_logs

Use this data source to query log list under the specified log stream within HuaweiCloud.

## Example Usage

### Basic log query

```hcl
var "log_group_id" {}
var "log_stream_id" {}
var "start_time" {}
var "end_time" {}

data "huaweicloud_lts_logs" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  start_time    = var.start_time
  end_time      = var.end_time
}
```

### Query with custom time enabled

```hcl
var "log_group_id" {}
var "log_stream_id" {}
var "start_time" {}
var "end_time" {}

data "huaweicloud_lts_logs" "test" {
  log_group_id           = var.log_group_id
  log_stream_id          = var.log_stream_id
  start_time             = var.start_time
  end_time               = var.end_time
  is_custom_time_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the LTS logs are located.  
  If omitted, the provider-level region will be used.

* `log_group_id` - (Required, String) Specifies the ID of the log group to which the logs belong.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream to which the logs belong.

* `start_time` - (Required, String) Specifies the start time for querying log list, in RFC3339 format.

* `end_time` - (Required, String) Specifies the end time for querying log list, in RFC3339 format.

* `labels` - (Optional, Map) Specifies the labels in key/value format to be queried.

* `keywords` - (Optional, String) Specifies the keywords for exact search.

* `is_custom_time_enabled` - (Optional, Boolean) Specifies whether to enable the custom time function
  for the log stream structured configuration.  
  Defaults to **false**.

  -> If the structured configuration of this log stream has enabled the custom time function, this parameter must
     be set to **true**.

* `highlight` - (Optional, Boolean) Specifies whether to highlight the keyword in the logs.  
  Defaults to **true**.

* `is_desc` - (Optional, Boolean) Specifies whether to sort the logs in descending order.  
  Defaults to **false**.

* `is_iterative` - (Optional, Boolean) Specifies whether to enable iterative query.  
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of logs that match the filter parameters.  
  The [logs](#lts_logs_attr) structure is documented below.

<a name="lts_logs_attr"></a>
The `logs` block supports:

  * `content` - The content of the log.

  * `line_num` - The line number of the log.

  * `labels` - The labels associated with the log.
