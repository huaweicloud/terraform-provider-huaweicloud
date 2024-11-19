---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_physical_sessions_kill"
description: |-
  Manages a DDM physical sessions kill resource within HuaweiCloud.
---

# huaweicloud_ddm_physical_sessions_kill

Manages a DDM physical sessions kill resource within HuaweiCloud.

-> **NOTE:** Deleting physical sessions kill is not supported. If you destroy a resource of physical sessions kill,
the resource is only removed from the state.

## Example Usage

```hcl
variable instance_id {}
variable process_id {}

resource "huaweicloud_ddm_physical_sessions_kill" "test" {
  instance_id = var.instance_id
  process_ids = [var.process_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.
  Changing this creates a new resource.

* `process_ids` - (Required, List, ForceNew) Specifies the list of process IDs.

## Attribute Reference

* `id` - The resource ID.
