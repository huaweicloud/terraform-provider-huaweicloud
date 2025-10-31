---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_check_agency_role_assignment"
description: |-
    Use this data source to query an IAM agency role assignment within HuaweiCloud.
---

# huaweicloud_identity_check_agency_role_assignment

Use this data source to query an IAM agency role assignment within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

### Check whether the agency has global service permissions

```hcl
variable "agency_id" {}
variable "role_id" {}
variable "domain_id" {}

data "huaweicloud_identity_check_agency_role_assignment" "test" {
  agency_id = var.agency_id
  role_id   = var.role_id
  domain_id = var.domain_id
}
```

### Check whether the agency has project service permissions

```hcl
variable "agency_id" {}
variable "role_id" {}
variable "project_id" {}

data "huaweicloud_identity_check_agency_role_assignment" "test" {
  agency_id  = var.agency_id
  role_id    = var.role_id
  project_id = var.project_id
}
```

### Check whether the agency has the specified permissions for all projects

```hcl
variable "agency_id" {}
variable "role_id" {}

data "huaweicloud_identity_check_agency_role_assignment" "test" {
  agency_id  = var.agency_id
  role_id    = var.role_id
  project_id = "all"
}
```

## Argument Reference

* `agency_id` - (Required, String) Specifies the agency id.

* `role_id` - (Required, String) Specifies the role id.

* `domain_id` - (Optional, String) Specifies the domain id.

* `project_id` - (Optional, String) Specifies the project id.
  If `project_id` is set to **all**, it means to check whether the agency has the specified permissions for all
  projects, including existing and future projects.

Exactly one of `domain_id` or `project_id` must be specified.

## Attribute Reference

* `result` - Indicates whether the agency has the role assignment in the scope.
