---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_batch_delete_domain_sets"
description: |-
  Manages a resource to batch delete domain sets within HuaweiCloud.
---

# huaweicloud_cfw_batch_delete_domain_sets

Manages a resource to batch delete domain sets within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch delete domain sets. Deleting this resource will not clear
  the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. After the successful execution of the resource, it is necessary to pay attention to the value of the `data`
  attribute. If the value of `data` is empty, it means that the current operation has not taken effect.

## Example Usage

```hcl
variable "object_id" {}
variable "set_ids" {
  type = list(string)
}

resource "huaweicloud_cfw_batch_delete_domain_sets" "test" {
  object_id = var.object_id
  set_ids   = var.set_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID.

* `set_ids` - (Required, List, NonUpdatable) Specifies the IDs of domain sets to be deleted.

* `fw_instance_id` - (Optional, String, NonUpdatable) Specifies the firewall ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate on assets under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.

* `data` - The domain sets for batch deletion.

  The [data](#Batch_Delete_Domain_Sets_Data) structure is documented below.

<a name="Batch_Delete_Domain_Sets_Data"></a>
The `data` block supports:

* `response_data` - The domain sets for batch deletion.

  The [response_data](#Batch_Delete_Domain_Sets_Response_Data) structure is documented below.

<a name="Batch_Delete_Domain_Sets_Response_Data"></a>
The `response_data` block supports:

* `name` - The name of the domain set for batch deletion.

* `id` - The ID of the domain set for batch deletion.
