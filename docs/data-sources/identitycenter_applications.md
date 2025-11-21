---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_applications"
description: |-
  Use this data source to get the Identity Center applications.
---

# huaweicloud_identitycenter_applications

Use this data source to get the Identity Center applications.

## Example Usage

```hcl
var instance_id {}

data "huaweicloud_identitycenter_applications" "test" {
  instance_id    = var.instance_id
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the Identity Center instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - The list of the applications.
  The [applications](#applications_struct) structure is documented below.

<a name="applications_struct"></a>
The `applications` block supports:

* `application_urn` - The urn of the application.

* `application_provider_urn` - The urn of the application provider.

* `assignment_config` - The assignment configuration of the application.
  The [assignment_config](#assignment_config_struct) structure is documented below.

* `created_date` - The created date of the application.

* `description` - The description of the application.

* `instance_urn` - The urn of the Identity Center instance.

* `name` - The name of the application.

* `status` - The status of the application.

* `application_account` - The account of the application.

* `portal_options` - The portal options of the application.
  The [portal_options](#portal_options_struct) structure is documented below.

<a name="assignment_config_struct"></a>
The `assignment_config` block supports:

* `assignment_required` - Whether the application requires assignment.

<a name="portal_options_struct"></a>
The `portal_options` block supports:

* `visible` - Whether the application instance is visible.

* `visibility` - Whether the application instance is visible.

* `sign_in_options` - The sign in options of the application.
  The [sign_in_options](#sign_in_options_struct) structure is documented below.

<a name="sign_in_options_struct"></a>
The `sign_in_options` block supports:

* `origin` - This determines how IAM Identity Center navigates the user to the target application.

* `application_url` - The URL that accepts authentication requests for an application.
