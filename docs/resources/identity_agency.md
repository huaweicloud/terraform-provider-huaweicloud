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
variable "agency_name" {}
variable "domain_name" {}

resource "huaweicloud_identity_agency" "test" {
  name                  = var.agency_name
  description           = "Created by terraform script"
  delegated_domain_name = var.domain_name

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

* `name` - (Required, String, NonUpdatable) Specifies the name of agency.  
  The valid length is limited from `1` to `64`.

* `delegated_domain_name` - (Required, String) Specifies the name of delegated user domain.

* `description` - (Optional, String) Specifies the description (supplementary information) of the agency.  
  The valid length is limited from `0` to `255`, and excluding these characters: `@#$%^&*<>\`.

* `duration` - (Optional, String) Specifies the validity period of the agency.  
  The valid values are as follows:
  + **FOREVER**,
  + **ONEDAY**
  + A specific days, e.g. **20**.

  Default to **FOREVER**.

* `project_role` - (Optional, List) Specifies the roles assignment for the agency which the projects are used to
  grant.  
  The [project_role](#iam_agency_project_role) structure is documented below.

* `domain_roles` - (Optional, List) Specifies the roles assignment for the agency which the domain are used to grant.

* `all_resources_roles` - (Optional, List) Specifies the roles assignment for the agency which the all resources are
  used to grant.  
  Each assignment includes enterprise projects, region-specific projects, and global services under your account.

* `enterprise_project_roles` - (Optional, List) Specifies the roles assignment for the agency which the enterprise
  projects are used to grant.  
  The [enterprise_project_roles](#iam_agency_enterprise_project_roles) structure is documented below.

<a name="iam_agency_project_role"></a>
The `project_role` block supports:

* `project` - (Required, String) Specifies the name of project.

* `roles` - (Required, List) Specifies the list of role names used for assignment in a specified project.

<a name="iam_agency_enterprise_project_roles"></a>
The `enterprise_project_roles` block supports:

* `enterprise_project` - (Required, String) Specifies the name of enterprise project.

* `roles` - (Required, List) Specifies the list of role names used for assignment in a specified enterprise project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The agency ID.

* `expire_time` - The expiration time of agency.

* `create_time` - The creation time of the agency.

## Import

Agencies can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_identity_agency.test <id>
```

Note that the imported state may not be identical to your resource definition, due to `all_resources_roles` and
`enterprise_project_roles` field are missing from the API response. It is generally recommended running `terraform plan`
after importing an agency. You can then decide if changes should be applied to the agency, or the resource definition
should be updated to align with the agency. Also you can ignore changes as below.

```hcl
resource "huaweicloud_identity_agency" "test" {
  ...

  lifecycle {
    ignore_changes = [
      all_resources_roles,
      enterprise_project_roles,
    ]
  }
}
```
