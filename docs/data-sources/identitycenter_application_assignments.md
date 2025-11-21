---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_assignments"
description: |-
  Use this data source to get the Identity Center application assignments.
---

# huaweicloud_identitycenter_application_assignments

Use this data source to get the Identity Center application assignments.

## Example Usage

```hcl
var instance_id {}
var principal_id {}

data "huaweicloud_identitycenter_application_assignments" "test" {
  instance_id    = var.instance_id
  principal_id   = var.principal_id
  principal_type = "USER"
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the Identity Center instance.

* `principal_id` - (Required, String) Specifies the ID of the principal.

* `principal_type` - (Required, String) Specifies the type of the principal. Value options: **USER**, **GROUP**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `application_assignments` - The list of the application assignments.
  The [application_assignments](#application_assignments_struct) structure is documented below.

<a name="application_assignments_struct"></a>
The `application_assignments` block supports:

* `application_urn` - The urn of the application.

* `principal_id` - The ID of the principal.

* `principal_type` - The type of the principal.
