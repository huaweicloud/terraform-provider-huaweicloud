---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instance_delete_nic"
description: |-
  Using this resource to delete a network interface card (NIC) from a protected instance in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_protected_instance_delete_nic

Using this resource to delete a network interface card (NIC) from a protected instance in SDRS within HuaweiCloud.

-> This is a one-time action resource to delete a NIC from a protected instance. Deleting this resource will
not change the current NIC configuration, but will only remove the resource information from the tfstate file.

-> Using this resource may cause unexpected changes to the ECS security group used to protect the instance.
Before using this resource, use `lifecycle` to ignore unexpected changes to the `security_group_ids` field in
resource `huaweicloud_compute_instance`.The following restrictions apply before using this resource:
<br/>1. Status of the protection group must be **available** or **protected**.
<br/>2. Status of the protected instance must be **available** or **protected**.
<br/>3. The primary NIC cannot be deleted.

## Example Usage

```hcl
variable "protected_instance_id" {}
variable "nic_id" {}

resource "huaweicloud_sdrs_protected_instance_delete_nic" "test" {
  protected_instance_id = var.protected_instance_id
  nic_id                = var.nic_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `protected_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the protected instance to delete the
  NIC from.

* `nic_id` - (Required, String, NonUpdatable) Specifies the ID of the NIC port to delete.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
