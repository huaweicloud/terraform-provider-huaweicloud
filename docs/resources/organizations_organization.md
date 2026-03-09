---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_organization"
description: |-
  Manages an Organizations organization resource within HuaweiCloud.
---

# huaweicloud_organizations_organization

Manages an Organizations organization resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_organizations_organization" "test" {}
```

## Argument Reference

The following arguments are supported:

* `enabled_policy_types` - (Optional, List) Specifies the list of policy types to be enabled in the
  organization root.  
  The valid values are as follows:
  + **service_control_policy**: Service control policy.
  + **tag_policy**: Tag policy.

* `root_tags` - (Optional, Map) Specifies the key/value pairs to be attached to the root.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the organization.

* `master_account_id` - Indicates the unique ID of the organization's management account.

* `master_account_name` - Indicates the name of the organization's management account.

* `created_at` - Indicates the time when the organization was created.

* `root_id` - Indicates the ID of the root.

* `root_name` - Indicates the name of the root.

* `root_urn` - Indicates the urn of the root.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.

## Import

The Organizations organization can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_organization.test <id>
```
