---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_delete_fault_instance"
description: |-
  Manages a CBH delete fault instance resource within HuaweiCloud.
---

# huaweicloud_cbh_delete_fault_instance

Manages a CBH delete fault instance resource within HuaweiCloud.

-> This resource is a one-time action resource to delete a faulty CBH instance. Deleting this resource
will not recover the deleted CBH instance, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_cbh_delete_fault_instance" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to delete the faulty CBH instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the faulty CBH instance to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `instance_id`).
