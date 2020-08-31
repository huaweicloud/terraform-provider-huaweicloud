---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_stream"
sidebar_current: "docs-huaweicloud-resource-lts-stream"
description: |-
  log stream management
---

# huaweicloud\_lts\_stream

Manage a log stream resource within HuaweiCloud.

## Example Usage

### create a log stream

```hcl
resource "huaweicloud_lts_group" "test_group" {
	group_name  = "test_group"
	ttl_in_days = 1
}
resource "huaweicloud_lts_stream" "test_stream" {
  group_id = huaweicloud_lts_group.test_group.id
  stream_name = "testacc_stream"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required)
  Specifies the ID of a created log group.
  Changing this parameter will create a new resource.

* `stream_name` - (Required)
  Specifies the log stream name.
  Changing this parameter will create a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The log stream ID.

* `group_id` - See Argument Reference above.

* `stream_name` - See Argument Reference above.

* `filter_count` - Number of log stream filters.

## Import

Log stream can be imported using the lts group ID and stream ID separated by a slash, e.g.

```
$ terraform import huaweicloud_lts_stream.stream_1 393f2bfd-2244-11ea-adb7-286ed488c87f/72855918-20b1-11ea-80e0-286ed488c880
```
