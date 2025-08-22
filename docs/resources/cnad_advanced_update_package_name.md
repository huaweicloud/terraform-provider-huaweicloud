---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_update_package_name"
description: |-
  Updates the name of a CNAD advanced package within HuaweiCloud.
---

# huaweicloud_cnad_advanced_update_package_name

Updates the name of a CNAD advanced package within HuaweiCloud.

-> This resource is a one-time action resource for updating the name of a CNAD package. Deleting this resource
   will not revert the package name, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "package_id" {}
variable "package_name" {}

resource "huaweicloud_cnad_advanced_update_package_name" "test" {
  package_id = var.package_id
  name       = var.package_name
}
```

## Argument Reference

The following arguments are supported:

* `package_id` - (Required, String, NonUpdatable) Specifies the ID of the CNAD package.

* `name` - (Required, String, NonUpdatable) Specifies the new name for the CNAD package.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
