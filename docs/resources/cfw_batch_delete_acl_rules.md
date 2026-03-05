---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_batch_delete_acl_rules"
description: |-
  Manages a resource to batch delete ACL rules within HuaweiCloud.
---

# huaweicloud_cfw_batch_delete_acl_rules

Manages a resource to batch delete ACL rules within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch delete ACL rules. Deleting this resource will not clear
  the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. After the successful execution of the resource, it is necessary to pay attention to the value of the `data`
  attribute. If the value of `data` is empty, it means that the current operation has not taken effect.

## Example Usage

```hcl
variable "object_id" {}
variable "rule_ids" {
  type = list(string)
}

resource "huaweicloud_cfw_batch_delete_acl_rules" "test" {
  object_id = var.object_id
  rule_ids  = var.rule_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID.

* `rule_ids` - (Required, List, NonUpdatable) Specifies the IDs of the ACL rules to be deleted.

* `fw_instance_id` - (Optional, String, NonUpdatable) Specifies the firewall ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.

* `data` - The ACL rules for batch deletion.

  The [data](#Batch_Delete_Acl_Rules_Data) structure is documented below.

<a name="Batch_Delete_Acl_Rules_Data"></a>
The `data` block supports:

* `response_data` - The rules for batch deletion.

  The [response_data](#Batch_Delete_Acl_Rules_Response_Data) structure is documented below.

<a name="Batch_Delete_Acl_Rules_Response_Data"></a>
The `response_data` block supports:

* `name` - The name of ACL rule for batch deletion.

* `id` - The ID of ACL rule for batch deletion.
