---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_log_cache"
description: |-
  Manages a GaussDB DR log cache resource within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_log_cache

Manages a GaussDB DR log cache resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_dr_log_cache" "test" {
  instance_id   = var.instance_id
  disaster_type = "stream"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `disaster_type` - (Required, String, NonUpdatable) Specifies the disaster recovery type.
  The valid values are as follows:
    + **stream**: Stream disaster recovery.

* `xlog_keep_ratio` - (Optional, Int, NonUpdatable) Specifies the ratio of log retention space to remaining usable disk
  capacity. The value ranges from **1** to **99**. The default value is **0** (the default value of the primary instance).

  -> This parameter does not take effect for disaster recovery instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `instance_id`.
