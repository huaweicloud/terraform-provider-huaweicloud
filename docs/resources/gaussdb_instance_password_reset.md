---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_password_reset"
description: |-
  Use this resource to reset the password for a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_password_reset

Use this resource to reset the password for a GaussDB instance within HuaweiCloud.

-> This resource is a one-time action resource for resetting the GaussDB instance password. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_instance_password_reset" "test" {
  instance_id = var.instance_id
  password    = "Test@12345678"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `password` - (Required, String, NonUpdatable) Specifies the database root user password.
  The password must meet the following requirements:
  + 8~32 characters.
  + Must contain at least three of the following: uppercase letters, lowercase letters, digits, and special characters
    `~!@#%^*-_-=+?,`.

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.
