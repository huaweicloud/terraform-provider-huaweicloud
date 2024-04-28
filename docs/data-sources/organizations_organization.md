---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_organization"
description: ""
---

# huaweicloud_organizations_organization

Use this data source to get the Organization info and the root info.

## Example Usage

```hcl
data "huaweicloud_organizations_organization" "test"{
}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `urn` - Indicates the uniform resource name of the organization.

* `master_account_id` - Indicates the unique ID of the organization's management account.

* `master_account_name` - Indicates the name of the organization's management account.

* `created_at` - Indicates the time when the organization was created.

* `root_id` - Indicates the ID of the root.

* `root_name` - Indicates the name of the root.

* `root_urn` - Indicates the urn of the root.

* `root_tags` - Indicates the key/value attached to the root.

* `enabled_policy_types` - Indicates the list of enabled Organizations policy types in the Organization Root.
