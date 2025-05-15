---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_cold_data_eviction"
description: |-
  Use this resource to update the cold data eviction duration of the SFS turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_cold_data_eviction

Use this resource to update the cold data eviction duration of the SFS turbo within HuaweiCloud.

-> 1. The current resource is a one-time resource, and destroying this resource will not change the current status.
<br/>2. Before use this resource, please ensure the SFS turbo already bound storage backends.
<br/>3. This resource is only available for the following SFS Turbo types:
  **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,**250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**.

## Example Usage

```hcl
variable "share_id" {}
variable "action" {}
variable "gc_time" {}

resource "huaweicloud_sfs_turbo_cold_data_eviction" "test" {
  share_id = var.share_id
  action   = var.action
  gc_time  = var.gc_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, NonUpdatable) Specifies the ID of the SFS Turbo.

* `action` - (Required, String, NonUpdatable) Specifies the operation type.
  Currently, only **config_gc_time** is supported.

* `gc_time` - (Required, Int, NonUpdatable) Specifies the cold data eviction duration, in hour.
  The value ranges from `1` to `100,000,000`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
