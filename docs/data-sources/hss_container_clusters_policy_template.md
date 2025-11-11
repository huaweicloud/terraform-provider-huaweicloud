---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_clusters_policy_template"
description: |-
  Use this data source to get the HSS container clusters policy template within HuaweiCloud.
---

# huaweicloud_hss_container_clusters_policy_template

Use this data source to get the HSS container clusters policy template within HuaweiCloud.

## Example Usage

```hcl
variable "policy_template_id" {}

data "huaweicloud_hss_container_clusters_policy_template" "test" {
  policy_template_id = var.policy_template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `policy_template_id` - (Required, String) Specifies the policy template ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, same as `policy_template_id`.

* `template_name` - The template name.

* `template_type` - The template type.

* `description` - The template description.

* `target_kind` - The policy template application resource type. Multiple resource types are separated by semicolons.

* `tag` - The tag.

* `level` - The recommendation level.

* `constraint_template` - The policy template content.
