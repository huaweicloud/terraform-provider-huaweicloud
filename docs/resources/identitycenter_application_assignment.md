---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_assignment"
description: |-
  Manages an Identity Center application assignment resource within HuaweiCloud.
---

# huaweicloud_identitycenter_application_assignment

Manages an Identity Center application assignment resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_instance_id" {}
variable "principal_id" {}

resource "huaweicloud_identitycenter_application_assignment" "test"{
  instance_id             = var.instance_id
  application_instance_id = var.application_instance_id
  principal_id            = var.principal_id
  principal_type          = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center instance.

* `application_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the application instance.

* `principal_id` - (Required, String, NonUpdatable) Specifies the ID of the user or user group that belongs to IAM
  Identity Center.

* `principal_type` - (Required, String, NonUpdatable) Specifies the type of the user or user group.
  Value options: **USER**, **GROUP**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `application_urn` - The urn of the application.

## Import

The Identity Center application assignment can be imported using the `instance_id`, `application_instance_id`
and `principal_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_identitycenter_application_assignment.test <instance_id>/<application_instance_id>/<principal_id>
```
