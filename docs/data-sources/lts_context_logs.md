---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_context_logs"
description: |-
  Use this data source to query context log list under the specified log stream within HuaweiCloud.
---

# huaweicloud_lts_context_logs

Use this data source to query context log list under the specified log stream within HuaweiCloud.

## Example Usage

### Query logs by the specified line number

```hcl
var "log_group_id" {}
var "log_stream_id" {}
var "line_num" {}

data "huaweicloud_lts_context_logs" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  line_num      = var.line_num
}
```

### Query logs with custom time field

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

data "huaweicloud_lts_context_logs" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  line_num      = try(data.huaweicloud_lts_logs.test.logs[0].line_num, null)
  time          = try(data.huaweicloud_lts_logs.test.logs[0].labels.__time__, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the context logs are located.  
  If omitted, the provider-level region will be used.

* `log_group_id` - (Required, String) Specifies the ID of the log group to which the logs belong.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream to which the logs belong.

* `line_num` - (Optional, String) Specifies the sequence number of a log line.

* `time` - (Optional, String) Specifies the time field of the custom time function, in millisecond timestamp.  

  -> If the structured configuration of this log stream has enabled the custom time function,
     this parameter is required.

* `backwards_size` - (Optional, Int) Specifies the number of logs before the start log.  
  The valid value ranges from `0` to `500`. Defaults to **100**.

* `forwards_size` - (Optional, Int) Specifies the number of logs after the start log.  
  The valid value ranges from `0` to `500`. Defaults to **100**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of context logs that match the filter parameters.  
  The [logs](#lts_context_logs_attr) structure is documented below.

* `backwards_count` - The number of logs queried backward based on `line_num`.

* `forwards_count` - The number of logs queried forward based on `line_num`.

* `total_count` - The total number of logs, including the starting log specified in the request parameters.

<a name="lts_context_logs_attr"></a>
The `logs` block supports:

  * `content` - The original log data.

  * `line_num` - The log line sequence number.

  * `labels` - The labels contained in this log entry.
