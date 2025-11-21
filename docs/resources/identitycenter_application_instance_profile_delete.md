---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_instance_profile_delete"
description: |-
  Manages an Identity Center application instance profile delete resource within HuaweiCloud.
---

# huaweicloud_identitycenter_application_instance_profile_delete

Manages an Identity Center application instance profile delete resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
  Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_instance_id" {}
variable "profile_id" {}

resource "huaweicloud_identitycenter_application_instance_profile_delete" "test" {
  instance_id             = var.instance_id
  application_instance_id = var.application_instance_id
  profile_id              = var.profile_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center instance.

* `application_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center
  application instance.

* `profile_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center application
  instance profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
