---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_record_delete"
description: |-
  Manages a GaussDB DR record delete resource within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_record_delete

Manages a GaussDB DR record delete resource within HuaweiCloud.

-> This resource is a one-time action resource for deleting a disaster recovery task record. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_gaussdb_dr_record_delete" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to operate the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_id` - (Required, String, NonUpdatable) Specifies the disaster recovery task ID to be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the `task_id` argument.
