---
subcategory: "LakeFormation"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lakeformation_instance_recover"
description: |-
  Use this resource to recover a LakeFormation instance from the recycle bin within HuaweiCloud.
---

# huaweicloud_lakeformation_instance_recover

Use this resource to recover a LakeFormation instance from the recycle bin within HuaweiCloud.

-> This resource is only a one-time action resource for recovering the LakeFormation instance from the recycle bin.
   Deleting this resource will not clear the corresponding request record, but will only remove the resource
   information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_lakeformation_instance_recover" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the instance needs to be recovered is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to be recovered
  from the recycle bin.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
