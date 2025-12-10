---
subcategory: "Security Master (SecMaster)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_baseline_checkitems"
description: |-
  Use this data source to get the list of SecMaster baseline checkitems.
---

# huaweicloud_secmaster_baseline_checkitems

Use this data source to get the list of SecMaster baseline checkitems.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_baseline_checkitems" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `catalog_uuid` - (Optional, String) Specifies the catalog UUID.

* `compliance_package_id` - (Optional, String) Specifies the compliance package ID.

* `sort_by` - (Optional, String) Specifies the field to sort by.

* `order` - (Optional, String) Specifies the sorting order. Valid values are **DESC** and **ASC**.

* `name` - (Optional, String) Specifies the checkitem name for filtering.

* `suggestion` - (Optional, String) Specifies the checkitem suggestion for filtering.

* `type` - (Optional, Int) Specifies the checkitem type for filtering. Valid values are:
  + **0**: Built-in.
  + **1**: Customized.

* `source_list` - (Optional, List) Specifies the source list for filtering. Valid values are:
  + **0**: Kotlin.
  + **2**: Playbook flow.
  + **3**: Manual.
  + **4**: Host access.

* `condition` - (Optional, List) Specifies the search condition.

  The [condition](#condition_struct) structure is documented below.

* `query_mode` - (Optional, String) Specifies the query mode.

* `severity` - (Optional, String) Specifies the severity level for filtering. Valid values are:
  + **Fatal**: Fatal.
  + **High**: High.
  + **Medium**: Medium.
  + **Low**: Low.
  + **Tips**: Tips.

<a name="condition_struct"></a>
The `condition` block supports:

* `conditions` - (Optional, List) Specifies the list of expressions.

  The [conditions](#conditions_struct) structure is documented below.

* `logics` - (Optional, List) Specifies the list of expression names.

<a name="conditions_struct"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the expression name.

* `data` - (Optional, List) Specifies the list of expression contents.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `builtin_checkitem_num` - The number of built-in checkitems.

* `checkitem_num` - The total number of checkitems.

* `customized_checkitem_num` - The number of customized checkitems.

* `checkitems` - The list of checkitems.

  The [checkitems](#checkitems_struct) structure is documented below.

<a name="checkitems_struct"></a>
The `checkitems` block supports:

* `aggregation_handle_status` - The aggregation handle status.

* `audit_procedure` - The audit procedure.

* `impact` - The impact.

* `cloud_server` - The cloud service.

* `description` - The description.

* `level` - The severity level. Valid values are:
  + **informational**
  + **low**
  + **medium**
  + **high**
  + **fatal**

* `method` - The check method. Valid values are:
  + **0**: Automatic item.
  + **3**: Script process/logic app.

* `name` - The checkitem name.

* `source` - The source. Valid values are:
  + **0**: default.
  + **3**: playbook.

* `workflow_id` - The workflow ID.

* `spec_checkitem_list` - The specification checkitem list.

  The [spec_checkitem_list](#spec_checkitem_list_struct) structure is documented below.

<a name="spec_checkitem_list_struct"></a>
The `spec_checkitem_list` block supports:

* `checkitem_uuid` - The checkitem UUID.

* `create_time` - The creation time.

* `language` - The language.

* `name` - The specification name.

* `remove_time` - The removal time.

* `specification_uuid` - The specification UUID.

* `uuid` - The UUID.
