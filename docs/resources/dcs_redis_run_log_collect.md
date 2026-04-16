---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_redis_run_log_collect"
description: |-
  Manages a DCS redis run log collect resource within HuaweiCloud.
---

# huaweicloud_dcs_redis_run_log_collect

Manages a DCS redis run log collect resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_redis_run_log_collect" "test" {
  instance_id = var.instance_id
  log_type    = "run"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the redis run log collect.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `log_type` - (Required, String, NonUpdatable) Specifies the type of log to collect.
  Currently, only **run** (Redis running log) is supported.

* `query_time` - (Optional, Int, NonUpdatable) Specifies the date offset, indicating to query logs from the past n
  days. Valid values are **0**, **1**, **3**, and **7**.  
  + **0**: Query today's logs (default).
  + **1**: Query logs from the past 1 day.
  + **3**: Query logs from the past 3 days.
  + **7**: Query logs from the past 7 days.

* `replication_id` - (Optional, String, NonUpdatable) Specifies the replica ID. This parameter is required when the
  instance is not a single-node instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.
