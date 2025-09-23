---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_hls_configuration"
description: |-
  Manages a Live HLS configuration resource within HuaweiCloud.
---

# huaweicloud_live_hls_configuration

Manages a Live HLS configuration resource within HuaweiCloud.

-> This resource is an operational resource, and destroying it will not change the current HLS configuration.

## Example Usage

```hcl
variable "domain_name" {}
variable "app_name" {}

resource "huaweicloud_live_hls_configuration" "test" {
  domain_name = var.domain_name

  application {
    name          = var.app_name
    hls_fragment  = 5
    hls_ts_count  = 5
    hls_min_frags = 5
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name.
  Changing this parameter will create a new resource.

* `application` - (Required, List) Specifies app configuration of the ingest domain.

  The [application](#hls_application) structure is documented below.

<a name="hls_application"></a>
The `application` block supports:

* `name` - (Required, String, ForceNew) Specifies the application name.
  Changing this parameter will create a new resource.

* `hls_fragment` - (Required, Int) Specifies the HLS slice duration in seconds.

* `hls_ts_count` - (Required, Int) Specifies the number of ts slices in each M3U8 file.

* `hls_min_frags` - (Required, Int) Specifies the minimum number of ts shards in each M3U8 file.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.

## Import

The resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_hls_configuration.test <domain_name>
```
