---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_instance_profiles"
description: |-
  Use this data source to get the Identity Center application instance profiles.
---

# huaweicloud_identitycenter_application_instance_profiles

Use this data source to get the Identity Center application instance profiles.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_instance_id" {}

data "huaweicloud_identitycenter_application_instance_profiles" "test" {
  instance_id             = var.instance_id
  application_instance_id = var.application_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the ID of the Identity Center instance.

* `application_instance_id` - (Required, String) Specifies the ID of the Identity Center application instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The name of the profile.

* `status` - The status of the profile.

* `profile_id` - The ID of the profile.
