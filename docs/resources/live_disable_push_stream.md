---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_disable_push_stream"
description: |-
  Manages a disable push stream resource within HuaweiCloud.
---

# huaweicloud_live_disable_push_stream

Manages a disable push stream resource within HuaweiCloud.

-> Creating the resource indicates disable a push stream, deleting the resource indicates resume a push stream.

## Example Usage

```hcl
variable "domain_name" {}
variable "app_name" {}
variable "stream_name" {}

resource "huaweicloud_live_disable_push_stream" "test" {
  domain_name = var.domain_name
  app_name    = var.app_name
  stream_name = var.stream_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name of the disabling push stream.
  Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the application name of the disabling push stream.
  Changing this parameter will create a new resource.

* `stream_name` - (Required, String, ForceNew) Specifies the stream name of the disabling push stream.
  The stream name is not allowed to be `*`.
  Changing this parameter will create a new resource.

* `resume_time` - (Optional, String) Specifies the time of resuming push stream.
  The time is in UTC, the format is **yyyy-mm-ddThh:mm:ssZ**. e.g. **2024-06-01T15:03:01Z**
  If this parameter is not specified, the default value is `7` days. The maximum value is `90` days.
  The `resume_time` cannot be earlier than the current time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, UUID format.

## Import

The resource can be imported using `domain_name`, `app_name` and `stream_name`, separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_live_disable_push_stream.test <domain_name>/<app_name>/<stream_name>
```
