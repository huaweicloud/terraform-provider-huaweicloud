---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_certificate"
description: |-
  Manages an Identity Center application certificate resource within HuaweiCloud.
---

# huaweicloud_identitycenter_application_certificate

Manages an Identity Center application certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "application_instance_id" {}
variable "instance_id" {}

resource "huaweicloud_identitycenter_application_certificate" "test"{
  application_instance_id = var.application_instance_id
  instance_id             = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM Identity Center instance.

* `application_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the application instance.

* `status` - (Optional, String) Whether if you need to create an active certificate. Value options:
  **ACTIVE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `algorithm` - The algorithm of the application instance certificate.

* `certificate` - The certificate of the application instance.

* `certificate_id` - The ID of the application instance certificate.

* `expiry_date` - The expiry date of the application instance certificate.

* `key_size` - The key size of the application instance certificate.

* `issue_date` - The issue date of the application instance certificate.

## Import

The IdentityCenter application certificate can be imported using the `instance_id`, `application_instance_id`
and `certificate_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_identitycenter_application_certificate.test <instance_id>/<application_instance_id>/<certificate_id>
```
