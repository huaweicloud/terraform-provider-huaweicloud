---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instance_group_assign"
description: |-
  Use this resource to assign instances to a DAS instance group within HuaweiCloud.
---

# huaweicloud_das_instance_group_assign

Use this resource to assign instances to a DAS instance group within HuaweiCloud.

-> This resource is only a one-time action resource for assigning instances to a DAS instance group.
   Deleting this resource will not clear the corresponding request record, but will only remove
   the resource information from the tfstate file.

## Example Usage

```hcl
variable "group_id" {}
variable "instance_ids" {
  type = list(string)
}

resource "huaweicloud_das_instance_group_assign" "test" {
  group_id     = var.group_id
  instance_ids = var.instance_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DAS instance group is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `group_id` - (Required, String, NonUpdatable) Specifies the instance group ID.

* `instance_ids` - (Required, List, NonUpdatable) Specifies the list of instance IDs to be assigned to the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
