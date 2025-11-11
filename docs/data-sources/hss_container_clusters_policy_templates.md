---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_clusters_policy_templates"
description: |-
  Use this data source to get the list of HSS container clusters protection policy templates within HuaweiCloud.
---

# huaweicloud_hss_container_clusters_policy_templates

Use this data source to get the list of HSS container clusters protection policy templates within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_clusters_policy_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `template_name` - (Optional, String) Specifies the template name.

* `template_type` - (Optional, String) Specifies the template type.

* `target_kind` - (Optional, String) Specifies the policy template application resource type.
  Multiple resource types are separated by semicolons.

* `tag` - (Optional, String) Specifies the tag.

* `level` - (Optional, String) Specifies the recommendation level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of policy templates.

* `data_list` - The list of policy templates.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `id` - The template ID.

* `template_name` - The template name.

* `template_type` - The template type.

* `description` - The template description.

* `target_kind` - The policy template application resource type. Multiple resource types are separated by semicolons.

* `tag` - The tag.

* `level` - The recommendation level.

* `constraint_template` - The policy template content.
