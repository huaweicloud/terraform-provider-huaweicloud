---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_slow_logs"
description: |-
  Use this data source to get the list of DDS instance slow logs.
---

# huaweicloud_dds_slow_logs

Use this data source to get the list of DDS instance slow logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dds_slow_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.

* `instance_id` - (Required, String) Specifies the ID of the instance.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `collection_keywords` - (Optional, List) Specifies the fuzzy search for logs based on multiple database table name
  keywords, indicating that at least one keyword is matched.

* `database_keywords` - (Optional, List) Specifies the fuzzy search for logs based on multiple database keywords,
  indicating that at least one keyword is matched.

* `keywords` - (Optional, List) Specifies the full-text log search based on multiple keywords, indicating that all
  keywords are matched.

* `max_cost_time` - (Optional, Int) Specifies the logs can be searched based on the maximum execution duration.
  Unit is ms.

* `min_cost_time` - (Optional, Int) Specifies the logs can be searched based on the minimum execution duration.
  Unit is ms.

* `node_id` - (Optional, String) Specifies the node ID.

* `operate_type` - (Optional, String) Specifies the statement type. Valid values are **insert**, **query**, **update**,
  **remove**, **getmore**, **command** and **killcursors**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_logs` - Indicates the list of the slow logs.
  The [slow_logs](#attrblock--slow_logs) structure is documented below.

<a name="attrblock--slow_logs"></a>
The `slow_logs` block supports:

* `node_id` - Indicates the node ID.

* `node_name` - Indicates the node name.

* `collection` - Indicates the name of the database table which the log belongs to.

* `cost_time` - Indicates the execution time. Unit is ms.

* `database` - Indicates the name of the database which the log belongs to.

* `docs_returned` - Indicates the number of returned documents.

* `docs_scanned` - Indicates the number of scanned documents.

* `lock_time` - Indicates the lock wait time. Unit is ms.

* `log_time` - Indicates the time of the slow log in the **yyyy-mm-ddThh:mm:ssZ** format.

* `operate_type` - Indicates the statement type.

* `whole_message` - Indicates the statement.
