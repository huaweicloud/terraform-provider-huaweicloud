---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_agency"
description: |-
  Manages an agency resource within HuaweiCloud.
---

# huaweicloud_identity_agency

Manages an agency resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

### Delegate another HUAWEI CLOUD account to perform operations on your resources

```hcl
resource "huaweicloud_identity_agency" "agency" {
  name                  = "test_agency"
  description           = "test agency"
  delegated_domain_name = "***"

  project_role {
    project = "cn-north-1"
    roles   = ["Tenant Administrator"]
  }
  domain_roles        = ["Anti-DDoS Administrator"]
  all_resources_roles = ["Server Administrator"]
  enterprise_project_roles {
    enterprise_project = "test_enterprise_project"
    roles              = ["CCE ReadOnlyAccess"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of agency. The name is a string of 1 to 64 characters.
  Changing this will create a new agency.

* `description` - (Optional, String) Specifies the supplementary information about the agency. The value is a string of
  0 to 255 characters, excluding these characters: '**@#$%^&*<>\\**'.

* `delegated_domain_name` - (Required, String) Specifies the name of delegated user domain.

* `duration` - (Optional, String) Specifies the validity period of an agency. The valid value are *FOREVER*, *ONEDAY*
  or the specific days, for example, "20". The default value is *FOREVER*.

* `project_role` - (Optional, List) Specifies an array of one or more roles and projects which are used to grant
  permissions to agency on project. The structure is documented below.

* `domain_roles` - (Optional, List) Specifies an array of one or more role names which stand for the permissions to be
  granted to agency on domain.

* `all_resources_roles` - (Optional, List) Specifies an array of one or more role names which stand for the permissions
  to be granted to agency on all resources, including those in enterprise projects, region-specific projects,
  and global services under your account.

* `enterprise_project_roles` - (Optional, List) Specifies an array of one or more roles and enterprise projects which
  are used to grant permissions to agency on project. The structure is documented below.

The `project_role` block supports:

* `project` - (Required, String) Specifies the name of project.

* `roles` - (Required, List) Specifies an array of role names.

The `enterprise_project_roles` block supports:

* `enterprise_project` - (Required, String) Specifies the name of enterprise project.

* `roles` - (Required, List) Specifies an array of role names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The agency ID.
* `expire_time` - The expiration time of agency.
* `create_time` - The time when the agency was created.

## Import

Agencies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_agency.agency 0b97661f9900f23f4fc2c00971ea4dc0
```

Note that the imported state may not be identical to your resource definition, due to `all_resources_roles` and
`enterprise_project_roles` field are missing from the API response. It is generally recommended running `terraform plan`
after importing an agency. You can then decide if changes should be applied to the agency, or the resource definition
should be updated to align with the agency. Also you can ignore changes as below.

```hcl
resource "huaweicloud_identity_agency" "agency" {
    ...

  lifecycle {
    ignore_changes = [all_resources_roles, enterprise_project_roles]
  }
}
```
