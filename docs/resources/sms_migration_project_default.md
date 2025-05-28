---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_migration_project_default"
description: |-
  Manages an SMS set default migration project resource within HuaweiCloud.
---

# huaweicloud_sms_migration_project_default

Manages an SMS set default migration project resource within HuaweiCloud.

~> Deleting set default migration project resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "migration_project_id" {}

resource "huaweicloud_sms_migration_project_default" "test" {
  mig_project_id = var.migration_project_id
}
```

## Argument Reference

The following arguments are supported:

* `mig_project_id` - (Required, String, NonUpdatable) Specifies the migrate project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The migrate project ID, which equals to `mig_project_id`.
