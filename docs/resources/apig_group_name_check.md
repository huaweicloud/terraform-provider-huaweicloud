---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_group_name_check"
description: |-
  Use this resource to check whether the API group name already exists within HuaweiCloud.
---

# huaweicloud_apig_group_name_check

Use this resource to check whether the API group name already exists within HuaweiCloud.

-> This resource is only a one-time resource for checking the API group name. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}

resource "huaweicloud_apig_group_name_check" "test" {
  instance_id = var.instance_id
  group_name  = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the API group belongs.

* `group_name` - (Required, String) Specifies the name of the API group to be checked.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
