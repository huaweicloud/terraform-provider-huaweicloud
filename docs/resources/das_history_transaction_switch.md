---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_history_transaction_switch"
description: |-
  Use this resource to enable or disable the history transaction switch within HuaweiCloud.
---

# huaweicloud_das_history_transaction_switch

Use this resource to enable or disable the history transaction switch within HuaweiCloud.

-> This resource is a one-time action resource for switching the history transaction. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_das_history_transaction_switch" "test" {
  instance_id    = var.instance_id
  status         = "Enabled"
  datastore_type = "MySQL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the history transaction switch is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `datastore_type` - (Required, String, NonUpdatable) Specifies the database type. The valid value is **MySQL**.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the database instance.

* `status` - (Required, String, NonUpdatable) Specifies the switch status of the history transaction.  
  The valid values are as follows:
  + **Enabled**: Enable the history transaction switch.
  + **Disabled**: Disable the history transaction switch.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
