---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_projects"
description: ""
---

# huaweicloud_identity_projects

Use this data source to query the IAM project list within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

### Obtain project information by name

```hcl
data "huaweicloud_identity_projects" "test" {
  name = "cn-north-4_demo"
}
```

### Obtain special project information by name

```hcl
data "huaweicloud_identity_projects" "test" {
  name = "MOS" // The project for OBS Billing
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the IAM project name to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `projects` - The details of the query projects. The structure is documented below.

The `projects` block supports:

* `id` - The IAM project ID.

* `name` - The IAM project name.

* `enabled` - Whether the IAM project is enabled.
