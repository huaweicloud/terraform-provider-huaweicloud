---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_stream_delay"
description: |-
  Manages a Live stream delay resource within HuaweiCloud.
---

# huaweicloud_live_stream_delay

Manages a Live stream delay resource within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

resource "huaweicloud_live_stream_delay" "test" {
  domain_name = var.domain_name
  app_name    = "live"
  delay       = 2000
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the streaming domain name.

  Changing this will create a new resource.

* `delay` - (Required, Int, ForceNew) Specifies the delay time, in ms. The options are as follows:
  + `2,000` (low).
  + `4,000` (medium).
  + `6,000` (high).

* `app_name` - (Optional, String, ForceNew) Specifies the application name. Defaults to **live**.

  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

The resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_stream_delay.test <domain_name>
```
