---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_disable_push_streams"
description: |-
  Use this data source to get the list of disabled push streams.
---

# huaweicloud_live_disable_push_streams

Use this datasource to get the list of disabled push streams.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_live_disable_push_streams" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Required, String) Specifies the ingest domain name of the disabling push stream.

* `app_name` - (Optional, String) Specifies the application name of the disabling push stream.

* `stream_name` - (Optional, String) Specifies the stream name of the disabling push stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `blocks` - The list of the disabled push streams.

  The [blocks](#blocks_struct) structure is documented below.

<a name="blocks_struct"></a>
The `blocks` block supports:

* `app_name` - The application name of the disabling push stream.

* `stream_name` - The stream name of the disabling push stream.

* `resume_time` - The time of the resuming push stream.
  The format is **yyyy-mm-ddThh:mm:ssZ**. e.g. **2024-09-01T15:30:20Z**.
