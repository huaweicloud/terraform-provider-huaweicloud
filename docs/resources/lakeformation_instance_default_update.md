---
subcategory: "LakeFormation"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lakeformation_instance_default_update"
description: |-
  Use this resource to set a LakeFormation instance as default within HuaweiCloud.
---

# huaweicloud_lakeformation_instance_default_update

Use this resource to set a LakeFormation instance as default within HuaweiCloud.

-> This resource is only a one-time action resource for setting the LakeFormation instance as default. Deleting
   this resource will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_lakeformation_instance_default_update" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the instance needs to be set as default
  is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to be set as default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
