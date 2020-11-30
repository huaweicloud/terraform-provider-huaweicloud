---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_agency

Manages an agency resource within huawei cloud.
This is an alternative to `huaweicloud_iam_agency_v3`

## Example Usage

```hcl
resource "huaweicloud_identity_agency" "agency" {
  name                  = "test_agency"
  description           = "test agency"
  delegated_domain_name = "***"

  project_role {
    project = "cn-north-1"
    roles = [
      "Tenant Administrator",
    ]
  }
  domain_roles = [
    "Anti-DDoS Administrator",
  ]
}
```

**Note**: It can not set `tenant_name` in `provider "huaweicloud"` when
   using this resource.

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, String) The name of agency. The name is a string of 1 to 64
    characters.

* `description` - (Optional, String) Provides supplementary information about the
    agency. The value is a string of 0 to 255 characters.

* `delegated_domain_name` - (Required, String) The name of delegated domain.

* `project_role` - (Optional, List) An array of roles and projects which are used to
    grant permissions to agency on project. The structure is documented below.

* `domain_roles` - (optional, List) An array of role names which stand for the
    permissionis to be granted to agency on domain.

The `project_role` block supports:

* `project` - (Required, String) The name of project

* `roles` - (Required, List) An array of role names

**note**:
    one or both of `project_role` and `domain_roles` must be input when
creating an agency.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The agency ID.
* `duration` - Validity period of an agency. The default value is null,
    indicating that the agency is permanently valid.
* `expire_time` - The expiration time of agency.
* `create_time` - The time when the agency was created.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 5 minute.