---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_groups"
description: |-
  Use this data source to get the list of CodeArts deploy application groups.
---

# huaweicloud_codearts_deploy_application_groups

Use this data source to get the list of CodeArts deploy application groups.

## Example Usage

```hcl
variable "project_id" {}

data "huaweicloud_codearts_deploy_application_groups" "test" {
  project_id = var.project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the application group list.
  The [groups](#attrblock--groups) structure is documented below.

<a name="attrblock--groups"></a>
The `groups` block supports:

* `id` - Indicates the application group ID.

* `name` - Indicates the application group name.

* `application_count` - Indicates the total number of applications in the group.

* `ordinal` - Indicates the group sorting field.

* `parent_id` - Indicates the parent application group ID.

* `path` - Indicates the group path.

* `children` - Indicates the child group name list.

* `created_by` - Indicates the ID of the group creator.

* `updated_by` - Indicates the ID of the user who last updates the group.
