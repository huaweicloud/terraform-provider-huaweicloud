---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_region_projects"
description: |-
  Use this data source to query the project information of a tenant within HuaweiCloud.
---

# huaweicloud_cbr_region_projects

Use this data source to query the project information of a tenant within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_region_projects" "example" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the datasource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `projects` - The list of project information.
 The [projects](#projects_struct) structure is documented below.

* `links` - The link address.
  The [links](#links_struct) structure is documented below.

<a name="projects_struct"></a>
The `projects` block supports:

* `domain_id` - The domain ID.

* `is_domain` - The domain level sign.

* `parent_id` - The ID of the specific project or account ID of a specific system project.

* `name` - The backup region name.

* `description` - The description.

* `id` - The project ID.

* `enabled` - The enabling status of the project.

* `links` - The link address.
  The [links](#links_struct) structure is documented below.

<a name="links_struct"></a>
The `links` block supports:

* `self` - The link address.
