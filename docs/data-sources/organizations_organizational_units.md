---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_organizational_units"
description: ""
---

# huaweicloud_organizations_organizational_units

Use this data source to get the list of child organizational units under the specified parent OU.

## Example Usage

```hcl
data "huaweicloud_organizations_organization" "org" {}

data "huaweicloud_organizations_organizational_units" "test" {
  parent_id = data.huaweicloud_organizations_organization.org.root_id
}
```

## Argument Reference

The following arguments are supported:

* `parent_id` - (Required, String) Specifies the ID of root or organizational unit.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `children` - The list of child organizational units.
  The [children](#OrganizationalUnits) structure is documented below.

<a name="OrganizationalUnits"></a>
The `children` block supports:

* `id` - The ID of the organizational unit.

* `name` - The name of the organizational unit.

* `urn` - The uniform resource name of the organizational unit.

* `created_at` - The time when the organizational unit was created.
