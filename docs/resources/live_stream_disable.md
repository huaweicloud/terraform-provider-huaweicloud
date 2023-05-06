---
subcategory: "Live"
---

# huaweicloud_live_stream_disable

Manages a Live stream disable within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

resource "huaweicloud_live_stream_disable" "test"{
  domain_name = var.domain_name
  app_name    = "live"
  stream_name = "test_stream"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name.

  Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the application name.

  Changing this parameter will create a new resource.

* `stream_name` - (Required, String, ForceNew) Specifies the stream name(not *).

  Changing this parameter will create a new resource.

* `resume_time` - (Optional, String) Specifies the time to resume stream push.
  The format is yyyy-mm-ddThh:mm:ssZ (UTC time). Default value is 7 days. The maximum value is 90 days.
  The value must greater than the current time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The live stream disable can be imported using the domain name, app name and stream name separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_live_stream_disable.test <domain_name>/<app_name>/<stream_name>
```
