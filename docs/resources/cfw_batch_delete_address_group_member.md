---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_batch_delete_address_group_member"
description: |-
  Manages a resource to batch delete address group members within HuaweiCloud.
---

# huaweicloud_cfw_batch_delete_address_group_member

Manages a resource to batch delete address group members within HuaweiCloud.

-> This resource is a one-time action resource used to batch delete address group members. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "set_id" {}
variable "address_item_ids" {
  type = list(string)
}

resource "huaweicloud_cfw_batch_delete_address_group_member" "test" {
  set_id           = var.set_id
  address_item_ids = var.address_item_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `set_id` - (Required, String, NonUpdatable) Specifies the address group ID.

* `address_item_ids` - (Required, List, NonUpdatable) Specifies the IDs of address group members to be deleted.

* `fw_instance_id` - (Optional, String, NonUpdatable) Specifies the firewall ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate on assets under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `set_id`.
