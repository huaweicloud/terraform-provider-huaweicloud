---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iam_agency_v3"
sidebar_current: "docs-huaweicloud-resource-iam-agency-v3"
description: |-
  Manages an agency resource within huawei cloud.
---

# huaweicloud\_iam\_agency\_v3

Manages an agency resource within huawei cloud.

## Example Usage

```hcl
resource "huaweicloud_iam_agency_v3" "agency" {
  name = "test_agency"
  description = "test agency"
  delegated_domain_name = "***"
  project_role = [
    {
      project = "cn-north-1"
      roles = [
        "Tenant Administrator",
      ]
    }
  ]
  domain_roles = [
      "Anti-DDoS Administrator",
  ]
}
```

**Note**: It can not set `tenant_name` in `provider "huaweicloud"` when
   using this resource.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of agency. The name is a string of 1 to 64
    characters.

* `description` - (Optional) Provides supplementary information about the
    agency. The value is a string of 0 to 255 characters.

* `delegated_domain_name` - (Required) The name of delegated domain.

* `project_role` - (Optional) An array of roles and projects which are used to
    grant permissions to agency on project. The structure is documented below.

* `domain_roles` - (optional) An array of role names which stand for the
    permissionis to be granted to agency on domain.

The `project_role` block supports:

* `project` - (Required) The name of project

* `roles` - (Required) An array of role names

**note**:
    one or both of `project_role` and `domain_roles` must be input when
creating an agency.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `delegated_domain_name` - See Argument Reference above.
* `project_role` - See Argument Reference above.
* `domain_roles` - See Argument Reference above.
* `duration` - Validity period of an agency. The default value is null,
    indicating that the agency is permanently valid.
* `expire_time` - The expiration time of agency
* `create_time` - The time when the agency was created.
* `id` - The agency ID.
