---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_catalog_applications"
description: |
  Use this data source to get the Identity Center applications in catalog.
---

# huaweicloud_identitycenter_catalog_applications

Use this data source to get the Identity Center applications in catalog.

## Example Usage

```hcl
data "huaweicloud_identitycenter_catalog_applications" "test"{}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - The list of the applications.The [applications](#applications_struct) structure is documented below.

<a name="applications_struct"></a>
The `applications` block supports:

* `application_id` - The ID of the application.

* `display` - The display information of the application.The [display](#display_struct) structure is documented below.

* `application_type` - The type of the application.

<a name="display_struct"></a>
The `display` block supports:

* `description` - The description of the application.

* `display_name` - The display name of the application.

* `icon` - The icon of the application.
