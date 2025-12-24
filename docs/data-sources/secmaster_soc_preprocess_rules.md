---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_preprocess_rules"
description: |-
  Use this data source to get the list of the preprocess rules.
---

# huaweicloud_secmaster_soc_preprocess_rules

Use this data source to get the list of the preprocess rules.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_soc_preprocess_rules" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the preprocess rule name.

* `mapping_id` - (Optional, String) Specifies the mapping ID.

* `mapper_ids` - (Optional, List) Specifies the mapper ID list.

* `start_time` - (Optional, String) Specifies the query start time.

* `end_time` - (Optional, String) Specifies the query end time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of the preprocess rules.

  The [data](#rules_data_struct) structure is documented below.

<a name="rules_data_struct"></a>
The `data` block supports:

* `id` - The preprocess rule ID.

* `name` - The preprocess rule name.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `mapping_id` - The mapping ID.

* `mapper_id` - The mapper ID.

* `mapper_type_id` - The mapper type ID.

* `action` - The preprocessing options.

* `expression` - The expression.

* `create_time` - The creation time.

* `update_time` - The update time.
