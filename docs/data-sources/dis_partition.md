---
subcategory: "Data Ingestion Service (DIS)"
---

# huaweicloud_dis_partitions

Use this data source to get all the partitions of a stream.

## Example Usage

### list all the partitions of a stream

```hcl
variable stream_name {}

data "huaweicloud_dis_partitions" "partition" {
  stream_name = var.stream_name
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the partitions. If omitted, the provider-level region will
  be used.

* `stream_name` - (Required, String) Specifies the name of the DIS stream.

## Attributes Reference

The following attributes are exported:

* `id` - The data source ID.

* `partitions` - The information of stream partitions. Structure is documented below.

The `partitions` block contains:

* `id` - The ID of the partition.

* `status` - The status of the partition.

* `hash_range` - Possible value range of the hash key used by each partition.

* `sequence_number_range` - Sequence number range of each partition.
