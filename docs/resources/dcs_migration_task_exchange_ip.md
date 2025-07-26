---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_migration_task_exchange_ip"
description: |-
  Manages a DCS IP addresses exchange of the source and target instances during incremental data migration resource within
  HuaweiCloud.
---

# huaweicloud_dcs_migration_task_exchange_ip

Manages a DCS IP addresses exchange of the source and target instances during incremental data migration resource within
HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_dcs_migration_task_exchange_ip" "test"{
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_id` - (Required, String, NonUpdatable) Specifies the ID of the migration task.

* `exchanged_ip` - (Optional, List, NonUpdatable) Specifies the list of IP address to be switched during data migration.

* `is_exchange_domain` - (Optional, Bool, NonUpdatable) Specifies whether to switch the domain name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `task_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
