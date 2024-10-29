---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_error_logs"
description: |-
  Use this data source to get the list of DDS error logs.
---

# huaweicloud_dds_error_logs

Use this data source to get the list of DDS error logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dds_error_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.

* `instance_id` - (Required, String) Specifies the ID of the instance.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `severity` - (Optional, String) Specifies the log level. Valid values are **Warning** and **Error**.
  If it is left blank, logs of all levels can be queried.

* `node_id` - (Optional, String) Specifies the node ID.
  Nodes that can be queried:
  + Shard nodes in a cluster instance.
  + All nodes in a replica set or single node instance.

  If it is left blank, all nodes in the instance can be queried.

* `keywords` - (Optional, List) Specifies the full-text log search based on multiple keywords, indicating that all
  keywords are matched. Only fuzzy search by keyword prefix is supported. A maximum of 10 keywords are supported.
  Each keyword can contain a maximum of 512 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `error_logs` - Indicates the list of the error logs.
  The [error_logs](#attrblock--error_logs) structure is documented below.

<a name="attrblock--error_logs"></a>
The `error_logs` block supports:

* `node_id` - Indicates the node ID.

* `node_name` - Indicates the node name.

* `raw_message` - Indicates the error description.

* `severity` - Indicates the error log level.

* `log_time` - Indicates the time of the error log in the **yyyy-mm-ddThh:mm:ssZ** format.
