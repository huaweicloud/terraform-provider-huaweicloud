---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_security_class"
description: |-
  Use this data source to get the security class list within HuaweiCloud.
---

# huaweicloud_dsc_security_class

Use this data source to get the security class list within HuaweiCloud.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_dsc_security_class" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `template_id` - (Required, String) Specifies the template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `classification_trees` - The classification tree list.

  The [classification_trees](#classification_trees_struct) structure is documented below.

* `template_category` - The template category.

* `template_name` - The template name.

<a name="classification_trees_struct"></a>
The `classification_trees` block supports:

* `children` - The child classification list.

  The [children](#classification_trees_children_struct) structure is documented below.

* `children_nums` - The number of child classifications.

* `create_time` - The creation time.

* `depth` - The classification depth.

* `id` - The classification ID.

* `name` - The classification name.

* `parent_id` - The parent classification ID.

* `project_id` - The project ID.

* `root_id` - The root classification ID.

* `rule_nums` - The number of rules.

* `update_time` - The update time.

<a name="classification_trees_children_struct"></a>
The `children` block supports:

* `children_nums` - The number of child classifications.

* `create_time` - The creation time.

* `depth` - The classification depth.

* `id` - The classification ID.

* `name` - The classification name.

* `parent_id` - The parent classification ID.

* `project_id` - The project ID.

* `root_id` - The root classification ID.

* `rule_nums` - The number of rules.

* `update_time` - The update time.
