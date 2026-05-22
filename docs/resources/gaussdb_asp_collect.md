---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_asp_collect"
description: |-
  Manages a GaussDB ASP collect resource within HuaweiCloud.
---

# huaweicloud_gaussdb_asp_collect

Manages a GaussDB ASP collect resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_asp_collect" "test" {
  instance_id = var.instance_id
  start_time  = "2024-01-01T00:00:00+0800"
  end_time    = "2024-12-31T23:59:59+0800"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ASP collect. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `start_time` - (Required, String, NonUpdatable) Specifies the start time for ASP collect. The format is
  **yyyy-mm-ddThh:mm:ssZ**, where T indicates the start of a certain time, and Z indicates the time zone offset.

* `end_time` - (Required, String, NonUpdatable) Specifies the end time for ASP collect. The format is
  **yyyy-mm-ddThh:mm:ssZ**, where T indicates the start of a certain time, and Z indicates the time zone offset.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default timeout is 30 minutes.
