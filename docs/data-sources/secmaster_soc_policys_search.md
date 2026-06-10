---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_policys_search"
description: |-
  Use this data source to search SOC policys within HuaweiCloud.
---

# huaweicloud_secmaster_soc_policys_search

Use this data source to search SOC policys within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_soc_policys_search" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `condition` - (Optional, List) Specifies the query conditions.

  The [condition](#policys_condition_struct) structure is documented below.

* `sort` - (Optional, List) Specifies the sort conditions.

  The [sort](#policys_sort_struct) structure is documented below.

* `group_by` - (Optional, List) Specifies the aggregation conditions.

  The [group_by](#policys_group_by_struct) structure is documented below.

<a name="policys_condition_struct"></a>
The `condition` block supports:

* `conditions` - (Optional, List) Specifies the list of query conditions.

  The [conditions](#policys_conditions_struct) structure is documented below.

* `logics` - (Optional, List) Specifies the list of condition names.

<a name="policys_conditions_struct"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the condition name.

* `data` - (Optional, List) Specifies the condition values.

<a name="policys_sort_struct"></a>
The `sort` block supports:

* `sort_by` - (Optional, String) Specifies the sort field.

* `order` - (Optional, String) Specifies the sort direction.

<a name="policys_group_by_struct"></a>
The `group_by` block supports:

* `group_by_fields` - (Optional, List) Specifies the aggregation fields.

* `group_by_hit` - (Optional, List) Specifies the aggregation result mapping.

  The [group_by_hit](#policys_group_by_hit_struct) structure is documented below.

<a name="policys_group_by_hit_struct"></a>
The `group_by_hit` block supports:

* `source` - (Optional, String) Specifies the source field.

* `dest` - (Optional, String) Specifies the destination field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The policys data in JSON string format.
