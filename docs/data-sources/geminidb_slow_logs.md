---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_slow_logs"
description: |-
  Use this data source to query the list of slow logs for GeminiDB Cassandra instances within HuaweiCloud.
---

# huaweicloud_geminidb_slow_logs

Use this data source to query the list of slow logs for GeminiDB Cassandra instances within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_slow_logs" "test" {
  instance_id = var.instance_id
  start_date  = "2018-08-06T10:41:14+0800"
  end_date    = "2018-08-07T10:41:14+0800"
}
```

### Filter by Node ID

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_gemini_db_slow_logs" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_date  = "2018-08-06T10:41:14+0800"
  end_date    = "2018-08-07T10:41:14+0800"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the slow logs.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the GeminiDB Cassandra instance ID.

* `start_date` - (Required, String) Specifies the start time, format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_date` - (Required, String) Specifies the end time, format is **yyyy-mm-ddThh:mm:ssZ**.

* `node_id` - (Optional, String) Specifies the node ID. If empty, it means querying all nodes under the instance.

* `type` - (Optional, String) Specifies the statement type. If empty, it means querying all statement types.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_log_list` - The list of slow logs.
  The [slow_log_list](#slow_logs_struct) structure is documented below.

<a name="slow_logs_struct"></a>
The `slow_log_list` block supports:

* `time` - The execution time.

* `database` - The database name.

* `query_sample` - The execution syntax.

* `type` - The statement type.

* `start_time` - The occurrence time, UTC time.
