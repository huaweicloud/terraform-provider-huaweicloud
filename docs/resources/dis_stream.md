---
subcategory: "Data Ingestion Service (DIS)"
---

# huaweicloud\_dis\_stream

DIS Stream management
This is an alternative to `huaweicloud_dis_stream_v2`

## Example Usage

### create a stream that type is BLOB

```hcl
resource "huaweicloud_dis_stream" "stream" {
  stream_name     = "terraform_test_dis_stream"
  partition_count = 1
}
```

### create a stream that type is JSON

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

* `partition_count` - (Required) Number of the expect partitions. NOTE: Each stream can be scaled up
  and down a total of five times within one hour. After the stream is
  successfully scaled up or down, it cannot be scaled up or down again
  within the next one hour.

* `stream_name` - (Required) Name of the DIS stream to be created.

* `region` - (Optional) The region in which to create the DIS stream resource. If omitted, the provider-level region will be used. Changing this creates a new DIS Stream resource.

* `auto_scale_max_partition_count` - (Optional) Maximum number of partition for automatic scaling.  Changing this parameter will create a new resource.

* `auto_scale_min_partition_count` - (Optional) Minimum number of partition for automatic scaling.  Changing this parameter will create a new resource.

* `compression_format` - (Optional) Data compression type. The value is one of snappy, gzip and zip.  Changing this parameter will create a new resource.

* `csv_delimiter` - (Optional) Field separator for CSV file.  Changing this parameter will create a new resource.

* `data_schema` - (Optional) User's JOSN, CSV format data schema, described with Avro schema.  Changing this parameter will create a new resource.

* `data_type` - (Optional) Data type of the data putting into the stream. The value is one of
  BLOB, JSON and CSV.  Changing this parameter will create a new resource.

* `retention_period` - (Optional) The number of hours for which data from the stream will be retained
  in DIS.  Changing this parameter will create a new resource.

* `stream_type` - (Optional) Stream Type. The value is COMMON(means 1M bandwidth) or
  ADVANCED(means 5M bandwidth).  Changing this parameter will create a new resource.

* `tags` - (Optional) List of tags for the newly created DIS stream. Structure is documented below. Changing this parameter will create a new resource.

The `tags` block supports:

* `key` - (Optional) The key of tag.  Changing this parameter will create a new resource.

* `value` - (Optional) The value of tag.  Changing this parameter will create a new resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created` - Timestamp at which the DIS stream was created.

* `readable_partition_count` - Total number of readable partitions (including partitions in ACTIVE state only).

* `writable_partition_count` - Total number of writable partitions (including partitions in ACTIVE and DELETED states).
