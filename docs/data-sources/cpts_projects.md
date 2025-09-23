---
subcategory: "Cloud Performance Test Service (CPTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpts_projects"
description: |-
  Use this datasource to get a list of CPTS projects within HuaweiCloud.
---

# huaweicloud_cpts_projects

Use this datasource to get a list of CPTS projects within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cpts_projects" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `projects` - The project details.

  The [projects](#projects_struct) structure is documented below.

<a name="projects_struct"></a>
The `projects` block supports:

* `name` - The project name.

* `source` - The project source.

* `variables_no_file` - The file variable.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `description` - The description.

* `id` - The project ID.
