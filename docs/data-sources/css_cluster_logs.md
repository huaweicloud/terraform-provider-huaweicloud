---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_logs"
description: |-
  Use this data source to search logs of a CSS cluster within HuaweiCloud.
---

# huaweicloud_css_cluster_logs

Use this data source to search logs of a CSS cluster within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "instance_name" {}

data "huaweicloud_css_cluster_logs" "test" {
  cluster_id    = var.cluster_id
  instance_name = var.instance_name
  log_type      = "instance"
}
```

### Query Logs with Filters

```hcl
variable "cluster_id" {}
variable "instance_name" {}

data "huaweicloud_css_cluster_logs" "test" {
  cluster_id    = var.cluster_id
  instance_name = var.instance_name
  log_type      = "instance"
  level         = "ERROR"
  limit         = 10000
  time_index    = "2099-12-31T23:59:59,999"
  keyword       = "keyword"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the cluster logs.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

* `instance_name` - (Required, String) Specifies the node name in the cluster.

* `log_type` - (Required, String) Specifies the queried log type.  
  The valid values are as follows:
  + **deprecation**: Queries deprecation logs.
  + **indexingSlow**: Queries slow indexing logs.
  + **searchSlow**: Queries slow query logs.
  + **instance**: Queries run logs.

* `level` - (Optional, String) Specifies the queried log level. Defaults to **ALL**.
  The valid values are as follows:
  + **DEBUG**: Queries DEBUG logs.
  + **INFO**: Queries INFO logs.
  + **WARN**: Queries WARN logs.
  + **ERROR**: Queries ERROR logs.
  + **ALL**: Queries all log levels, including DEBUG, INFO, WARN, ERROR, and TRACE.

  ->**Note:** When `log_type` is set to **deprecation**, only **ALL** is supported. Separate multiple log levels
    using vertical bars (|), for example, `WARN|ERROR`. Do not use this method with **ALL**.

* `limit` - (Optional, Int) Specifies the maximum number of log records to be queried.
  The value ranges from **1** to **10000**. Defaults to **100**.

* `time_index` - (Optional, String) Specifies the time threshold to filter logs generated before this time.
  The value is in **yyyy-MM-ddTHH:mm:ss,SSS** format (UTC time zone), e.g. **2001-01-01T00:00:00,000**.

* `keyword` - (Optional, String) Specifies the keyword used to filter the log content.

-> **Note:** The arguments `log_type` and `level` will not affect the result when the cluster type is **Logstash**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `completed` - Whether all log files have been queried. Values are as follows:
  + **true**: All log files have been queried. There are no more results.
  + **false**: Only some of the log files have been queried. The query result is returned because the number of
    requested records has been reached or the log size has reached 1 MB.
  
* `type` - The log type. Returns when the cluster type is **Elasticsearch** or **OpenSearch**.

* `log_list` - The list of log records. Returns when the cluster type is **Elasticsearch** or **OpenSearch**.
  The [log_list](#log_list_struct) structure is documented below.

<a name="log_list_struct"></a>
  The `log_list` block supports:

* `content` - The log content.

* `date` - The log time. The value is in **yyyy-MM-ddTHH:mm:ss,SSS** format (UTC time zone).

* `level` - The log level.

* `instance_log` - The log content merged from filtered log records. Returns when the cluster type is **Logstash**.
