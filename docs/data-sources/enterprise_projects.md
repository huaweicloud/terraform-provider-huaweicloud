---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_projects"
description: ""
---

# huaweicloud_enterprise_projects

Use this data source to get the list of EPS enterprise projects.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_enterprise_projects" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the enterprise project name. Fuzzy search is supported.

* `enterprise_project_id` - (Optional, String) Specifies the ID of an enterprise project.
  The value **0** indicates enterprise project default.

* `status` - (Optional, Int) Specifies the status of an enterprise project. The valid values are as follows:
  + **1**: Enabled.
  + **2**: Disabled.

* `type` - (Optional, String) Specifies the type of an enterprise project. The valid values are as follows:
  + **prod**: Commercial project.
  + **poc**: Test project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `enterprise_projects` - All projects that match the filter parameters.
  The [enterprise_projects](#Enterprise_Projects) structure is documented below.

<a name="Enterprise_Projects"></a>
The `enterprise_projects` block supports:

* `id` - The ID of the enterprise project.

* `name` - The name of the enterprise project.

* `status` - The status of the enterprise project.

* `type` - The type of the enterprise project.

* `description` - The description of the enterprise project.

* `created_at` - The time (UTC) when the enterprise project was created. Example: 2023-11-15T06:49:06Z

* `updated_at` - The time (UTC) when the enterprise project was modified. Example: 2023-11-16T02:21:36Z
