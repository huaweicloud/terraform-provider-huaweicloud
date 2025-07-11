---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vault_set_resource"
description: |-
  Using this resource to configure resource backup settings for a CBR vault within HuaweiCloud.
---

# huaweicloud_cbr_vault_set_resource

Using this resource to configure resource backup settings for a CBR vault within HuaweiCloud.

-> This resource is only a one-time action resource to configure resource backup settings. Deleting this resource will
not change the backup settings of the resources, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "vault_id" {}
variable "resource_ids" {
  type = list(string)
}
variable "action" {}

resource "huaweicloud_cbr_vault_set_resource" "test" {
  vault_id     = var.vault_id
  resource_ids = var.resource_ids
  action       = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted,
  the provider-level region will be used. Changing this will create new resource.

* `vault_id` - (Required, String, NonUpdatable) Specifies the ID of the CBR vault to configure resource backup settings.

* `resource_ids` - (Required, List, NonUpdatable) Specifies the list of resource IDs for which to configure backup settings.

* `action` - (Required, String, NonUpdatable) Specifies the action to configure backup settings. Valid values:
  + **suspend**: Enable backup for the resources.
  + **unsuspend**: Disable backup for the resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the resource.
