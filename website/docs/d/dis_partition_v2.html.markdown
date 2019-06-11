---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dis_partition_v2"
sidebar_current: "docs-huaweicloud-datasource-dis-partition-v2"
description: |-
  Get all the partitions of a stream
---

# huaweicloud\_dis\_partition\_v2

Get all the partitions of a stream

## Example Usage

### list all the partitions of a stream

```hcl

data "huaweicloud_dis_partition_v2" "partition" {
  stream_name = "{{ stream_name }}"
}
```

## Argument Reference

* `stream_name` -
  (Required)
  Name of the DIS stream.

## Attributes Reference

The following attributes are exported:

* `partitions` - The information of stream partitions. Structure is documented below.

The `partitions` block contains:

* `id` -  The ID of the partition.

* `status` - The status of the partition.

* `hash_range` - Possible value range of the hash key used by each partition.

* `sequence_number_range` - Sequence number range of each partition..
