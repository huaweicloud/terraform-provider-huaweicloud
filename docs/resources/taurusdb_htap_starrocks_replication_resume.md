---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_replication_resume"
description: |-
  Manages a TaurusDB HTAP StarRocks replication resume resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_replication_resume

Manages a TaurusDB HTAP StarRocks replication resume resource within HuaweiCloud.

-> This resource is a one-time action resource to resume a data synchronization task for a StarRocks instance.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "task_name" {}

resource "huaweicloud_taurusdb_htap_starrocks_replication_resume" "test" {
  instance_id = var.instance_id
  task_name   = var.task_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the StarRocks instance ID.

* `task_name` - (Required, String, NonUpdatable) Specifies the synchronization task name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is a UUID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
