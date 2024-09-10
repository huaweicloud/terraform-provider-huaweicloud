---
subcategory: "Cloud Performance Test Service (CPTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpts_project"
description: ""
---

# huaweicloud_cpts_project

Manages a project resource within HuaweiCloud CPTS.

## Example Usage

```hcl
resource "huaweicloud_cpts_project" "test" {
  name        = "project_name"
  description = "cpts project description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the project resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies a name for the project, which can contain a maximum of `42` characters.

* `description` - (Optional, String) Specifies the description of the project, which can contain a maximum of
  `50` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.

## Import

Projects can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cpts_project.test 1090
```
