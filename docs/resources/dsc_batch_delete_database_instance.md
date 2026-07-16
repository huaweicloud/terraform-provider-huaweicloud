---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_batch_delete_database_instance"
description: |-
  Manages a resource to batch delete DSC database instances within HuaweiCloud.
---

# huaweicloud_dsc_batch_delete_database_instance

Manages a resource to batch delete DSC database instances within HuaweiCloud.

-> This resource is a one-time action resource used to batch delete DSC database instances.
Deleting this resource will not restore the deleted database instances or undo the delete action, but will only
remove the resource information from the tf state file.

## Example Usage

```hcl
variable "db_ids" {
  type = list(string)
}

resource "huaweicloud_dsc_batch_delete_database_instance" "test" {
  db_ids = var.db_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `db_ids` - (Required, List, NonUpdatable) Specifies the database instance IDs to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
