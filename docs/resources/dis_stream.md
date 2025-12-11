---
subcategory: "Data Ingestion Service (DIS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dis_stream"
description: ""
---

# huaweicloud_dis_stream

Manages DIS Stream resource within HuaweiCloud.

## Example Usage

### Create a stream that type is BLOB

```hcl
resource "huaweicloud_dis_stream" "stream" {
  stream_name     = "terraform_test_dis_stream"
  partition_count = 1
}
```

### Create a stream that type is JSON

```hcl
resource "huaweicloud_dis_stream" "stream" {
  stream_name     = "terraform_test_dis_stream"
  partition_count = 1
  data_type       = "JSON"
  data_schema     = "{\"type\":\"record\",\"name\":\"RecordName\",\"fields\":[{\"name\":\"id\",\"type\":\"string\",\"doc\":\"Type inferred from '\\\"2017/10/11 11:11:11\\\"'\"},{\"name\":\"info\",\"type\":{\"type\":\"array\",\"items\":{\"type\":\"record\",\"name\":\"info\",\"fields\":[{\"name\":\"date\",\"type\":\"string\",\"doc\":\"Type inferred from '\\\"2018/10/11 11:11:11\\\"'\"}]}},\"doc\":\"Type inferred from '[{\\\"date\\\":\\\"2018/10/11 11:11:11\\\"}]'\"}]}"
}
```

## Argument Reference

The following arguments are supported:

* `stream_name` - (Required, String, ForceNew) Name of the DIS stream to be created.
  Changing this parameter will create a new resource.

* `partition_count` - (Required, Int) Number of the expect partitions. NOTE: Each stream can be scaled up and down a
  total of five times within one hour. After the stream is successfully scaled up or down, it cannot be scaled up or
  down again within the next one hour.

* `stream_type` - (Optional, String, ForceNew) Stream Type. The value is COMMON(means 1M bandwidth) or ADVANCED(means 5M
  bandwidth). Changing this parameter will create a new resource.

* `region` - (Optional, String, ForceNew) The region in which to create the DIS stream resource. If omitted, the
  provider-level region will be used. Changing this creates a new DIS Stream resource.

* `retention_period` - (Optional, Int, ForceNew) The number of hours for which data from the stream will be retained in
  DIS. Value range: `24` to `72`. Unit: **hour**. Default:`24`. Changing this parameter will create a new resource.

* `data_type` - (Optional, String, ForceNew) Data type of the data putting into the stream. The value is one of
  **BLOB**, **JSON** and **CSV**. Changing this parameter will create a new resource.

* `auto_scale_max_partition_count` - (Optional, Int, ForceNew) Maximum number of partition for automatic scaling.
  Changing this parameter will create a new resource.

* `auto_scale_min_partition_count` - (Optional, Int, ForceNew) Minimum number of partition for automatic scaling.
  Changing this parameter will create a new resource.

* `data_schema` - (Optional, String, ForceNew) User's JSON, CSV format data schema, described with Avro schema. Changing
  this parameter will create a new resource.

* `compression_format` - (Optional, String, ForceNew) Data compression type. The value is one of snappy, gzip and zip.
  Changing this parameter will create a new resource.

* `csv_delimiter` - (Optional, String, ForceNew) Field separator for CSV file. Changing this parameter will create a new
  resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the dis stream, Value 0
  indicates the default enterprise project. Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a resource ID in UUID format.

* `created` - Timestamp at which the DIS stream was created.

* `readable_partition_count` - Total number of readable partitions (including partitions in ACTIVE state only).

* `writable_partition_count` - Total number of writable partitions (including partitions in ACTIVE and DELETED states).

* `status` - Status of stream: **CREATING**,**RUNNING**,**TERMINATING**,**TERMINATED**,**FROZEN**.

* `stream_id` - Indicates a stream ID in UUID format.

* `partitions` - The information of stream partitions. Structure is documented below.

The `partitions` block contains:

* `id` - The ID of the partition.

* `status` - The status of the partition.

* `hash_range` - Possible value range of the hash key used by each partition.

* `sequence_number_range` - Sequence number range of each partition.

## Import

Dis stream can be imported by `stream_name`. For example,

```bash
terraform import huaweicloud_dis_stream.example _abc123
```
