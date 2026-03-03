---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_sync_event"
description: |-
  Manages an OMS synchronization event resource within HuaweiCloud.
---

# huaweicloud_oms_sync_event

Manages an OMS synchronization event resource within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "sync_task_id" {}
variable "object_keys" {
  type = list(string)
}

resource "huaweicloud_oms_sync_event" "test" {
  sync_task_id = var.sync_task_id
  object_keys  = var.object_keys
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `sync_task_id` - (Required, String, NonUpdatable) Specifies the synchronization task ID.

* `object_keys` - (Required, List, NonUpdatable) Specifies the list of URL-encoded names of objects to be synchronized.
  A maximum of `10` objects can be included.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
